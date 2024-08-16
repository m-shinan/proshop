package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m-shinan/project-shop/database"
	"github.com/razorpay/razorpay-go"
)

var razorpayKey = "rzp_test_hy3toARxu0rJLv"
var razorpaySecret = "BH0gPLym0gkIfTMOEymBzZpg"

func CreateOrderRazorPay(c *gin.Context) {
	log.Println("CreateOrderRazorPay called")

	client := razorpay.NewClient(razorpayKey, razorpaySecret)

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		AddressID   uint    `json:"address_id" binding:"required"`
		TotalAmount float64 `json:"total_amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Request binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Request data: AddressID=%d, TotalAmount=%f\n", req.AddressID, req.TotalAmount)

	user := struct {
		Email   string
		Contact string
	}{}
	query := "SELECT email, phone FROM users WHERE id=$1"
	if err := database.DB.QueryRow(query, userID).Scan(&user.Email, &user.Contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
		return
	}

	amount := req.TotalAmount * 100
	orderData := map[string]interface{}{
		"amount":   int(amount),
		"currency": "INR",
		"receipt":  "order_rcptid_" + time.Now().Format("20060102150405"),
	}

	order, err := client.Order.Create(orderData, nil)
	if err != nil {
		log.Println("Error creating Razorpay order:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Razorpay order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orderID":      order["id"],
		"amount":       order["amount"],
		"user":         gin.H{"email": user.Email, "contact": user.Contact},
		"redirect_url": "/user/orders",
	})
}

func ConfirmRazorpayPayment(c *gin.Context) {
	log.Println("ConfirmRazorpayPayment called")

	var req struct {
		RazorpayOrderID   string  `json:"razorpay_order_id" binding:"required"`
		RazorpayPaymentID string  `json:"razorpay_payment_id" binding:"required"`
		RazorpaySignature string  `json:"razorpay_signature" binding:"required"`
		AddressID         uint    `json:"address_id" binding:"required"`
		TotalAmount       float64 `json:"total_amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Request binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify Razorpay signature here
	signature := req.RazorpaySignature
	message := req.RazorpayOrderID + "|" + req.RazorpayPaymentID
	hash := hmac.New(sha256.New, []byte(razorpaySecret))
	hash.Write([]byte(message))
	expectedSignature := hex.EncodeToString(hash.Sum(nil))

	if signature != expectedSignature {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Payment verification failed"})
		return
	}

	// Create order in database
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Save order details and create order items
	cartItems, err := GetCartItems(userID.(uint))
	if err != nil || len(cartItems) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve cart items"})
		return
	}

	orderID, err := CreateOrderRaz(userID.(uint), req.TotalAmount, req.AddressID, "Razorpay")
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

	err = ClearCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Payment verified and order created successfully", "orderID": orderID})
}

func CreateOrderRaz(userID uint, amount float64, addressID uint, paymentMethod string) (uint, error) {
	query := `INSERT INTO orders (user_id, amount, currency, status, payment_method, address_id, created_at, updated_at)
              VALUES ($1, $2, 'INR', 'pending', $3, $4, $5, $6)
              RETURNING id`

	var orderID uint
	err := database.DB.QueryRow(query, userID, amount, paymentMethod, addressID, time.Now(), time.Now()).Scan(&orderID)
	return orderID, err
}
