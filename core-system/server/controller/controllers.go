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

if !exists {
	err := models.AppError{
		Message:   "no request ID found",
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
