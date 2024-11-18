package api

import (
	"encoding/json"
	"math/big"
	"net/http"
	"ppba_project/pkg/circuit"
	"ppba_project/pkg/db"
)

type EnrollmentRequest struct {
	Username       string `json:"username"`
	BiometricProof string `json:"biometric_proof"`
}

func EnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	var req EnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Convert the biometric proof to a big.Int
	biometricProof := new(big.Int)
	_, ok := biometricProof.SetString(req.BiometricProof, 10) // Base 10 for decimal strings
	if !ok {
		http.Error(w, "Invalid biometric proof format", http.StatusBadRequest)
		return
	}

	// Call the GenerateProof function
	proof, err := circuit.GenerateProof(biometricProof)
	if err != nil {
		http.Error(w, "Failed to generate proof", http.StatusInternalServerError)
		return
	}

	// Save the proof to the database
	// err = db.SaveProof(proof)
	// if err != nil {
	// 	http.Error(w, "Failed to save proof to database", http.StatusInternalServerError)
	// 	return
	// }
	err = db.SaveProof(req.Username,proof)
	if err != nil {
		http.Error(w, "Failed to save proof to database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Enrollment successful"})
}