package models

// User represents the user data structure for JSON request/response
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
