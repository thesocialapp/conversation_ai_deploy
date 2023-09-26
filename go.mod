module conversational_ai

go 1.19

// Web framework
require github.com/gin-gonic/gin v1.9.1  

// Structured logging
require github.com/sirupsen/logrus v1.9.3   

// OpenTelemetry tracing
require go.opentelemetry.io/otel v1.11.2

// Prometheus metrics
require (
	github.com/prometheus/client_golang v1.14.0
	github.com/prometheus/client_golang/prometheus/promauto v0.7.0
)

// Exponential backoff 
require github.com/cenkalti/backoff/v4 v4.2.0

// Testing
require (
	github.com/golang/mock v1.6.0
	github.com/stretchr/testify v1.8.1 
)