package main

import (
	"log"
	"net/http"

	"be_kafe/config"
	"be_kafe/handlers"
	"be_kafe/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database connection
	config.InitDB()
	defer config.DB.Close()

	// Initialize router
	router := mux.NewRouter()

	// Public endpoints (no authentication required)
	router.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// Protected endpoints (authentication required)
	router.HandleFunc("/user", middleware.AuthMiddleware(handlers.GetUserHandler)).Methods("GET")
	router.HandleFunc("/menu", middleware.AuthMiddleware(handlers.CreateMenuHandler)).Methods("POST")
	router.HandleFunc("/menu/{id}", middleware.AuthMiddleware(handlers.UpdateMenuHandler)).Methods("PUT")
	router.HandleFunc("/menu/{id}", middleware.AuthMiddleware(handlers.DeleteMenuHandler)).Methods("DELETE")

	// Start server
	log.Println("Starting server on port 8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
