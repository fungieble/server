package models

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Username    string `json:"username"`
	Description string `json:"description"`
}
