package ws

// EventType constants for WebSocket events
const (
	EventTypeActivity = "activity"
	EventTypeMessage  = "message"
	EventTypeAlert    = "alert"
)

type WebSocketEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

