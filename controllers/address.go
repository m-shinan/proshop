package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-shinan/project-shop/database"
	"github.com/m-shinan/project-shop/models"
)

// type AddressData struct {
// 	UserID       int    `json:"user_id"`
// 	Name         string `json:"name" form:"name" binding:"required"`
// 	MobileNumber string `json:"mobileNumber" form:"mobileNumber" binding:"required"`
// 	HomeName     string `json:"homeName" form:"homeName" binding:"required"`
// 	Place        string `json:"place" form:"place" binding:"required"`
// 	Landmark     string `json:"landmark" form:"landmark"`
// 	City         string `json:"city" form:"city" binding:"required"`
// 	District     string `json:"district" form:"district" binding:"required"`
// 	State        string `json:"state" form:"state" binding:"required"`
// 	Country      string `json:"country" form:"country" binding:"required"`
// 	PostalCode   string `json:"postalCode" form:"postalCode" binding:"required"`
// }

// func AddAddress(c *gin.Context) {
// 	userID, _ := c.Get("userID") // Use this if the userID is set in the context
// 	var address AddressData

// 	if err := c.Bind(&address); err != nil {
// 		c.HTML(http.StatusBadRequest, "user_profile.html", gin.H{"AddressError": "Invalid input data"})
// 		return
// 	}

// 	intUserID, ok := userID.(int)
// 	if !ok {
// 		c.HTML(http.StatusInternalServerError, "user_profile.html", gin.H{"AddressError": "User ID is missing"})
// 		return
// 	}
// 	address.UserID = intUserID

// 	_, err := database.DB.Exec(`INSERT INTO addresses (user_id, name, mobile_number, home_name, place, landmark, city, district, state, country, postal_code)
//                                 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
// 		address.UserID, address.Name, address.MobileNumber, address.HomeName, address.Place, address.Landmark, address.City, address.District, address.State, address.Country, address.PostalCode)
// 	if err != nil {
// 		c.HTML(http.StatusInternalServerError, "user_profile.html", gin.H{"AddressError": "Unable to add address"})
// 		return
// 	}

// 	c.Redirect(http.StatusSeeOther, "/user/profile")
// }

// func AddAddress(c *gin.Context) {
// 	userID, _ := c.Get("userID") // Use this if the userID is set in the context
// 	var address AddressData

// 	if err := c.Bind(&address); err != nil {
// 		log.Printf("Bind error: %v", err)
// 		c.HTML(http.StatusBadRequest, "user_profile.html", gin.H{"AddressError": "Invalid input data"})
// 		return
// 	}

// 	intUserID, ok := userID.(int)
// 	if !ok {
// 		c.HTML(http.StatusInternalServerError, "user_profile.html", gin.H{"AddressError": "User ID is missing"})
// 		return
// 	}
// 	address.UserID = intUserID

// 	log.Printf("Address data: %+v", address)

// 	_, err := database.DB.Exec(`INSERT INTO addresses (user_id, name, mobile_number, home_name, place, landmark, city, district, state, country, postal_code)
//                                 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
// 		address.UserID, address.Name, address.MobileNumber, address.HomeName, address.Place, address.Landmark, address.City, address.District, address.State, address.Country, address.PostalCode)
// 	if err != nil {
// 		c.HTML(http.StatusInternalServerError, "user_profile.html", gin.H{"AddressError": "Unable to add address"})
// 		return
// 	}

// 	c.Redirect(http.StatusSeeOther, "/user/profile")
// }

// func AddAddress(c *gin.Context) {
// 	userID, _ := c.Get("userID")
// 	var address AddressData

// 	if err := c.Bind(&address); err != nil {
// 		log.Printf("Bind error: %v", err)
// 		c.HTML(http.StatusBadRequest, "user_profile.html", gin.H{"AddressError": "Invalid input data"})
// 		return
// 	}

// 	intUserID, ok := userID.(int)
// 	if !ok {
// 		log.Println("User ID is missing")
// 		c.HTML(http.StatusInternalServerError, "user_profile.html", gin.H{"AddressError": "User ID is missing"})
// 		return
// 	}
// 	address.UserID = intUserID

// 	log.Printf("Address data: %+v", address)

// 	_, err := database.DB.Exec(`INSERT INTO addresses (user_id, name, mobile_number, home_name, place, landmark, city, district, state, country, postal_code)
//                                 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
// 		address.UserID, address.Name, address.MobileNumber, address.HomeName, address.Place, address.Landmark, address.City, address.District, address.State, address.Country, address.PostalCode)
// 	if err != nil {
// 		log.Printf("Database error: %v", err)
// 		c.HTML(http.StatusInternalServerError, "user_profile.html", gin.H{"AddressError": "Unable to add address"})
// 		return
// 	}

// 	c.Redirect(http.StatusSeeOther, "/user/profile")
// }

// func AddAddress(c *gin.Context) {
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		log.Println("User ID is missing")
// 		c.HTML(http.StatusInternalServerError, "user_profile.html", gin.H{"AddressError": "User ID is missing"})
// 		return
// 	}

// 	var address AddressData
// 	if err := c.Bind(&address); err != nil {
// 		log.Printf("Bind error: %v", err)
// 		c.HTML(http.StatusBadRequest, "user_profile.html", gin.H{"AddressError": "Invalid input data"})
// 		return
// 	}

// 	intUserID, ok := userID.(int)
// 	if !ok {
// 		log.Println("User ID is not an integer")
// 		c.HTML(http.StatusInternalServerError, "user_profile.html", gin.H{"AddressError": "User ID is invalid"})
// 		return
// 	}
// 	address.UserID = intUserID

// 	log.Printf("Address data: %+v", address)

// 	_, err := database.DB.Exec(`INSERT INTO addresses (user_id, name, mobile_number, home_name, place, landmark, city, district, state, country, postal_code)
//                                 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
// 		address.UserID, address.Name, address.MobileNumber, address.HomeName, address.Place, address.Landmark, address.City, address.District, address.State, address.Country, address.PostalCode)
// 	if err != nil {
// 		log.Printf("Database error: %v", err)
// 		c.HTML(http.StatusInternalServerError, "user_profile.html", gin.H{"AddressError": "Unable to add address"})
// 		return
// 	}

//		c.Redirect(http.StatusSeeOther, "/user/profile")
//	}
//
//	type AddressData struct {
//		UserID       int    `json:"user_id"`
//		Name         string `json:"name" form:"name" binding:"required"`
//		MobileNumber string `json:"mobileNumber" form:"mobileNumber" binding:"required"`
//		HomeName     string `json:"homeName" form:"homeName" binding:"required"`
//		Place        string `json:"place" form:"place" binding:"required"`
//		Landmark     string `json:"landmark" form:"landmark"`
//		City         string `json:"city" form:"city" binding:"required"`
//		District     string `json:"district" form:"district" binding:"required"`
//		State        string `json:"state" form:"state" binding:"required"`
//		Country      string `json:"country" form:"country" binding:"required"`
//		PostalCode   string `json:"postalCode" form:"postalCode" binding:"required"`
//	}
type AddressData struct {
	UserID       int    `json:"user_id"`
	Name         string `json:"name" form:"name" binding:"required"`
	MobileNumber string `json:"mobileNumber" form:"mobileNumber" binding:"required"`
	HomeName     string `json:"homeName" form:"homeName" binding:"required"`
	Place        string `json:"place" form:"place" binding:"required"`
	Landmark     string `json:"landmark" form:"landmark"`
	City         string `json:"city" form:"city" binding:"required"`
	District     string `json:"district" form:"district" binding:"required"`
	State        string `json:"state" form:"state" binding:"required"`
	Country      string `json:"country" form:"country" binding:"required"`
	PostalCode   string `json:"postalCode" form:"postalCode" binding:"required"`
}

// AddAddress handles the POST request to add a new address
func AddAddress(c *gin.Context) {
	// Get the userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("User ID is missing")
		c.HTML(http.StatusInternalServerError, "user_profile.html", gin.H{"AddressError": "User ID is missing"})
		return
	}

	// Convert userID to integer
	var intUserID uint
	switch v := userID.(type) {
	case float64:
		intUserID = uint(v)
	case uint:
		intUserID = v
	case int:
		intUserID = uint(v)
	default:
		log.Printf("Unexpected type for User ID: %T", v)
		c.HTML(http.StatusInternalServerError, "user_profile.html", gin.H{"AddressError": "User ID is invalid"})
		return
	}

	// Bind and validate address data
	var address AddressData
	if err := c.Bind(&address); err != nil {
		log.Printf("Bind error: %v", err)
		c.HTML(http.StatusBadRequest, "user_profile.html", gin.H{"AddressError": "Invalid input data"})
		return
	}
	address.UserID = int(intUserID) // Ensure the UserID is set correctly

	// Insert address into the database
	_, err := database.DB.Exec(`INSERT INTO addresses (user_id, name, mobile_number, home_name, place, landmark, city, district, state, country, postal_code) 
                                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		address.UserID, address.Name, address.MobileNumber, address.HomeName, address.Place, address.Landmark, address.City, address.District, address.State, address.Country, address.PostalCode)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.HTML(http.StatusInternalServerError, "user_profile.html", gin.H{"AddressError": "Unable to add address"})
		return
	}

	// Redirect or respond with success
	c.Redirect(http.StatusSeeOther, "/user/profile")
}

func EditAddressPage(c *gin.Context) {
	addressID := c.Param("id")

	var address models.Address
	err := database.DB.QueryRow("SELECT id, user_id, name, mobile_number, home_name, place, landmark, city, district, state, country, postal_code FROM addresses WHERE id = $1", addressID).Scan(
		&address.ID, &address.UserID, &address.Name, &address.MobileNumber, &address.HomeName, &address.Place, &address.Landmark, &address.City, &address.District, &address.State, &address.Country, &address.PostalCode,
	)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Unable to fetch address data"})
		return
	}

	c.HTML(http.StatusOK, "edit_address.html", gin.H{
		"Address": address,
	})
}

func EditAddress(c *gin.Context) {
	addressID := c.Param("id")
	userID := c.MustGet("userID").(uint)

	name := c.PostForm("name")
	mobileNumber := c.PostForm("mobile_number")
	homeName := c.PostForm("home_name")
	place := c.PostForm("place")
	landmark := c.PostForm("landmark")
	city := c.PostForm("city")
	district := c.PostForm("district")
	state := c.PostForm("state")
	country := c.PostForm("country")
	postalCode := c.PostForm("postal_code")

	_, err := database.DB.Exec("UPDATE addresses SET name = $1, mobile_number = $2, home_name = $3, place = $4, landmark = $5, city = $6, district = $7, state = $8, country = $9, postal_code = $10 WHERE id = $11 AND user_id = $12",
		name, mobileNumber, homeName, place, landmark, city, district, state, country, postalCode, addressID, userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Unable to update address"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/user/profile")
}

func DeleteAddress(c *gin.Context) {
	addressID := c.Param("id")
	userID := c.MustGet("userID").(uint)

	_, err := database.DB.Exec("DELETE FROM addresses WHERE id = $1 AND user_id = $2", addressID, userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Unable to delete address"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/user/addresses")
}
