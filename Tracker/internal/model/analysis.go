package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AnalysisResult represents the detailed analysis results
type AnalysisResult struct {
	BehaviorSummary string          `bson:"behaviorSummary" json:"behaviorSummary"`
	ActivityMetrics ActivityMetrics `bson:"activityMetrics" json:"activityMetrics"`
	Recommendations []string        `bson:"recommendations" json:"recommendations"`
	TimeFrame       TimeFrame       `bson:"timeFrame" json:"timeFrame"`
}

// ActivityMetrics contains statistical metrics about the activity
type ActivityMetrics struct {
	TotalEvents      int     `bson:"totalEvents" json:"totalEvents"`
	ActiveTime       float64 `bson:"activeTime" json:"activeTime"`
	IdleTime         float64 `bson:"idleTime" json:"idleTime"`
	FocusScore       float64 `bson:"focusScore" json:"focusScore"`
	ProductivityRate float64 `bson:"productivityRate" json:"productivityRate"`
}

// TimeFrame represents a period of analysis
type TimeFrame struct {
	Start time.Time `bson:"start" json:"start"`
	End   time.Time `bson:"end" json:"end"`
}

// NewAnalysis creates a new analysis instance
func NewAnalysis(userID string, activityID primitive.ObjectID) *Analysis {
	return &Analysis{
		UserID:     userID,
		ActivityID: activityID,
		CreatedAt:  time.Now(),
		Tags:       make([]string, 0),
	}
}

// SetBehavior updates the behavior type and confidence
func (a *Analysis) SetBehavior(behaviorType string, confidence float64) {
	a.BehaviorType = behaviorType
	a.Confidence = confidence
}

// AddTags adds new tags to the analysis
func (a *Analysis) AddTags(tags ...string) {
	a.Tags = append(a.Tags, tags...)
}

// IsValid checks if the analysis has valid required fields
func (a *Analysis) IsValid() bool {
	return a.UserID != "" &&
		a.ActivityID != primitive.NilObjectID &&
		a.BehaviorType != "" &&
		a.Confidence > 0 &&
		a.Confidence <= 1.0
}
