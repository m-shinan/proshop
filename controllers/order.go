package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/m-shinan/project-shop/database"
)

// Order struct represents the order details
type Order struct {
	ID                int         `json:"id"` // Changed from OrderID to ID
	Amount            float64     `json:"amount"`
	Payment           string      `json:"payment"`
	Status            string      `json:"status"`
	CreatedAt         string      `json:"created_at"`
	AddressName       string      `json:"address_name"`
	AddressPlace      string      `json:"address_place"`
	AddressCity       string      `json:"address_city"`
	AddressDistrict   string      `json:"address_district"`
	AddressState      string      `json:"address_state"`
	AddressCountry    string      `json:"address_country"`
	AddressPostalCode string      `json:"address_postal_code"`
	Email             string      `json:"email"` // Added this field for user email
	OrderItems        []OrderItem `json:"order_items"`
}

func (o Order) FullAddress() string {
	return o.AddressName + ", " + o.AddressPlace + ", " + o.AddressCity + ", " + o.AddressDistrict + ", " + o.AddressState + ", " + o.AddressCountry + " - " + o.AddressPostalCode
}

// OrderItem struct represents an item in an order
type OrderItem struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Amount      float64 `json:"amount"`
}

// UserOrders fetches and displays the orders for the authenticated user
func UserOrders(c *gin.Context) {

	// Assuming user ID is stored in a session or JWT token
	userID := c.MustGet("userID").(uint) // Adjust according to your user session management

	// Query to fetch orders
	rows, err := database.DB.Query(`
		SELECT o.id, o.amount, o.payment_method, o.status, o.created_at, 
		       a.name, a.place, a.city, a.district, a.state, a.country, a.postal_code
		FROM orders o
		JOIN addresses a ON o.address_id = a.id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.Amount, &order.Payment, &order.Status, &order.CreatedAt,
			&order.AddressName, &order.AddressPlace, &order.AddressCity, &order.AddressDistrict,
			&order.AddressState, &order.AddressCountry, &order.AddressPostalCode,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order"})
			return
		}

		// Query to fetch order items
		itemRows, err := database.DB.Query(`
			SELECT p.product_name, oi.quantity, oi.amount
			FROM order_items oi
			JOIN products p ON oi.product_id = p.id
			WHERE oi.order_id = $1
		`, order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
			return
		}
		defer itemRows.Close()

		var items []OrderItem
		for itemRows.Next() {
			var item OrderItem
			err := itemRows.Scan(&item.ProductName, &item.Quantity, &item.Amount)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order item"})
				return
			}
			items = append(items, item)
		}
		order.OrderItems = items
		orders = append(orders, order)
	}

	c.HTML(http.StatusOK, "user_orders.html", gin.H{
		"orders": orders,
	})
}

// CancelOrder deletes an order from the orders table and redirects to the orders page
func CancelOrder(c *gin.Context) {
	userID := c.MustGet("userID").(uint) // Adjust according to your user session management

	// Retrieve the order_id from form data
	orderID := c.PostForm("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	// Check if the order exists and retrieve the owner ID
	var ownerID uint
	err := database.DB.QueryRow(`
		SELECT user_id
		FROM orders
		WHERE id = $1
	`, orderID).Scan(&ownerID)
	if err != nil {
		log.Printf("Error fetching order owner: %v", err) // Log the specific error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order owner"})
		return
	}

	// Check if the order belongs to the current user
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to cancel this order"})
		return
	}

	// Delete related order items
	_, err = database.DB.Exec(`
		DELETE FROM order_items
		WHERE order_id = $1
	`, orderID)
	if err != nil {
		log.Printf("Error deleting order items: %v", err) // Log the specific error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order items"})
		return
	}

	// Delete the order from the database
	_, err = database.DB.Exec(`
		DELETE FROM orders
		WHERE id = $1
	`, orderID)
	if err != nil {
		log.Printf("Error deleting order: %v", err) // Log the specific error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	// Redirect to the orders page
	c.Redirect(http.StatusSeeOther, "/user/orders")
}

// AdminOrders fetches and displays all orders with their details for the admin
// AdminOrders fetches and displays all orders with their items for the admin
func AdminOrders(c *gin.Context) {
	// Query to fetch orders along with their associated user emails
	rows, err := database.DB.Query(`
		SELECT o.id, o.amount, o.payment_method, o.status, o.created_at, 
		       a.name, a.place, a.city, a.district, a.state, a.country, a.postal_code, 
		       u.email
		FROM orders o
		JOIN addresses a ON o.address_id = a.id
		JOIN users u ON o.user_id = u.id
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.Amount, &order.Payment, &order.Status, &order.CreatedAt,
			&order.AddressName, &order.AddressPlace, &order.AddressCity, &order.AddressDistrict,
			&order.AddressState, &order.AddressCountry, &order.AddressPostalCode,
			&order.Email,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order"})
			return
		}

		// Query to fetch order items
		itemRows, err := database.DB.Query(`
			SELECT p.product_name, oi.quantity, oi.amount
			FROM order_items oi
			JOIN products p ON oi.product_id = p.id
			WHERE oi.order_id = $1
		`, order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
			return
		}
		defer itemRows.Close()

		var items []OrderItem
		for itemRows.Next() {
			var item OrderItem
			err := itemRows.Scan(&item.ProductName, &item.Quantity, &item.Amount)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order item"})
				return
			}
			items = append(items, item)
		}
		order.OrderItems = items
		orders = append(orders, order)
	}

	c.HTML(http.StatusOK, "admin_orders.html", gin.H{
		"Orders": orders,
	})
}

// UpdateOrderStatus updates the status of an order
func UpdateOrderStatus(c *gin.Context) {
	orderID := c.PostForm("order_id")
	newStatus := c.PostForm("new_status")

	// Update the status in the database
	_, err := database.DB.Exec(`
		UPDATE orders 
		SET status = $1 
		WHERE id = $2
	`, newStatus, orderID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/orders")
}

// AdminOrdersByUser fetches and displays orders for a specific user in the admin panel
func AdminOrdersByUser(c *gin.Context) {
	userID := c.Param("user_id") // Get user ID from the URL parameter

	// Query to fetch orders for the specified user
	rows, err := database.DB.Query(`
        SELECT o.id, o.amount, o.payment_method, o.status, o.created_at, 
               a.name, a.place, a.city, a.district, a.state, a.country, a.postal_code
        FROM orders o
        JOIN addresses a ON o.address_id = a.id
        WHERE o.user_id = $1
    `, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.Amount, &order.Payment, &order.Status, &order.CreatedAt,
			&order.AddressName, &order.AddressPlace, &order.AddressCity, &order.AddressDistrict,
			&order.AddressState, &order.AddressCountry, &order.AddressPostalCode,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order"})
			return
		}

		// Query to fetch order items
		itemRows, err := database.DB.Query(`
            SELECT p.product_name, oi.quantity, oi.amount
            FROM order_items oi
            JOIN products p ON oi.product_id = p.id
            WHERE oi.order_id = $1
        `, order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
			return
		}
		defer itemRows.Close()

		var items []OrderItem
		for itemRows.Next() {
			var item OrderItem
			err := itemRows.Scan(&item.ProductName, &item.Quantity, &item.Amount)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order item"})
				return
			}
			items = append(items, item)
		}
		order.OrderItems = items
		orders = append(orders, order)
	}

	// Render the admin orders page with the specific user's orders
	c.HTML(http.StatusOK, "admin_orders.html", gin.H{
		"Orders": orders,
	})
}
