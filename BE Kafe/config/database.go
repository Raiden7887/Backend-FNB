package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Global database connection
var DB *sql.DB

// ConnectDB establishes a connection to the MySQL database
func ConnectDB() (*sql.DB, error) {
	database, err := sql.Open("mysql", "root:26j13a05@tcp(localhost:3306)/db_kafe")
	if err != nil {
		return nil, err
	}
	if err := database.Ping(); err != nil {
		return nil, err
	}
	return database, nil
}

// InitDB initializes the database connection
func InitDB() {
	var err error
	DB, err = ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")
}
