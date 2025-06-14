package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Event represents a user activity event
type Event struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID    string                 `bson:"userId" json:"userId"`
	Type      string                 `bson:"type" json:"type"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
	Metadata  map[string]interface{} `bson:"metadata" json:"metadata"`
}

// Activity represents a collection of events with AI analysis
type Activity struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"userId" json:"userId"`
	StartTime time.Time          `bson:"startTime" json:"startTime"`
	EndTime   time.Time          `bson:"endTime" json:"endTime"`
	Events    []Event            `bson:"events" json:"events"`
	Analysis  *Analysis          `bson:"analysis" json:"analysis"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// Analysis represents AI-generated analysis of user activity
type Analysis struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       string             `bson:"userId" json:"userId"`
	ActivityID   primitive.ObjectID `bson:"activityId" json:"activityId"`
	BehaviorType string             `bson:"behaviorType" json:"behaviorType"`
	Confidence   float64            `bson:"confidence" json:"confidence"`
	Summary      string             `bson:"summary" json:"summary"`
	Tags         []string           `bson:"tags" json:"tags"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}

// ActivityRequest represents the incoming request for activity operations
type ActivityRequest struct {
	Title       string  `bson:"title" binding:"required"`
	Description string  `bson:"description"`
	Category    string  `bson:"category" binding:"required"`
	Duration    float64 `bson:"duration" binding:"required,gt=0"`
	Date        string  `bson:"date" binding:"required"`
}

// BehaviorType constants
const (
	BehaviorFocused      = "focused"
	BehaviorIdle         = "idle"
	BehaviorMultitasking = "multitasking"
	BehaviorDistracted   = "distracted"
)

// EventType constants
const (
	EventMouseMove  = "mouse_move"
	EventClick      = "click"
	EventKeyPress   = "key_press"
	EventScroll     = "scroll"
	EventTabFocus   = "tab_focus"
	EventTabBlur    = "tab_blur"
	EventPageLoad   = "page_load"
	EventPageUnload = "page_unload"
)
