package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m-shinan/project-shop/database"
	"github.com/m-shinan/project-shop/models"
)

type CODOrderRequest struct {
	AddressID   uint    `json:"address_id" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required"`
	CouponCode  string  `json:"coupon_code"`
}

func CreateOrderCod(c *gin.Context) {
	var req CODOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("Received COD Order Request:", req)

	userID, exists := c.MustGet("userID").(uint)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	discountedAmount, err := ApplyDiscount(req.TotalAmount, req.CouponCode, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply discount"})
		return
	}

	cartItems, err := GetCartItems(userID)
	if err != nil || len(cartItems) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve cart items"})
		return
	}

	orderID, err := CreateOrder(userID, discountedAmount, req.AddressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create order"})
		return
	}

	for _, item := range cartItems {
		err = CreateOrderItem(orderID, item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create order items"})
			return
		}
	}

	err = ClearCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Order placed successfully",
		"order_id":     orderID,
		"redirect_url": "/user/orders",
	})
}

func ApplyDiscount(totalAmount float64, couponCode string, userID uint) (float64, error) {
	if couponCode == "" {
		return totalAmount, nil
	}

	query := `SELECT discount_amount FROM coupons WHERE coupon_code = $1 AND expiry_date > NOW()`
	var discountAmount float64
	err := database.DB.QueryRow(query, couponCode).Scan(&discountAmount)
	if err != nil {
		return totalAmount, err
	}

	discountedAmount := totalAmount - discountAmount
	if discountedAmount < 0 {
		discountedAmount = 0
	}
	return discountedAmount, nil
}

func GetCartItems(userID uint) ([]models.CartItem, error) {
	query := `SELECT p.id, p.product_name, p.price, c.quantity
              FROM cart c
              JOIN products p ON c.product_id = p.id
              WHERE c.user_id = $1`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []models.CartItem
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(&item.ProductID, &item.ProductName, &item.ProductPrice, &item.Quantity)
		if err != nil {
			return nil, err
		}
		item.TotalPrice = item.ProductPrice * float64(item.Quantity)
		cartItems = append(cartItems, item)
	}
	return cartItems, nil
}

func CreateOrder(userID uint, amount float64, addressID uint) (uint, error) {
	query := `INSERT INTO orders (user_id, amount, currency, status, payment_method, address_id, created_at, updated_at)
              VALUES ($1, $2, 'INR', 'pending', 'COD', $3, $4, $5)
              RETURNING id`

	var orderID uint
	err := database.DB.QueryRow(query, userID, amount, addressID, time.Now(), time.Now()).Scan(&orderID)
	return orderID, err
}

func CreateOrderItem(orderID uint, item models.CartItem) error {
	query := `INSERT INTO order_items (order_id, product_id, quantity, amount)
              VALUES ($1, $2, $3, $4)`

	_, err := database.DB.Exec(query, orderID, item.ProductID, item.Quantity, item.TotalPrice)
	return err
}

func ClearCart(userID uint) error {
	query := `DELETE FROM cart WHERE user_id = $1`
	_, err := database.DB.Exec(query, userID)
	return err
}
