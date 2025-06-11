package controllers

import (
	"context"
	"net/http"
	"time"

	"Tracker/internal/database"
	"Tracker/internal/model"
	"Tracker/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActivityController handles activity-related operations
type ActivityController struct {
	aiService *services.AIService
}

// NewActivityController creates a new activity controller
func NewActivityController() (*ActivityController, error) {
	aiService, err := services.NewAIService()
	if err != nil {
		return nil, err
	}
	return &ActivityController{
		aiService: aiService,
	}, nil
}

// CreateActivity creates a new activity
func (c *ActivityController) CreateActivity(ctx *gin.Context) {
	var req model.ActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	activity := model.Activity{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Duration:    req.Duration,
		Date:        date,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := database.GetCollection().InsertOne(context.Background(), activity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	activity.ID = result.InsertedID.(primitive.ObjectID)
	ctx.JSON(http.StatusCreated, activity)
}

// GetActivities retrieves all activities
func (c *ActivityController) GetActivities(ctx *gin.Context) {
	cursor, err := database.GetCollection().Find(context.Background(), bson.M{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var activities []model.Activity
	if err := cursor.All(context.Background(), &activities); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, activities)
}

// GetActivity retrieves a specific activity
func (c *ActivityController) GetActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var activity model.Activity
	err = database.GetCollection().FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&activity)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	ctx.JSON(http.StatusOK, activity)
}

// UpdateActivity updates an existing activity
func (c *ActivityController) UpdateActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req model.ActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":       req.Title,
			"description": req.Description,
			"category":    req.Category,
			"duration":    req.Duration,
			"date":        date,
			"updated_at":  time.Now(),
		},
	}

	result, err := database.GetCollection().UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		update,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Activity updated successfully"})
}

// DeleteActivity deletes an activity
func (c *ActivityController) DeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := database.GetCollection().DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}

// GetSuggestions generates activity suggestions
func (c *ActivityController) GetSuggestions(ctx *gin.Context) {
	preferences := ctx.Query("preferences")
	if preferences == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Preferences are required"})
		return
	}

	suggestions, err := c.aiService.GetActivitySuggestions(preferences)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}
