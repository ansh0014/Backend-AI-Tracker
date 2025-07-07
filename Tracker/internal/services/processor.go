package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// EventProcessor handles real-time event processing
type EventProcessor struct {
	events    map[string][]UserEvent
	mutex     sync.RWMutex
	batchSize int
	analyzer  *ActivityAnalyzer
}

func (p *EventProcessor) ProcessBatchEvents(ctx *gin.Context, userID string, events []UserEvent) (*ActivityAnalysis, error) {
	if len(events) == 0 {
		return nil, errors.New("no events to process")
	}

	// Group events by type
	eventsByType := make(map[string][]UserEvent)
	for _, event := range events {
		eventsByType[event.Type] = append(eventsByType[event.Type], event)
	}

	// Analyze patterns
	analysis, err := p.analyzer.AnalyzeActivity(ctx.Request.Context(), userID, events)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %v", err)
	}

	return analysis, nil
}

// NewEventProcessor creates a new event processor instance
func NewEventProcessor(analyzer *ActivityAnalyzer) *EventProcessor {
	return &EventProcessor{
		events:    make(map[string][]UserEvent),
		batchSize: 100,
		analyzer:  analyzer,
	}
}

// ProcessEvent handles a new incoming event
func (p *EventProcessor) ProcessEvent(ctx context.Context, userID string, event UserEvent) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Add event to user's event list
	if _, exists := p.events[userID]; !exists {
		p.events[userID] = make([]UserEvent, 0, p.batchSize)
	}

	p.events[userID] = append(p.events[userID], event)

	// Process batch if threshold reached
	if len(p.events[userID]) >= p.batchSize {
		return p.processBatch(ctx, userID)
	}

	return nil
}

// processBatch handles a batch of events for analysis
func (p *EventProcessor) processBatch(ctx context.Context, userID string) error {
	events := p.events[userID]

	// Clear events before processing to prevent duplicates
	p.events[userID] = make([]UserEvent, 0, p.batchSize)

	// Analyze the batch
	_, err := p.analyzer.AnalyzeActivity(ctx, userID, events)
	if err != nil {
		return err
	}

	// TODO: Store analysis results in database

	return nil
}

// GetRecentEvents retrieves recent events for a user
func (p *EventProcessor) GetRecentEvents(userID string, duration time.Duration) []UserEvent {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if events, exists := p.events[userID]; exists {
		cutoff := time.Now().Add(-duration)
		recent := make([]UserEvent, 0)

		for _, event := range events {
			if event.Timestamp.After(cutoff) {
				recent = append(recent, event)
			}
		}

		return recent
	}

	return nil
}
