package services

import (
	"context"
	"sync"
	"time"
)

// EventProcessor handles real-time event processing
type EventProcessor struct {
	events    map[string][]UserEvent
	mutex     sync.RWMutex
	batchSize int
	analyzer  *ActivityAnalyzer
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
