package main

import (
	"log"
	"net/http"
	"os"

	"Tracker/internal/config"
	routes "Tracker/router" // Changed from "Tracker/router" to "Tracker/routes"
)

func main() {
	// Load configuration
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize router
	router := routes.SetupRouter() // Removed .Routes() call

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
