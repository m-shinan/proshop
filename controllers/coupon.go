package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-shinan/project-shop/database"
	"github.com/m-shinan/project-shop/models"
)

func ViewCoupons(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, code, discount_amount, description, user_limit FROM coupons")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch coupons"})
		return
	}
	defer rows.Close()

	var coupons []models.Coupon
	for rows.Next() {
		var coupon models.Coupon
		err := rows.Scan(&coupon.ID, &coupon.Code, &coupon.DiscountAmount, &coupon.Description, &coupon.UserLimit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan coupon"})
			return
		}
		coupons = append(coupons, coupon)
	}

	c.HTML(http.StatusOK, "coupon.html", gin.H{
		"coupons": coupons,
	})
}

func AddCoupon(c *gin.Context) {
	var coupon models.Coupon

	// Bind form data to Coupon struct
	if err := c.ShouldBind(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "details": err.Error()})
		return
	}

	// Validate data
	if coupon.Code == "" || coupon.DiscountAmount <= 0 || coupon.UserLimit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Save coupon to the database (adjust query as needed for your DB setup)
	_, err := database.DB.Exec(`
		INSERT INTO coupons (code, discount_amount, description, user_limit)
		VALUES ($1, $2, $3, $4)
	`, coupon.Code, coupon.DiscountAmount, coupon.Description, coupon.UserLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save coupon", "details": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/coupons")
}

func EditCoupon(c *gin.Context) {
	id := c.PostForm("id")

	var coupon models.Coupon
	if err := c.ShouldBind(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	_, err := database.DB.Exec(
		"UPDATE coupons SET code = $1, discount_amount = $2, description = $3, user_limit = $4, updated_at = NOW() WHERE id = $5",
		coupon.Code, coupon.DiscountAmount, coupon.Description, coupon.UserLimit, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coupon"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/coupons")
}

func DeleteCoupon(c *gin.Context) {
	id := c.Param("id")

	_, err := database.DB.Exec("DELETE FROM coupons WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupon"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/coupons")

}

func EditCouponForm(c *gin.Context) {
	id := c.Param("id")

	var coupon models.Coupon
	err := database.DB.QueryRow(`
        SELECT id, code, discount_amount, description, user_limit
        FROM coupons
        WHERE id = $1
    `, id).Scan(&coupon.ID, &coupon.Code, &coupon.DiscountAmount, &coupon.Description, &coupon.UserLimit)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch coupon"})
		return
	}

	c.HTML(http.StatusOK, "coupon.html", gin.H{
		"EditingCoupon": true,
		"Coupon":        coupon,
	})
}
