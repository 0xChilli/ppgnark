package circuit

import (
	"fmt"
	"bytes"
	"github.com/consensys/gnark-crypto/ecc"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	r1cs2 "github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/mimc"
	"math/big"
	_"reflect"
)

type Circuit struct {
	BiometricExtract frontend.Variable `gnark:",private"`
	Proof     frontend.Variable `gnark:",public"`
}

func (circuit *Circuit) Define(api frontend.API) error {
	mimc, _ := mimc.NewMiMC(api)
	mimc.Write(circuit.BiometricExtract)
	computedHash := mimc.Sum()
	api.AssertIsEqual(circuit.Proof, computedHash)
	return nil
}

// GenerateProof computes the zk-SNARK proof for a given PreImage
func GenerateProof(BiometricExtract *big.Int) ([]byte, error) {
	// Compute the hash of the pre-image using external MiMC
	hash := mimcHash(BiometricExtract)

	// Define and compile the circuit
	var circuit Circuit
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs2.NewBuilder, &circuit)
	if err != nil {
		return nil, fmt.Errorf("failed to compile circuit: %v", err)
	}

	// Setup keys
	pk, _, err := groth16.Setup(r1cs)
	if err != nil {
		return nil, fmt.Errorf("failed to setup Groth16: %v", err)
	}

	// Assignment for the circuit
	assignment := &Circuit{
		BiometricExtract: frontend.Variable(BiometricExtract),
		Proof:     frontend.Variable(hash),
	}

	// Create a witness
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("failed to create witness: %v", err)
	}

	// Generate the proof
	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		return nil, fmt.Errorf("failed to generate proof: %v", err)
	}
	// Serialize the proof
    var buf bytes.Buffer
    if _, err := proof.WriteTo(&buf); err != nil {
        return nil, fmt.Errorf("failed to serialize proof: %v", err)
    }
    proofBytes := buf.Bytes()

    return proofBytes, nil
}

// mimcHash computes the MiMC hash of the input data
func mimcHash(data *big.Int) string {
	f := bn254.NewMiMC()
	f.Write(data.Bytes())
	hash := f.Sum(nil)
	hashInt := new(big.Int).SetBytes(hash)
	return hashInt.String()
}
