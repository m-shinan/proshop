package models

import "time"

type Categories struct {
	ID           uint   `json:"id" form:"id"`
	CategoryName string `json:"categoryname" form:"category_name" binding:"required"`
}

type Products struct {
	ID           uint       `json:"id" form:"id"`
	ProductName  string     `json:"productname" form:"product_name"`
	Description  string     `json:"description" form:"description"`
	Stock        uint       `json:"stock" form:"stock"`
	Price        uint       `json:"price" form:"price"`
	Category     Categories `json:"category"`
	CategoryId   uint       `json:"category_id" form:"category_id"`
	Image        string     `json:"image"`
	IsInCart     bool
	IsInWishlist bool
}

type CartItem struct {
	ProductID     uint
	ProductName   string
	ProductPrice  float64
	Quantity      int
	TotalPrice    float64
	ImageFilename string
}

type Wishlist struct {
	ID        uint `json:"id"`
	User      Users
	UserID    uint `json:"user_id"`
	Product   Products
	ProductID uint `json:"product_id"`
}

type Coupon struct {
	ID             int       `json:"id"`
	Code           string    `json:"code" form:"code"`                       // Coupon Code
	DiscountAmount float64   `json:"discount_amount" form:"discount_amount"` // Discount Amount
	Description    string    `json:"description" form:"description"`         // Description
	UserLimit      int       `json:"user_limit" form:"user_limit"`           // Limit for an individual user
	CreatedAt      time.Time `json:"created_at"`                             // Creation timestamp
	UpdatedAt      time.Time `json:"updated_at"`                             // Last update timestamp
}
