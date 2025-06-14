package model

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Username  string            `bson:"username" json:"username"`
    Email     string            `bson:"email" json:"email"`
    Settings  UserSettings      `bson:"settings" json:"settings"`
    CreatedAt time.Time         `bson:"createdAt" json:"createdAt"`
    UpdatedAt time.Time         `bson:"updatedAt" json:"updatedAt"`
}

// UserSettings contains user-specific settings
type UserSettings struct {
    TrackingEnabled bool              `bson:"trackingEnabled" json:"trackingEnabled"`
    Preferences     map[string]string `bson:"preferences" json:"preferences"`
    NotifyEmail     bool              `bson:"notifyEmail" json:"notifyEmail"`
    TimeZone        string            `bson:"timeZone" json:"timeZone"`
}