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
	"os"
	"os/signal"
	"syscall"

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
	// PostgreSQL podatki
	dbUser := "user"
	dbPassword := "pass"
	dbName := "imebaze"
	dbHost := "localhost"
	dbPort := "5432"

	// PostgreSQL connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Inicializacija baze
	if err := database.InitDB(connStr); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB()

	// Določi port
	router := api.SetupRouter()
	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on :8000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	if err := server.Close(); err != nil {
		log.Fatalf("Server Close(): %v", err)
	}
	log.Println("Server gracefully stopped")

	// URL za: GetStatusFromJWS
	url := "http://localhost:8000/api/status/testStatusId?index=1"
	pemEncodedPublicKey := `
-----BEGIN PUBLIC KEY-----
//prilagodi ključ
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

	req.SetBasicAuth("user", "pass") // replace with actual credentials

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
