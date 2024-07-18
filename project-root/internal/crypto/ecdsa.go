package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
)

// GenerateECDSAKey creates an ECDSA P-256 key and writes it in PEM format to a specified file.
func GenerateECDSAKey(filename string) error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate ECDSA key: %v", err)
	}

	der, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal ECDSA private key: %v", err)
	}

	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: der,
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	if err := pem.Encode(file, block); err != nil {
		return fmt.Errorf("failed to write PEM to file: %v", err)
	}

	return nil
}

// ReadPEMKeyAndSign reads the PEM key from the specified file, signs the message, and returns the signature in Base64URL format.
func ReadPEMKeyAndSign(filename, message string) (string, error) {
	pemData, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read PEM file: %v", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return "", fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse ECDSA private key: %v", err)
	}

	hash := sha256.Sum256([]byte(message))

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}

	signature := append(r.Bytes(), s.Bytes()...)
	base64URLSignature := base64.URLEncoding.EncodeToString(signature)

	return base64URLSignature, nil
}

// ReadPEMKeyAndVerify reads the PEM key from the specified file and verifies the signature of the message.
func ReadPEMKeyAndVerify(filename, message, base64URLSignature string) (bool, error) {
	pemData, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, fmt.Errorf("failed to read PEM file: %v", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return false, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse ECDSA private key: %v", err)
	}

	publicKey := &privateKey.PublicKey

	signature, err := base64.URLEncoding.DecodeString(base64URLSignature)
	if err != nil {
		return false, fmt.Errorf("failed to decode Base64URL signature: %v", err)
	}

	if len(signature) != 64 {
		return false, fmt.Errorf("invalid signature length")
	}

	r := big.NewInt(0).SetBytes(signature[:32])
	s := big.NewInt(0).SetBytes(signature[32:])

	hash := sha256.Sum256([]byte(message))

	valid := ecdsa.Verify(publicKey, hash[:], r, s)

	return valid, nil
}
