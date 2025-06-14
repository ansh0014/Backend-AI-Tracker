package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"Tracker/internal/database"
	"Tracker/internal/model"
	"Tracker/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ActivityController handles activity-related operations
type ActivityController struct {
	aiService    *services.AIService
	eventService *services.EventProcessor
}

// NewActivityController creates a new activity controller
func NewActivityController() (*ActivityController, error) {
	aiService, err := services.NewAIService(30 * time.Second)
	if err != nil {
		return nil, err
	}

	analyzer := services.NewActivityAnalyzer(aiService)
	eventProcessor := services.NewEventProcessor(analyzer)

	return &ActivityController{
		aiService:    aiService,
		eventService: eventProcessor,
	}, nil
}

// AnalyzeActivity analyzes user activity patterns
func (c *ActivityController) AnalyzeActivity(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Get time range from query params, default to last 5 minutes
	endTime := time.Now()
	startTime := endTime.Add(-5 * time.Minute)

	if timeRange := ctx.Query("timeRange"); timeRange != "" {
		duration, err := time.ParseDuration(timeRange)
		if err == nil {
			startTime = endTime.Add(-duration)
		}
	}

	// Get recent events
	events := c.eventService.GetRecentEvents(userID, endTime.Sub(startTime))

	// Perform analysis
	analysis, err := c.eventService.ProcessBatchEvents(ctx, userID, events)
	if err != nil {
		HandleError(ctx, fmt.Errorf("failed to analyze activity: %v", err))
		return
	}

	// Return successful response
	ctx.JSON(http.StatusOK, gin.H{
		"analysis": analysis,
		"timeFrame": gin.H{
			"start": startTime,
			"end":   endTime,
		},
	})
}

// GetActivitySummary retrieves activity summary for a user
func (c *ActivityController) GetActivitySummary(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Set up options for aggregation
	opts := options.Aggregate().SetMaxTime(2 * time.Second)

	// Create pipeline for activity summary
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"user_id": userID,
				"created_at": bson.M{
					"$gte": time.Now().Add(-24 * time.Hour),
				},
			},
		},
		{
			"$group": bson.M{
				"_id":             nil,
				"totalActivities": bson.M{"$sum": 1},
				"totalDuration":   bson.M{"$sum": "$duration"},
				"categories": bson.M{
					"$addToSet": "$category",
				},
			},
		},
	}

	cursor, err := database.GetCollection().Aggregate(ctx, pipeline, opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to aggregate activities"})
		return
	}
	defer cursor.Close(ctx)

	var summaries []bson.M
	if err = cursor.All(ctx, &summaries); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode summary"})
		return
	}

	if len(summaries) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No activities found"})
		return
	}

	ctx.JSON(http.StatusOK, summaries[0])
}

// CreateActivity creates a new activity
func (c *ActivityController) CreateActivity(ctx *gin.Context) {
	var req model.ActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
			"fields": map[string]string{
				"title":    "required",
				"category": "required",
				"duration": "required, must be greater than 0",
				"date":     "required, format: YYYY-MM-DD",
			},
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Preferences parameter is required",
			"details": map[string]string{
				"preferences": "Query parameter 'preferences' must be provided",
				"example":     "/api/suggestions?preferences=productivity,focus",
			},
		})
		return
	}

	// Create context with timeout
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	suggestions, err := c.aiService.GetActivitySuggestions(reqCtx, preferences)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	if len(suggestions) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"suggestions": []string{},
			"message":     "No suggestions found for given preferences",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"suggestions": suggestions,
		"count":       len(suggestions),
		"preferences": preferences,
	})
}

// HandleError is a helper function to handle common errors
func HandleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
	case errors.Is(err, context.DeadlineExceeded):
		ctx.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timeout"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// Close safely closes the controller's resources
func (c *ActivityController) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.aiService.Close(ctx); err != nil {
		return err
	}
	return nil
}
