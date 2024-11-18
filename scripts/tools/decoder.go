package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

// decodeHexProof decodes a hexadecimal string into bytes and converts it into a decimal number.
func decodeHexProof(hexProof string) ([]byte, *big.Int, error) {
	decoded, err := hex.DecodeString(hexProof)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode hex proof: %w", err)
	}

	// Convert the decoded bytes to a big.Int for decimal representation
	decimalProof := new(big.Int).SetBytes(decoded)

	return decoded, decimalProof, nil
}

func main() {
	hexProof := "c1a160f1995dd9c9ed3d2f9ca5fc090e0b933928aa25b5fa37012eb657fc1cc187be74da4c2230ba38cec4a2f54fb987aa8b7c9b7538c0b3581ded34d28a36e112e1a59d0271b00d9740f6648e9f1a41ecbf733308ecdc2e40805098c1708509a81fb4e5d447993f51c2bc05e21ac806e61cbd46473e5775e98f9b9d9a53d06c000000004000000000000000000000000000000000000000000000000000000000000000"

	decodedProof, decimalProof, err := decodeHexProof(hexProof)
	if err != nil {
		log.Fatalf("Error decoding hex proof: %v", err)
	}

	fmt.Printf("Decoded Proof (bytes): %v\n", decodedProof)
	fmt.Printf("Decoded Proof (decimal): %s\n", decimalProof.String())
}