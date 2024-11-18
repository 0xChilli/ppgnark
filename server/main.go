package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"ppba_project/pkg/db"
	"ppba_project/server/cache"
	"ppba_project/server/internal/api"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Initialize Redis cache with error handling
	if err := cache.InitCache("localhost:6379", "", 0); err != nil {
		log.Fatalf("Error initializing cache: %v", err)
	}

	log.Println("Cache is running...")
	// Additional server setup code here
	// Test setting and getting a cache value
    testKey := "test-key"
    testValue := "Hello, Redis!"

    // Set value in Redis
    if err := cache.SetCacheValue(testKey, testValue); err != nil {
        log.Fatalf("Failed to set cache value: %v", err)
    }

    // Retrieve value from Redis
    value, err := cache.GetCacheValue(testKey)
    if err != nil {
        log.Fatalf("Failed to get cache value: %v", err)
    }

    log.Printf("Retrieved value from cache: %s", value)

	
	
	// Initialize the database connection
	connString := os.Getenv("DB_CONN_STRING")
	if err := db.InitDB(connString); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.CloseDB()

	// Start the HTTP server and the API handlers
	http.HandleFunc("/enroll", api.EnrollmentHandler)
	http.HandleFunc("/verify", api.VerificationHandler)

	log.Println("Server started on: 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
