// routes.go

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/app/controllers" // Import controller package
	"import "github.com/app/core-system/server/controllers" // Import controller package
)

// AppError represents a generic error type with a message and request ID
type AppError struct {
	Message   string
	RequestID uuid.UUID
}

func (e *AppError) Error() string {
	return e.Message
}

// UserError extends from AppError
type UserError struct {
	AppError
}

// ErrorResponse is used to structure error messages
type ErrorResponse struct {
	Error     string    `json:"error"`
	RequestID uuid.UUID `json:"request_id"`
}

// Request and Response objects
type GetResponseRequest struct {
	Text string `json:"text"`
}

type GetResponseResponse struct {
	Response  string    `json:"response"`
	RequestID uuid.UUID `json:"request_id"`
}

// Middleware to generate and attach a request ID to the context
func requestIDMiddleware(c *gin.Context) {
	requestID := uuid.New() // Generating a new request ID
	c.Set("requestID", requestID)
	c.Next()
}

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Apply middleware to generate and attach request ID to context
	r.Use(requestIDMiddleware)

	// Existing route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Add new route to call Python API controller
	r.POST("/get-response", controllers.GetResponse)

	// Other routes...

	return r
}
