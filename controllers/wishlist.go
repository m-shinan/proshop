package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m-shinan/project-shop/database"
	"github.com/m-shinan/project-shop/models"
)

func Wishlist(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	db := database.DB
	rows, err := db.Query(`
        SELECT p.id, p.product_name, p.description, p.stock, p.price, p.category_id, p.image,
               COALESCE(cart.product_id IS NOT NULL, FALSE) AS is_in_cart
        FROM products p
        INNER JOIN wishlist w ON p.id = w.product_id
        LEFT JOIN cart ON p.id = cart.product_id AND cart.user_id = $1
        WHERE w.user_id = $1
        ORDER BY w.id DESC
    `, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wishlist"})
		return
	}
	defer rows.Close()

	var products []models.Products
	for rows.Next() {
		var product models.Products
		var isInCart bool
		if err := rows.Scan(&product.ID, &product.ProductName, &product.Description, &product.Stock, &product.Price, &product.CategoryId, &product.Image, &isInCart); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product"})
			return
		}
		product.IsInCart = isInCart
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows"})
		return
	}

	c.HTML(http.StatusOK, "wishlist.html", gin.H{"Products": products})
}

func AddToWishlist(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	db := database.DB
	_, err = db.Exec(`
        INSERT INTO wishlist (user_id, product_id) 
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING
    `, userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to wishlist"})
		return
	}



	c.Redirect(http.StatusFound, "/user/userHome")
}

func RemoveFromWishlist(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	db := database.DB
	_, err = db.Exec(`
        DELETE FROM wishlist 
        WHERE user_id = $1 AND product_id = $2
    `, userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from wishlist"})
		return
	}

	

	c.Redirect(http.StatusFound, "/user/userHome")
}
