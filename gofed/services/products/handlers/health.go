package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// HealthHandler returns the health status of the service
func HealthHandler(logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.WithField("endpoint", "/healthz").Info("Health check requested")

		response := HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Service:   "products-service",
			Version:   "1.0.0",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.WithError(err).Error("Failed to encode health response")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger.WithField("status", "healthy").Info("Health check completed")
	}
}
