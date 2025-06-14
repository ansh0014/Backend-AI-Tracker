package routes

import (
    "log"
    "time"

    "github.com/gin-gonic/gin"
)

// LoggerMiddleware logs request details
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start time
        startTime := time.Now()

        // Process request
        c.Next()

        // End time
        endTime := time.Now()

        // Execution time
        latencyTime := endTime.Sub(startTime)

        // Request details
        reqMethod := c.Request.Method
        reqURI := c.Request.RequestURI
        statusCode := c.Writer.Status()
        clientIP := c.ClientIP()

        log.Printf("| %3d | %13v | %15s | %s | %s |",
            statusCode,
            latencyTime,
            clientIP,
            reqMethod,
            reqURI,
        )
    }
}

// AuthMiddleware handles authentication
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }
        // TODO: Implement proper token validation
        c.Next()
    }
}

// ErrorMiddleware handles errors globally
func ErrorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            c.JSON(500, gin.H{
                "errors": c.Errors.Errors(),
            })
        }
    }
}