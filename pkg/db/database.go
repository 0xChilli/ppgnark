package db

// connection  , termination , insert user , delete user and most importantly get user function
import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID             int
	Username       string
	BiometricProof []byte
}

var Pool *pgxpool.Pool

// InitDB initializes a connection Pool to the PostgreSQL database.
func InitDB(connString string) error {
	var err error
	Pool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		return err
	}
	// Query to check the current database
	var currentDB string
	err = Pool.QueryRow(context.Background(), "SELECT current_database()").Scan(&currentDB)
	if err != nil {
		Pool.Close()
		return err
	}
	log.Printf("Connected to database: %s", currentDB)

	// Verify the connection
	if err := Pool.Ping(context.Background()); err != nil {
		Pool.Close()
		return err
	}
	var version string
	err = Pool.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	// Set the search path to include the public schema
	_, err = Pool.Exec(context.Background(), "SET search_path TO public")
	if err != nil {
		Pool.Close()
		return err
	}
	fmt.Printf("Connected to PostgreSQL version: %s\n", version)

	log.Println("Connected to PostgreSQL")
	return nil
}

// CloseDB closes the database connection Pool
func CloseDB() {
	if Pool != nil {
		Pool.Close()
	}
}
// func SaveProof(proof []byte) error {
// 	query := `INSERT INTO users (biometric_proof) VALUES ($1)`
// 	_, err := Pool.Exec(context.Background(), query, proof)
// 	if err != nil {
// 		return fmt.Errorf("saveProof: %v", err)
// 	}
// 	return nil
// }
// func SaveProof(proof []byte) error {
// 	query := `INSERT INTO users (biometric_proof) VALUES ($1)`
// 	log.Printf("Attempting to save proof: %v", proof) // Log proof for debugging
// 	_, err := Pool.Exec(context.Background(), query, proof)
// 	if err != nil {
// 		log.Printf("Error in SaveProof: %v", err) // Log detailed error
// 		return fmt.Errorf("saveProof: %v", err)
// 	}
// 	log.Println("Proof saved successfully")
// 	return nil
// }
func SaveProof(username string, proof []byte) error {
	query := `INSERT INTO users (username, biometric_proof) VALUES ($1, $2)`
	log.Printf("Attempting to save proof for username: %s with proof: %v", username, proof)
	_, err := Pool.Exec(context.Background(), query, username, proof)
	if err != nil {
		log.Printf("Error in SaveProof: %v", err)
		return fmt.Errorf("saveProof: %v", err)
	}
	log.Println("Proof saved successfully")
	return nil
}
// SaveEnrollment saves a new user into the database.
func SaveEnrollment(username string, biometricProof []byte) error {
	// Check if the user already exists
	var existingUser User
	err := Pool.QueryRow(context.Background(), "SELECT id FROM users WHERE username=$1", username).Scan(&existingUser.ID)
	if err == nil {
		return errors.New("user already enrolled")
	}

	// Insert the new user
	_, err = Pool.Exec(context.Background(),
		"INSERT INTO users (username, biometric_proof) VALUES ($1, $2)", username, biometricProof)
	if err != nil {
		log.Printf("Error inserting user: %v\n", err) // Log the error for debugging
		return err
	}

	return nil
}

// VerifyBiometric checks if the provided biometric proof matches the stored one for a given username.
func VerifyBiometric(username string, providedProof []byte) (bool, error) {
	// Retrieve the stored biometric proof for the user
	var storedProof []byte
	err := Pool.QueryRow(context.Background(), "SELECT biometric_proof FROM users WHERE username=$1", username).Scan(&storedProof)
	if err != nil {
		return false, errors.New("user not found or other error")
	}

	// Compare the provided proof with the stored proof
	if string(storedProof) == string(providedProof) {
		return true, nil // Match found
	}
	return false, nil // No match
}

// CreateUser adds a new user to the database
func CreateUser(ctx context.Context, username string, biometricProof []byte) error {
	query := `INSERT INTO users (username, biometric_proof) VALUES ($1, $2)`
	_, err := Pool.Exec(ctx, query, username, biometricProof)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}
func insertUser(Pool *pgxpool.Pool, user User) error {
	query := `INSERT INTO users (username, biometric_proof) VALUES ($1, $2) RETURNING id`

	// Execute the query
	err := Pool.QueryRow(context.Background(), query, user.Username, user.BiometricProof).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("insertUser: %v", err)
	}
	return nil
}

// GetUser fetches a user by username
func GetUser(ctx context.Context, username string) (User, error) {
	user := User{}
	query := `SELECT id, username, biometric_proof FROM users WHERE username=$1`
	row := Pool.QueryRow(ctx, query, username)
	err := row.Scan(&user.ID, &user.Username, &user.BiometricProof)
	if err != nil {
		return user, fmt.Errorf("failed to get user: %v", err)
	}
	return user, nil
}
