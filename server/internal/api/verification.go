package api

import (
    "encoding/json"
    "net/http"
	
    "ppba_project/pkg/db"
)

type VerificationRequest struct {
    Username       string `json:"username"`
    BiometricProof string `json:"biometric_proof"`
}

func VerificationHandler(w http.ResponseWriter, r *http.Request) {
    var req VerificationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Convert the provided biometric proof to a byte slice (assuming it's a string)
    biometricProof := []byte(req.BiometricProof)

    // Call the database function to verify the biometric proof
    matched, err := db.VerifyBiometric(req.Username, biometricProof)
    if err != nil {
        http.Error(w, "User not found or verification error", http.StatusNotFound)
        return
    }

    if matched {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Verification successful"})
    } else {
        http.Error(w, "Verification failed", http.StatusUnauthorized)
    }
}