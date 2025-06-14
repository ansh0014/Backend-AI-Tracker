package services

import (
	"context"
	"time"
)

// ActivityAnalysis represents the analysis results
type ActivityAnalysis struct {
	UserID          string    `json:"userId"`
	TimeFrame       TimeFrame `json:"timeFrame"`
	Behavior        string    `json:"behavior"`
	Confidence      float64   `json:"confidence"`
	Recommendations []string  `json:"recommendations"`
	AnalyzedAt      time.Time `json:"analyzedAt"`
}

// TimeFrame represents the analysis period
type TimeFrame struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// ActivityAnalyzer handles activity analysis
type ActivityAnalyzer struct {
	aiService *AIService
}

// NewActivityAnalyzer creates a new analyzer instance
func NewActivityAnalyzer(aiService *AIService) *ActivityAnalyzer {
	return &ActivityAnalyzer{
		aiService: aiService,
	}
}

// AnalyzeActivity analyzes user activity for a given timeframe
func (a *ActivityAnalyzer) AnalyzeActivity(ctx context.Context, userID string, events []UserEvent) (*ActivityAnalysis, error) {
	behavior, confidence := classifyBehavior(events)

	analysis := &ActivityAnalysis{
		UserID: userID,
		TimeFrame: TimeFrame{
			Start: time.Now().Add(-5 * time.Minute),
			End:   time.Now(),
		},
		Behavior:   behavior,
		Confidence: confidence,
		AnalyzedAt: time.Now(),
	}

	// Get AI recommendations based on behavior
	recommendations, err := a.aiService.GetActivitySuggestions(ctx, behavior)
	if err == nil {
		analysis.Recommendations = recommendations
	}

	return analysis, nil
}
