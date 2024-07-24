package models

type Admins struct {
	ID          uint   `json:"id"`
	Firstname   string `json:"first_name"`
	Lastname    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber int64  `json:"phone"`
	IsAdmin     bool   `JSON:"isadmin"`
}
