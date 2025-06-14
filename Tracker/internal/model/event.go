package model

import (
    "encoding/json"
    "time"
)

// EventMetadata represents specific metadata types for different events
type EventMetadata struct {
    URL        string  `json:"url,omitempty"`
    X          float64 `json:"x,omitempty"`
    Y          float64 `json:"y,omitempty"`
    ScrollDelta float64 `json:"scrollDelta,omitempty"`
    KeyCode     string  `json:"keyCode,omitempty"`
    TabID      string  `json:"tabId,omitempty"`
    PageTitle   string  `json:"pageTitle,omitempty"`
}

// NewEvent creates a new event with the current timestamp
func NewEvent(userID, eventType string, metadata map[string]interface{}) *Event {
    return &Event{
        UserID:    userID,
        Type:      eventType,
        Timestamp: time.Now(),
        Metadata:  metadata,
    }
}

// GetMetadata converts the generic metadata map to strongly typed EventMetadata
func (e *Event) GetMetadata() (*EventMetadata, error) {
    var metadata EventMetadata
    
    // Convert metadata map to JSON bytes
    jsonBytes, err := json.Marshal(e.Metadata)
    if err != nil {
        return nil, err
    }
    
    // Unmarshal into strongly typed struct
    if err := json.Unmarshal(jsonBytes, &metadata); err != nil {
        return nil, err
    }
    
    return &metadata, nil
}