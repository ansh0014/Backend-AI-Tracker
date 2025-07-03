package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Tracker/internal/model"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, implement proper origin checking
	},
}

// Handler handles WebSocket connections and events
type Handler struct {
	manager *Manager
}

// NewHandler creates a new WebSocket handler
func NewHandler(manager *Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

// HandleWebSocket upgrades HTTP connection to WebSocket and handles events
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get user ID from query parameters
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Create new client
	client := NewClient(conn, userID, h.manager)
	h.manager.RegisterClient(client)

	// Start goroutines for reading and writing
	go client.ReadPump()
	go client.WritePump()
}

// Client represents a WebSocket client
type Client struct {
	conn    *websocket.Conn
	userID  string
	manager *Manager
	send    chan []byte
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, userID string, manager *Manager) *Client {
	return &Client{
		conn:    conn,
		userID:  userID,
		manager: manager,
		send:    make(chan []byte, 256),
	}
}

// ReadPump pumps messages from the WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		c.manager.UnregisterClient(c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Parse event
		var event model.Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("error parsing event: %v", err)
			continue
		}

		// Set event metadata
		event.UserID = c.userID
		event.Timestamp = time.Now()

		// Process event
		c.manager.ProcessEvent(&event)
	}
}

// WritePump pumps messages to the WebSocket connection
func (c *Client) WritePump() {
	defer func() {
		c.conn.Close()
	}()
		for {
			select {
			case message, ok := <-c.send:
				if !ok {
					// The channel was closed, so close the WebSocket connection
					_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
	
				// Send the message to the WebSocket connection
				if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
					return
				}
			}
		}
	}

