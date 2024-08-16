package models

type Address struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	Name         string `json:"name"`
	MobileNumber string `json:"mobile_number"`
	HomeName     string `json:"home_name"`
	Place        string `json:"place"`
	Landmark     string `json:"landmark"`
	City         string `json:"city"`
	District     string `json:"district"`
	State        string `json:"state"`
	Country      string `json:"country"`
	PostalCode   string `json:"postal_code"`
}
