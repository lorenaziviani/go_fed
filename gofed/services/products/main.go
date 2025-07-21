package main

import (
	"net/http"
	"os"
	"products/graph"
	"products/handlers"
	"products/logger"
	"products/middleware"

	"products/metrics"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const defaultPort = "8082"

func main() {
	port := os.Getenv("PRODUCTS_SERVICE_PORT")
	if port == "" {
		port = defaultPort
	}

	// Setup logger
	logger := logger.SetupLogger()

	// Create resolver with semaphore
	resolver := graph.NewResolver()

	// Configure GraphQL
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Configure mux
	mux := http.NewServeMux()

	mux.HandleFunc("/query", srv.ServeHTTP)
	mux.HandleFunc("/healthz", handlers.HealthHandler(logger))
	mux.Handle("/metrics", promhttp.Handler())

	if os.Getenv("GRAPHQL_PLAYGROUND_ENABLED") != "false" {
		mux.HandleFunc("/", playground.Handler("GraphQL playground", "/query"))
	}

	// Middleware chain: Trace -> Metrics -> Logging
	handlerWithMiddleware := metrics.TraceMiddleware(
		metrics.MetricsMiddleware("products")(
			middleware.LoggingMiddleware(logger)(mux),
		),
	)

	logger.WithFields(map[string]interface{}{
		"port": port,
		"endpoints": []string{
			"http://localhost:" + port + "/ (GraphQL Playground)",
			"http://localhost:" + port + "/query (GraphQL Query)",
			"http://localhost:" + port + "/healthz (Health Check)",
			"http://localhost:" + port + "/metrics (Prometheus Metrics)",
		},
		"semaphore_max": resolver.Semaphore().MaxCount(),
		"features":      []string{"semaphore", "metrics", "tracing"},
	}).Info("Products service starting with semaphore, metrics and tracing")

	if err := http.ListenAndServe(":"+port, handlerWithMiddleware); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}
