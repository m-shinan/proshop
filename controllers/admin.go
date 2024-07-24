package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m-shinan/project-shop/auth"
	"github.com/m-shinan/project-shop/database"
	"github.com/m-shinan/project-shop/models"
	"golang.org/x/crypto/bcrypt"
)

type AdminData struct {
	Firstname       string `form:"first_name" binding:"required"`
	Lastname        string `form:"last_name" binding:"required"`
	Email           string `form:"email" binding:"required,email"`
	Password        string `form:"password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" binding:"required"`
	PhoneNumber     int    `form:"phone" binding:"required"`
}

func AdminSignUp(c *gin.Context) {
	var Data AdminData
	if err := c.ShouldBind(&Data); err != nil {
		fmt.Println("Error binding data:", err)
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		return
	}

	// Print Data to check if binding is correct
	fmt.Println("Form data:", Data)

	if Data.Password != Data.ConfirmPassword {
		c.JSON(400, gin.H{
			"error": "Passwords do not match",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Data.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Hashing password error",
		})
		return
	}

	db := database.DB

	// Check if admin with the same email already exists
	var existingAdminID uint
	err = db.QueryRow("SELECT id FROM admins WHERE email = $1", Data.Email).Scan(&existingAdminID)
	switch {
	case err == sql.ErrNoRows: // No admin with this email exists, proceed to insert
		_, err = db.Exec(
			"INSERT INTO admins (first_name, last_name, email, password, phone, is_admin) VALUES ($1, $2, $3, $4, $5, $6)",
			Data.Firstname, Data.Lastname, Data.Email, string(hashedPassword), Data.PhoneNumber, true, // assuming IsAdmin is true for signup
		)
		if err != nil {
			log.Printf("Error inserting admin: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create admin",
			})

			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/login")

	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return

	default: // Admin with this email already exists
		c.JSON(http.StatusConflict, gin.H{
			"error": "Admin with this email already exists",
		})
		return
	}
}

func AdminLogin(c *gin.Context) {
	type AdminData struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	var admin AdminData
	if err := c.ShouldBind(&admin); err != nil {
		c.JSON(400, gin.H{
			"error": "Login data binding error",
		})
		return
	}

	db := database.DB

	var checkAdmin models.Admins
	err := db.QueryRow("SELECT id, email, password FROM admins WHERE email = $1", admin.Email).Scan(&checkAdmin.ID, &checkAdmin.Email, &checkAdmin.Password)
	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{
			"error": "User not found",
		})
		return
	} else if err != nil {
		c.JSON(500, gin.H{
			"error": "Database query error",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(checkAdmin.Password), []byte(admin.Password))
	if err != nil {
		c.JSON(501, gin.H{
			"error": "Username and password invalid",
		})
		return
	}

	// Generating a JWT-token
	str := strconv.Itoa(int(checkAdmin.ID))
	tokenString := auth.TokenGeneration(str)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("AdminAuthorization", tokenString, 3600*24*30, "", "", false, true)

	c.Redirect(http.StatusFound, "/admin/adminvalidate")
}

func ValidateAdmin(c *gin.Context) {
	c.Get("admin")

	c.Redirect(http.StatusSeeOther, "/admin/adminHome")
}

func AdminLogout(c *gin.Context) {
	c.SetCookie("AdminAuthorization", "", -1, "", "", false, false)
	c.Redirect(http.StatusSeeOther, "/admin/login")
}
