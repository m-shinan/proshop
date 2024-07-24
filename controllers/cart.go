package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m-shinan/project-shop/database"
	"github.com/m-shinan/project-shop/models"
)

type CartData struct {
	CartItems []models.CartItem
	CartTotal float64
}

func Cart(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	log.Printf("Fetching cart for user ID: %d", userID)

	var cartData CartData
	query := `
        SELECT p.id, p.name, p.price, c.quantity 
        FROM cart c
        JOIN products p ON c.product_id = p.id
        WHERE c.user_id = $1`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart items"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(&item.ProductID, &item.ProductName, &item.ProductPrice, &item.Quantity)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan cart item"})
			return
		}
		item.TotalPrice = item.ProductPrice * float64(item.Quantity)
		cartData.CartTotal += item.TotalPrice
		cartData.CartItems = append(cartData.CartItems, item)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows"})
		return
	}

	log.Printf("Cart data retrieved successfully: %+v", cartData)
	c.HTML(http.StatusOK, "user_cart.html", cartData)
}

func PlusQuantity(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userID := c.MustGet("userID").(uint)

	query := `UPDATE cart SET quantity = quantity + 1 WHERE user_id = $1 AND product_id = $2`
	_, err = database.DB.Exec(query, userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increase quantity"})
		return
	}

	c.Redirect(http.StatusFound, "/user/cart")
}

// MinusQuantity decreases the quantity of a product in the cart
func MinusQuantity(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userID := c.MustGet("userID").(uint)

	query := `UPDATE cart SET quantity = CASE WHEN quantity > 1 THEN quantity - 1 ELSE quantity END WHERE user_id = $1 AND product_id = $2`
	_, err = database.DB.Exec(query, userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrease quantity"})
		return
	}

	c.Redirect(http.StatusFound, "/user/cart")
}
func AddToCart(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// userID, exists := c.Get("userID")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found"})
	// 	return
	// }

	// userIDStr, ok := userID.(string)
	// if !ok {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid User ID"})
	// 	return
	// }
	userID := c.MustGet("userID").(uint)

	query := `INSERT INTO cart (user_id, product_id) VALUES ($1, $2)`
	_, err = database.DB.Exec(query, userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart"})
}

func RemoveFromCart(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userID := c.MustGet("userID").(uint)

	query := `DELETE FROM cart WHERE user_id = $1 AND product_id = $2`
	_, err = database.DB.Exec(query, userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product removed from cart"})
}
