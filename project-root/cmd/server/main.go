package main

import (
	"log"
	"net/http"
	"project-root/internal/api"
	"project-root/internal/database"
)

func main() {
	// Initialize the database
	if err := database.InitDB("postgres://username:password@localhost/dbname?sslmode=disable"); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set up the router and start the server
	router := api.SetupRouter()
	log.Fatal(http.ListenAndServe(":8000", router))
}
