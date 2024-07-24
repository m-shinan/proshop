package models

import "time"

type Users struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber int    `json:"phone"`
	IsAdmin     bool   `json:"isadmin"`
	
	Isblocked   bool   `json:"isblocked"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
// Otp         string `json:"otp"`