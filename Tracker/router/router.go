package routes

import (
	"Tracker/internal/controllers"
	"Tracker/internal/ws"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all the routes for the application
func SetupRouter(manager *ws.Manager) *gin.Engine {
	router := gin.Default()

	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	// Initialize controller
	activityController, err := controllers.NewActivityController()
	if err != nil {
		panic(err)
	}

	// Activity routes
	activities := router.Group("/api/activities")
	{
		activities.POST("", activityController.CreateActivity)
		activities.GET("", activityController.GetActivities)
		activities.GET("/:id", activityController.GetActivity)
		activities.PUT("/:id", activityController.UpdateActivity)
		activities.DELETE("/:id", activityController.DeleteActivity)
	}

	// AI suggestions route
	router.GET("/api/suggestions", activityController.GetSuggestions)

	// WebSocket endpoint
	wsHandler := ws.NewHandler(manager)
	router.GET("/ws", func(c *gin.Context) {
		wsHandler.HandleWebSocket(c.Writer, c.Request)
	})

	return router
}
