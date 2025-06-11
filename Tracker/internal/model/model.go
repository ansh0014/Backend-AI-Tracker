package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Activity represents a user activity
type Activity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"`
	Duration    int                `bson:"duration" json:"duration"` // in minutes
	Date        time.Time          `bson:"date" json:"date"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// ActivityRequest represents the request body for creating/updating an activity
type ActivityRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"required"`
	Duration    int    `json:"duration" binding:"required"`
	Date        string `json:"date" binding:"required"`
}
