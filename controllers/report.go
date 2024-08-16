package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m-shinan/project-shop/database"
)

type CompletedOrder struct {
	ID           int
	CustomerName string
	Date         time.Time
	Total        float64
	Status       string
}

func SalesReport(c *gin.Context) {

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var (
		totalCompletedOrders int
		totalProductsSold    int
		totalEarnings        float64
		completedOrders      []CompletedOrder
	)

	query := `
        SELECT o.id, u.email, o.created_at, o.amount, o.status, SUM(oi.quantity)
        FROM orders o
        JOIN users u ON o.user_id = u.id
        JOIN order_items oi ON o.id = oi.order_id
        WHERE o.status = 'delivered'
    `

	if startDate != "" && endDate != "" {
		query += " AND o.created_at::date BETWEEN $1 AND $2"
	}

	query += " GROUP BY o.id, u.email"

	var rows *sql.Rows
	var err error

	if startDate != "" && endDate != "" {
		rows, err = database.DB.Query(query, startDate, endDate)
	} else {
		rows, err = database.DB.Query(query)
	}

	if err != nil {
		fmt.Println("Error querying database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving sales report"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order CompletedOrder
		var productsSold int

		err := rows.Scan(&order.ID, &order.CustomerName, &order.Date, &order.Total, &order.Status, &productsSold)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing sales report"})
			return
		}

		completedOrders = append(completedOrders, order)
		totalCompletedOrders++
		totalProductsSold += productsSold
		totalEarnings += order.Total
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Row iteration error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating sales report"})
		return
	}

	c.HTML(http.StatusOK, "report.html", gin.H{
		"TotalCompletedOrders": totalCompletedOrders,
		"TotalProductsSold":    totalProductsSold,
		"TotalEarnings":        totalEarnings,
		"CompletedOrders":      completedOrders,
	})
}
