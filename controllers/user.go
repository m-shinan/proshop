package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m-shinan/project-shop/auth"
	"github.com/m-shinan/project-shop/database"
	"github.com/m-shinan/project-shop/models"
	"golang.org/x/crypto/bcrypt"
)

type UserData struct {
	Firstname   string `form:"first_name" binding:"required"`
	Lastname    string `form:"last_name" binding:"required"`
	Email       string `form:"email" binding:"required,email"`
	Password    string `form:"password" binding:"required"`
	PhoneNumber int    `form:"phone" binding:"required"`
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Data.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Hashing password error",
		})
		return
	}

	db := database.DB

	// Check if user with the same email already exists
	var existingUserID uint
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", Data.Email).Scan(&existingUserID)
	switch {
	case err == sql.ErrNoRows: // No user with this email exists, proceed to insert
		otp := generateOTP()
		otpExpiry := time.Now().Add(15 * time.Minute) // OTP valid for 15 minutes

		// Insert user and OTP data
		_, err = db.Exec(
			"INSERT INTO users (first_name, last_name, email, password, phone, otp, otp_expiry) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			Data.Firstname, Data.Lastname, Data.Email, string(hashedPassword), Data.PhoneNumber, otp, otpExpiry,
		)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
			return
		}

		// Send OTP to user's email
		err = sendOTPEmail(Data.Email, otp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to send OTP email",
			})
			return
		}

		c.Redirect(http.StatusSeeOther, "/user/verifyotp")

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

// Generate a random OTP code
func generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// Send OTP email
func sendOTPEmail(email, otp string) error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAILPASS")
	to := email
	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is %s. It is valid for 15 minutes.", otp)

	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	log.Printf("Sending OTP to %s", email)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, msg)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}
	log.Printf("OTP email sent successfully to %s", email)
	return nil
}

func VerifyOTP(c *gin.Context) {
	var Data struct {
		Email string `form:"email" binding:"required,email"`
		OTP   string `form:"otp" binding:"required"`
	}

	if err := c.ShouldBind(&Data); err != nil {
		c.JSON(400, gin.H{"error": "Data binding error"})
		return
	}

	db := database.DB

	// Fetch the OTP and expiry from the database
	var storedOTP string
	var otpExpiry time.Time
	err := db.QueryRow("SELECT otp, otp_expiry FROM users WHERE email = $1", Data.Email).Scan(&storedOTP, &otpExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch OTP"})
		return
	}

	if time.Now().After(otpExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired"})
		return
	}

	if storedOTP != Data.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	// OTP is valid, update user status to verified or remove OTP
	_, err = db.Exec("UPDATE users SET otp = NULL, otp_expiry = NULL WHERE email = $1", Data.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully. You will be redirected to the login page."})
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

	c.Redirect(http.StatusSeeOther, "/")
}

func UserLogout(c *gin.Context) {
	c.SetCookie("UserAuthorization", "", -1, "", "", false, false)
	c.Redirect(http.StatusSeeOther, "/")
}

func UserHome(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	categoryID := c.DefaultQuery("category", "")
	searchQuery := c.DefaultQuery("search", "")

	products, err := GetProducts(c, userID, categoryID, searchQuery)
	if err {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	categories, err := GetCategories()
	if err {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.HTML(http.StatusOK, "user_home.html", gin.H{
		"Products":   products,
		"Categories": categories,
	})
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

////////////////// PROFILE //////////////////////////////

func UserProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	intUserID := int(userID)
	fmt.Println("UserProfile function called for userID:", intUserID)

	var user models.Users
	err := database.DB.QueryRow("SELECT id, first_name, last_name, email, phone, is_admin, is_blocked, created_at, updated_at FROM users WHERE id = $1", intUserID).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.IsAdmin, &user.Isblocked, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch user data"})
		return
	}

	var addresses []models.Address
	rows, err := database.DB.Query("SELECT id, user_id, name, mobile_number, home_name, place, landmark, city, district, state, country, postal_code FROM addresses WHERE user_id = $1", userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Unable to fetch addresses"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var address models.Address
		if err := rows.Scan(&address.ID, &address.UserID, &address.Name, &address.MobileNumber, &address.HomeName, &address.Place, &address.Landmark, &address.City, &address.District, &address.State, &address.Country, &address.PostalCode); err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Error scanning address data"})
			return
		}
		addresses = append(addresses, address)
	}

	// Render profile page with user data and addresses
	c.HTML(http.StatusOK, "user_profile.html", gin.H{
		"User":      user,
		"Addresses": addresses,
	})
}

func EditUserProfile(c *gin.Context) {
	userID, ok := c.MustGet("userID").(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve userID"})
		return
	}
	intUserID := int(userID) // Convert to int if necessary for database operations

	firstName := c.PostForm("firstName")
	lastName := c.PostForm("lastName")
	email := c.PostForm("email")
	phone := c.PostForm("phone")

	_, err := database.DB.Exec("UPDATE users SET first_name = $1, last_name = $2, email = $3, phone = $4 WHERE id = $5", firstName, lastName, email, phone, intUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update user data"})
		return
	}

	// Redirect to the profile page or another appropriate route
	c.Redirect(http.StatusSeeOther, "/user/profile")
}

//////////////////////  CHANGE PASSWORD  //////////////////

func UserChangePassword(c *gin.Context) {
	// Retrieve userID from context
	userID, ok := c.MustGet("userID").(uint)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/user/profile") // Redirect if there's an issue
		return
	}

	// Get current password and new password from form
	currentPassword := c.PostForm("currentPassword")
	newPassword := c.PostForm("newPassword")

	// Fetch user from database
	var user models.Users
	err := database.DB.QueryRow("SELECT id, password FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Password)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/user/profile?error=Unable to fetch user data")
		return
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/user/profile?error=Current password is incorrect")
		return
	}

	// Validate new password (e.g., minimum length)
	if len(newPassword) < 6 {
		c.Redirect(http.StatusSeeOther, "/user/profile?error=New password must be at least 6 characters long")
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/user/profile?error=Unable to hash new password")
		return
	}

	// Update password in the database
	_, err = database.DB.Exec("UPDATE users SET password = $1 WHERE id = $2", hashedPassword, userID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/user/profile?error=Unable to update password")
		return
	}

	// Redirect with success message
	c.Redirect(http.StatusSeeOther, "/user/profile?success=Password updated successfully")
}
