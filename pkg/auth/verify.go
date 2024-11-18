package auth

// compares the coming string to the already saved string
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type VerifyRequest struct {
	UserID string `json:"user_id"`
	Proof  string `json:"proof"`
}

func RunVerification() {
	// Placeholder: Collect user ID and proof from CLI args
	if len(os.Args) < 4 {
		fmt.Println("Usage: main verify <user_id> <proof>")
		os.Exit(1)
	}
	userID := os.Args[2]
	proof := os.Args[3]

	// Prepare verification data
	requestData := VerifyRequest{
		UserID: userID,
		Proof:  proof,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error encoding verification data:", err)
		os.Exit(1)
	}

	// Send request to server
	resp, err := http.Post("http://localhost:8080/verify", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending verification request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Verification failed with status:", resp.Status)
		os.Exit(1)
	}

	fmt.Println("Verification successful!")
}

// package api

// import (
//     "encoding/json"
//     "net/http"
//     "biometric-auth-project/db"
// )

// type VerificationRequest struct {
//     Username       string `json:"username"`
//     BiometricProof string `json:"biometric_proof"`
// }

// func VerificationHandler(w http.ResponseWriter, r *http.Request) {
//     var req VerificationRequest
//     if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//         http.Error(w, "Invalid request payload", http.StatusBadRequest)
//         return
//     }

//     if matched, err := db.VerifyBiometric(req.Username, req.BiometricProof); err != nil {
//         http.Error(w, "Verification failed", http.StatusInternalServerError)
//     } else if !matched {
//         http.Error(w, "Verification unsuccessful", http.StatusUnauthorized)
//     } else {
//         w.WriteHeader(http.StatusOK)
//         json.NewEncoder(w).Encode(map[string]string{"message": "Verification successful"})
//     }
// }
