package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

// Response represents a standard API response structure
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// RespondWithJSON sends a JSON response with the given status code and payload
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error marshalling response"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// RespondWithError sends a JSON error response
func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, Response{
		Status:  status,
		Message: message,
	})
}

// GetCurrentTimestamp returns the current UTC timestamp
func GetCurrentTimestamp() time.Time {
	return time.Now().UTC()
}

// FormatDate formats a time.Time to a string in ISO format
func FormatDate(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseDate parses a date string in ISO format to time.Time
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, dateStr)
}

// IsValidDate checks if a date string is in valid ISO format
func IsValidDate(dateStr string) bool {
	_, err := ParseDate(dateStr)
	return err == nil
}
