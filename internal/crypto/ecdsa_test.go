package crypto

import (
	"testing"
)

func TestECDSAFunctions(t *testing.T) {
	keyFilename := "ecdsa_test_key.pem"
	message := "Hello, World!"

	// Generate ECDSA key
	if err := GenerateECDSAKey(keyFilename); err != nil {
		t.Fatalf("Error generating ECDSA key: %v", err)
	}

	// Sign the message
	signature, err := ReadPEMKeyAndSign(keyFilename, message)
	if err != nil {
		t.Fatalf("Error signing message: %v", err)
	}

	// Verify the signature
	valid, err := ReadPEMKeyAndVerify(keyFilename, message, signature)
	if err != nil {
		t.Fatalf("Error verifying signature: %v", err)
	}

	if !valid {
		t.Fatalf("Expected signature to be valid")
	}
}
