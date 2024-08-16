package models

import "time"

type PaymentRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Receipt  string `json:"receipt"`
}

type PaymentResponse struct {
	ID       string `json:"id"`
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"`
	OrderID  string `json:"order_id"`
	Status   string `json:"status"`
}
type Order struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	RazorpayOrderID string    `json:"razorpay_order_id,omitempty"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	PaymentMethod   string    `json:"payment_method"`
	AddressID       uint      `json:"address_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID        uint    `json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Amount    float64 `json:"amount"`
}

// Payment represents the payment entity
type Payment struct {
	ID                uint      `json:"id"`
	OrderID           uint      `json:"order_id"`
	RazorpayPaymentID string    `json:"razorpay_payment_id"`
	Status            string    `json:"status"`
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency"`
	CreatedAt         time.Time `json:"created_at"`
}
