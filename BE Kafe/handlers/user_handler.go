package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"be_kafe/config"
	"be_kafe/models"
	"be_kafe/utils"

	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body into User struct
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Check if username already exists
	rows, err := config.DB.Query("SELECT username FROM users WHERE username = ?", user.Username)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if rows.Next() {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	_, err = config.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, hashedPassword)
	if err != nil {
		log.Printf("Database insert error: %v", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body into User struct
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Retrieve hashed password from database
	var storedHash string
	err := config.DB.QueryRow("SELECT password FROM users WHERE username = ?", req.Username).Scan(&storedHash)
	if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(req.Username)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send success response with token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"user": map[string]string{
			"username": req.Username,
		},
	})
}

// GetUserHandler retrieves user information
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get username from query parameter
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Query user from database
	var user models.User
	err := config.DB.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&user.Username)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Don't include password in response
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
