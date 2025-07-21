package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "graphql_requests_total",
			Help: "Total of GraphQL requests by service and endpoint",
		},
		[]string{"service", "endpoint", "operation_type"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "graphql_request_duration_seconds",
			Help:    "Duration of GraphQL requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "endpoint", "operation_type"},
	)

	ActiveRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "graphql_active_requests",
			Help: "Number of active GraphQL requests",
		},
		[]string{"service"},
	)

	ErrorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "graphql_errors_total",
			Help: "Total of GraphQL errors by service and type",
		},
		[]string{"service", "error_type"},
	)

	SemaphoreCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "semaphore_current",
			Help: "Number of goroutines using the semaphore",
		},
		[]string{"service"},
	)

	SemaphoreMax = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "semaphore_max",
			Help: "Maximum number of goroutines allowed by the semaphore",
		},
		[]string{"service"},
	)
)

// MetricsMiddleware - Middleware to collect metrics
func MetricsMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ActiveRequests.WithLabelValues(serviceName).Inc()
			defer ActiveRequests.WithLabelValues(serviceName).Dec()

			operationType := "query"
			if r.Method == "POST" {
				operationType = "mutation"
			}

			RequestCounter.WithLabelValues(serviceName, r.URL.Path, operationType).Inc()

			next.ServeHTTP(w, r)

			duration := time.Since(start).Seconds()
			RequestDuration.WithLabelValues(serviceName, r.URL.Path, operationType).Observe(duration)
		})
	}
}

// RecordError - Record error in metrics
func RecordError(serviceName, errorType string) {
	ErrorCounter.WithLabelValues(serviceName, errorType).Inc()
}

// UpdateSemaphoreMetrics - Update semaphore metrics
func UpdateSemaphoreMetrics(serviceName string, current, max int) {
	SemaphoreCurrent.WithLabelValues(serviceName).Set(float64(current))
	SemaphoreMax.WithLabelValues(serviceName).Set(float64(max))
}

// TraceMiddleware - Middleware to add TraceID to context
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate TraceID if not exists
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = generateTraceID()
		}

		// Add TraceID to context
		ctx := context.WithValue(r.Context(), "trace_id", traceID)
		ctx = context.WithValue(ctx, "start_time", time.Now())

		// Add TraceID to response header
		w.Header().Set("X-Trace-ID", traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetTraceID - Get TraceID from context
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return "unknown"
}

// GetStartTime - Get start time from context
func GetStartTime(ctx context.Context) time.Time {
	if startTime, ok := ctx.Value("start_time").(time.Time); ok {
		return startTime
	}
	return time.Now()
}

// generateTraceID - Generate unique ID for tracing
func generateTraceID() string {
	return uuid.New().String()
}
