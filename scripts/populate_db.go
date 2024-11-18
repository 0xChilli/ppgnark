package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database connection string and number of records
	dbURL := os.Getenv("DB_CONN_STRING")
	if dbURL == "" {
		log.Fatalf("DB_CONN_STRING is not set in the environment variables")
	}

	numRecords := 100 // Number of records to insert

	// Connect to the database
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}
	defer dbpool.Close()

	fmt.Println("Connected to the database successfully.")

	// Populate data
	err = populateUsers(dbpool, numRecords)
	if err != nil {
		log.Fatalf("Error populating users: %v", err)
	}

	fmt.Println("Database population completed successfully!")
}

func populateUsers(pool *pgxpool.Pool, count int) error {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Prepare data insertion
	sql := `INSERT INTO testusers (name, email, created_at) VALUES ($1, $2, $3)`
	for i := 0; i < count; i++ {
		name := gofakeit.Name()
		email := gofakeit.Email()
		createdAt := gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now())

		_, err := tx.Exec(ctx, sql, name, email, createdAt)
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}