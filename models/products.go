package models

type Categories struct {
	ID           uint   `json:"id" form:"id"`
	CategoryName string `json:"categoryname" form:"category_name" binding:"required"`
}

type Products struct {
	ID          uint       `json:"id" form:"id"`
	ProductName string     `json:"productname" form:"product_name"`
	Description string     `json:"description" form:"description"`
	Stock       uint       `json:"stock" form:"stock"`
	Price       uint       `json:"price" form:"price"`
	Category    Categories `json:"category"`
	CategoryId  uint       `json:"category_id" form:"category_id"`
	Image       string     `json:"image"`
}

type CartItem struct {
	ProductID    uint
	ProductName  string
	ProductPrice float64
	Quantity     int
	TotalPrice   float64
}

type Wishlist struct {
	ID        uint `json:"id"`
	User      Users
	UserID    uint `json:"user_id"`
	Product   Products
	ProductID uint `json:"product_id"`
}
