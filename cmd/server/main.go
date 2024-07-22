package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/korentmaj/go-ecdsa-status-netis-challenge/internal/api"
	"github.com/korentmaj/go-ecdsa-status-netis-challenge/internal/crypto"
	"github.com/korentmaj/go-ecdsa-status-netis-challenge/internal/database"
)

func main() {
	// Initialize the database
	if err := database.InitDB("postgres://username:password@localhost/dbname?sslmode=disable"); err != nil {
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
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvdK1xzZhN8/Yy6nl2+ld
...
...
...
H5BtuG4IoR3Ib+KHN0BFtF3GCcPx5+/OZ9xEkmYbljB1OwIDAQAB
-----END PUBLIC KEY-----
`

	publicKey, err := crypto.ParseECDSAPublicKeyFromPEM(pemEncodedPublicKey)
	if err != nil {
		log.Fatalf("Error parsing public key: %v", err)
	}

	status, err := crypto.GetStatusFromJWS(url, publicKey)
	if err != nil {
		log.Fatalf("Error getting status from JWS: %v", err)
	}

	fmt.Printf("Status: %v\n", status)
}
