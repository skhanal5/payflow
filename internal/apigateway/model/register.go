package model

type UserRegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Address  string `json:"address"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
	ZipCode  string `json:"zip_code"`
}