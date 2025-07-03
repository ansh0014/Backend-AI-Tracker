package main

import (
	"log"
	"net/http"
	"os"

	"Tracker/internal/config"
	"Tracker/internal/database"
	ws "Tracker/internal/ws"
	routes "Tracker/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize WebSocket manager and start it
	manager := ws.NewManager()
	go manager.Run()

	// Initialize router
	router := routes.SetupRouter(manager)

	// Add WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		// Use the already created manager
		handler := ws.NewHandler(manager)
		handler.HandleWebSocket(c.Writer, c.Request)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
