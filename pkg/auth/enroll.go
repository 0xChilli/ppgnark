package auth

//makes zk-snark proof and saves it to the database.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	_ "ppba_project/pkg/db"
)

type EnrollRequest struct {
	UserID        string `json:"user_id"`
	BiometricData string `json:"biometric_data"`
}

func RunEnrollment() {
	// Placeholder: Collect user ID and biometric data from CLI args
	if len(os.Args) < 4 {
		fmt.Println("Usage: main enroll <user_id> <biometric_data>")
		os.Exit(1)
	}
	userID := os.Args[2]
	biometricData := os.Args[3]

	// Prepare enrollment data
	requestData := EnrollRequest{
		UserID:        userID,
		BiometricData: biometricData,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error encoding enrollment data:", err)
		os.Exit(1)
	}

	// Send request to server
	resp, err := http.Post("http://localhost:8080/enroll", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending enrollment request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Enrollment failed with status:", resp.Status)
		os.Exit(1)
	}

	fmt.Println("Enrollment successful!")
}
