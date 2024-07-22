package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/korentmaj/go-ecdsa-status-netis-challenge/internal/api"
	"github.com/korentmaj/go-ecdsa-status-netis-challenge/internal/crypto"
	"github.com/korentmaj/go-ecdsa-status-netis-challenge/internal/database"
)

func ParseECDSAPublicKeyFromPEM(pemEncodedPubKey string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemEncodedPubKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("not ECDSA public key")
	}
}

func main() {
	// Define your PostgreSQL credentials
	dbUser := "ecdsa_user"
	dbPassword := "ecdsa_password"
	dbName := "ecdsadb"
	dbHost := "localhost"
	dbPort := "5432"

	// Construct the PostgreSQL connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Initialize the database
	if err := database.InitDB(connStr); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set up the router and start the server
	router := api.SetupRouter()
	go func() {
		log.Fatal(http.ListenAndServe(":8000", router))
	}()

	// Example usage of GetStatusFromJWS
	url := "http://localhost:8000/api/status/a_Q5JxCz#1"
	pemEncodedPublicKey := `
-----BEGIN PUBLIC KEY-----

-----END PUBLIC KEY-----
`

	publicKey, err := ParseECDSAPublicKeyFromPEM(pemEncodedPublicKey)
	if err != nil {
		log.Fatalf("Error parsing public key: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	// Add Basic Authentication header
	req.SetBasicAuth("your_username", "your_password") // Replace with actual credentials

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error getting status from JWS: received non-OK HTTP status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	status, err := crypto.ParseJWSResponse(body, publicKey)
	if err != nil {
		log.Fatalf("Error getting status from JWS: %v", err)
	}

	fmt.Printf("Status: %v\n", status)
}
