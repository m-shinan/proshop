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

type UserData struct {
	Firstname string `form:"first_name" binding:"required"`
	Lastname  string `form:"last_name" binding:"required"`
	Email     string `form:"email" binding:"required,email"`
	Password  string `form:"password" binding:"required"`

	PhoneNumber int `form:"phone" binding:"required"`
}

func UserSignUp(c *gin.Context) {
	var Data UserData
	if err := c.ShouldBind(&Data); err != nil {
		fmt.Println("Error binding data:", err)
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		return
	}

	// Print Data to check if binding is correct
	fmt.Println("Form data:", Data)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Data.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Hashing password error",
		})
		return
	}

	db := database.DB

	// Check if user with the same email already exists
	var existingAdminID uint
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", Data.Email).Scan(&existingAdminID)
	switch {
	case err == sql.ErrNoRows: // No user with this email exists, proceed to insert
		_, err = db.Exec(
			"INSERT INTO users (first_name, last_name, email, password, phone) VALUES ($1, $2, $3, $4, $5)",
			Data.Firstname, Data.Lastname, Data.Email, string(hashedPassword), Data.PhoneNumber, // assuming IsAdmin is false for signup
		)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})

			return
		}
		c.Redirect(http.StatusSeeOther, "/user/login")

	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return

	default: // User with this email already exists
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}
}

func UserLogin(c *gin.Context) {
	type UserData struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	var user UserData
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"error": "Login data binding error",
		})
		return
	}

	db := database.DB

	var checkUser models.Users
	err := db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", user.Email).Scan(&checkUser.ID, &checkUser.Email, &checkUser.Password)
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

	err = bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(501, gin.H{
			"error": "Username and password invalid",
		})
		return
	}

	// Generating a JWT-token
	str := strconv.Itoa(int(checkUser.ID))
	tokenString := auth.TokenGeneration(str)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAuthorization", tokenString, 3600*24*30, "", "", false, true)

	c.Redirect(http.StatusFound, "/user/uservalidate")
}

func ValidateUser(c *gin.Context) {
	c.Get("user")

	c.Redirect(http.StatusSeeOther, "/user/userHome")
}

func UserLogout(c *gin.Context) {
	c.SetCookie("UserAuthorization", "", -1, "", "", false, false)
	c.Redirect(http.StatusSeeOther, "/user/login")
}

func UserHome(c *gin.Context) {
	products, err := GetProducts(c)
	if err {
		// Error already handled inside GetProducts
		return
	}

	c.HTML(http.StatusOK, "user_home.html", gin.H{"Products": products})
}

//////////////////////////////  ADMIN USER MANAGEMENT ////////////////////////////

func AdminViewUsers(c *gin.Context) {
	db := database.DB

	// Fetch active users
	activeRows, err := db.Query(`
        SELECT id, first_name, last_name, email, phone, created_at, updated_at
        FROM users
        WHERE is_blocked = FALSE
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active users"})
		return
	}
	defer activeRows.Close()

	var activeUsers []models.Users
	for activeRows.Next() {
		var user models.Users
		if err := activeRows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan active user"})
			return
		}
		activeUsers = append(activeUsers, user)
	}

	// Fetch blocked users
	blockedRows, err := db.Query(`
        SELECT id, first_name, last_name, email, phone, created_at, updated_at
        FROM users
        WHERE is_blocked = TRUE
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blocked users"})
		return
	}
	defer blockedRows.Close()

	var blockedUsers []models.Users
	for blockedRows.Next() {
		var user models.Users
		if err := blockedRows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan blocked user"})
			return
		}
		blockedUsers = append(blockedUsers, user)
	}

	c.HTML(http.StatusOK, "admin_user_man.html", gin.H{
		"ActiveUsers":  activeUsers,
		"BlockedUsers": blockedUsers,
	})
}

func AdminBlockUsers(c *gin.Context) {
	userID := c.Param("id")
	db := database.DB

	_, err := db.Exec("UPDATE users SET is_blocked = TRUE WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to block user"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

func AdminUnblockUsers(c *gin.Context) {
	userID := c.Param("id")
	db := database.DB

	_, err := db.Exec("UPDATE users SET is_blocked = FALSE WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unblock user"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

func AdminDeleteUsers(c *gin.Context) {
	userID := c.Param("id")
	db := database.DB

	_, err := db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}
