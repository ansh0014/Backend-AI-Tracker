package services

import (
	"math"
	"time"
)

// UserEvent represents a single user activity event
type UserEvent struct {
	Type      string        `json:"type"`
	Timestamp time.Time     `json:"timestamp"`
	Metadata  EventMetadata `json:"metadata"`
}

// EventMetadata contains event-specific details
type EventMetadata struct {
	URL        string  `json:"url,omitempty"`
	X          float64 `json:"x,omitempty"`
	Y          float64 `json:"y,omitempty"`
	KeyPressed string  `json:"keyPressed,omitempty"`
	TabID      string  `json:"tabId,omitempty"`
}

// classifyBehavior analyzes events to determine user behavior
func classifyBehavior(events []UserEvent) (string, float64) {
	if len(events) == 0 {
		return "idle", 1.0
	}

	// Calculate metrics
	metrics := calculateMetrics(events)

	// Determine behavior based on metrics
	switch {
	case metrics.idleTime > 0.7:
		return "idle", metrics.idleTime
	case metrics.tabSwitches > 10 && metrics.activeTime > 0.8:
		return "multitasking", metrics.activeTime
	case metrics.focusedTime > 0.6 && metrics.tabSwitches < 5:
		return "focused", metrics.focusedTime
	default:
		return "distracted", math.Max(metrics.activeTime, 0.3)
	}
}

// behaviorMetrics contains analyzed behavior data
type behaviorMetrics struct {
	activeTime  float64
	idleTime    float64
	focusedTime float64
	tabSwitches int
}

// calculateMetrics processes events to extract behavior metrics
func calculateMetrics(events []UserEvent) behaviorMetrics {
	metrics := behaviorMetrics{}

	// Implementation details for calculating metrics
	// This would analyze event patterns, timing, and sequences

	return metrics
}
