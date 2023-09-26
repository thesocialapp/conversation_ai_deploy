// services/api.go

package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/app/models"
	"github.com/cenkalti/backoff/v4"
)

// Logger encapsulates logger logic
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{logrus.New()}
}

// Metrics
var (
	requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "api_request_count",
		Help: "The total number of processed requests",
	}, []string{"method", "endpoint", "status_code"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "api_request_duration_seconds",
		Help: "The duration of the processed request",
	}, []string{"method", "endpoint"})
)

// CallPythonAPI calls the Python service API endpoint
func CallPythonAPI(ctx context.Context, text string) (*models.PythonAPIResponse, error) {
	// Extract the requestID from the context
	requestID := ctx.Value("requestID").(string)

	// Create a new logger instance
	logger := NewLogger()

	tracer := otel.Tracer("CallPythonAPI")
	ctx, span := tracer.Start(ctx, "CallPythonAPI")
	defer span.End()

	// Propagate the trace context
	ctx = otel.SetTextMapPropagator(propagation.TraceContext{}).Extract(ctx, propagation.HeaderCarrier{})
	defer span.End()

	// 1. Construct request body
	requestBody := models.PythonAPIRequest{
		Text: text,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"requestID": requestID,
		}).Error("Error marshaling request body")
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	// Log outgoing API call
	logger.WithFields(logrus.Fields{
		"url":       "http://localhost:5000/process-text",
		"body":      jsonBody,
		"requestID": requestID,
	}).Info("Making API request to Python service")

	startTime := time.Now()

	// Define the operation to be retried with backoff.
	operation := func() error {
		// 2. Make HTTP request
		url := "http://localhost:5000/process-text"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return err
		}
		req.Header.Set("X-Request-ID", requestID) // Set X-Request-ID header
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return backoff.Permanent(err)
		}
		defer resp.Body.Close()

		// 3. Handle non-200 status
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("non-200 status code: %d", resp.StatusCode)
		}

		// 4. Parse response
		var apiResponse models.PythonAPIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResponse)
		if err != nil {
			return err
		}

		logger.WithFields(logrus.Fields{
			"response":  apiResponse,
			"requestID": requestID,
		}).Info("Received response from Python API")

		// Update metrics
		requestCount.WithLabelValues("POST", url, fmt.Sprint(resp.StatusCode)).Inc()
		requestDuration.WithLabelValues("POST", url).Observe(time.Since(startTime).Seconds())

		return nil
	}

	// Execute the operation with exponential backoff.
	err = backoff.Retry(operation, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"requestID": requestID,
		}).Error("Error after retries")
		return nil, err
	}

	return &apiResponse, nil
}
