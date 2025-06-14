package routes

import (
	"Tracker/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer ws.Close()

		// Handle WebSocket connection
		handleWebSocket(ws)
	}
}

// handleWebSocket processes WebSocket messages
func handleWebSocket(conn *websocket.Conn) {
	for {
		// Read message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Log error
			}
			break
		}

		// Process message
		if err := processWebSocketMessage(conn, messageType, message); err != nil {
			// Log error
			break
		}
	}
}

// processWebSocketMessage handles different types of WebSocket messages
func processWebSocketMessage(conn *websocket.Conn, messageType int, message []byte) error {
	// TODO: Implement message processing logic
	return conn.WriteMessage(messageType, message)
}

// HealthCheckHandler returns API health status
func HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   model.TimeFrame{}.Start,
		})
	}
}
