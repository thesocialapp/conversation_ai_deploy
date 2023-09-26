// controllers.go

package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/app/models"
	"github.com/app/services"
	"github.com/sirupsen/logrus" // Import logrus for structured logging
)

// Initialize the logger
var log = logrus.New()

// Envelope is used to structure the standard API response
type Envelope struct {
	Data  interface{}     `json:"data,omitempty"`
	Error string          `json:"error,omitempty"`
	Meta  *PaginationMeta `json:"meta,omitempty"`
}

// PaginationMeta holds metadata for pagination
type PaginationMeta struct {
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
	PerPage    int `json:"per_page"`
	Page       int `json:"page"`
}

// GetResponse calls Python API and handles response
func GetResponse(c *gin.Context) {
	env, status := getResponse(c)
	c.JSON(status, env)
}

func getResponse(c *gin.Context) (Envelope, int) {
	// Log received request
	requestID, exists := c.Get("requestID")
	log.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Info("Received a get response request")

	if !exists {
		err := models.AppError{
			Message: "no request ID found",
			RequestID: uuid.Nil,
		}
		log.Println(err.Error()) // log the error using old log
		// Log the error with logrus
		log.WithFields(logrus.Fields{
			"error":     err.Error(),
			"requestID": "nil",
		}).Error("No request ID found")
		return Envelope{Error: err.Error()}, http.StatusInternalServerError
	}

	// Bind the JSON body to the request structure
	var req models.GetPythonAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := models.AppError{
			Message: "invalid request payload",
			RequestID: requestID.(uuid.UUID),
		}
		log.Println(err.Error()) // log the error using old log
		// Log the error with logrus
		log.WithFields(logrus.Fields{
			"error":     err.Error(),
			"requestID": requestID,
		}).Error("Invalid request payload")
		return Envelope{Error: err.Error()}, http.StatusBadRequest
	}

	// Call Python API
	resp, err := services.CallPythonAPI(req.Text)

	// Handle errors
	if err != nil {
		err := models.AppError{
			Message: err.Error(),
			RequestID: requestID.(uuid.UUID),
		}
		log.Println(err.Error()) // log the error using old log
		// Log the error with logrus
		log.WithFields(logrus.Fields{
			"error":     err.Error(),
			"requestID": requestID,
		}).Error("Error calling Python API")
		return Envelope{Error: err.Error()}, http.StatusInternalServerError
	}

	// Log success using old log
	log.Printf("Successfully fetched response for requestID: %s", requestID)
	// Log success with logrus
	log.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Info("Successfully fetched response")

	// Process response and Return response in an Envelope
	responseText := resp.Text
	return Envelope{
		Data: models.GetResponseResponse{
			Response:  responseText,
			RequestID: requestID.(uuid.UUID),
		},
	}, http.StatusOK
}
