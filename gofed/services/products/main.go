package main

import (
	"net/http"
	"os"
	"products/graph"
	"products/handlers"
	"products/logger"
	"products/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
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

	// Configurar mux
	mux := http.NewServeMux()

	mux.HandleFunc("/query", srv.ServeHTTP)
	mux.HandleFunc("/healthz", handlers.HealthHandler(logger))

	if os.Getenv("GRAPHQL_PLAYGROUND_ENABLED") != "false" {
		mux.HandleFunc("/", playground.Handler("GraphQL playground", "/query"))
	}

	// Logging middleware
	handlerWithLogging := middleware.LoggingMiddleware(logger)(mux)

	logger.WithFields(map[string]interface{}{
		"port": port,
		"endpoints": []string{
			"http://localhost:" + port + "/ (GraphQL Playground)",
			"http://localhost:" + port + "/query (GraphQL Query)",
			"http://localhost:" + port + "/healthz (Health Check)",
		},
		"semaphore_max": resolver.Semaphore().MaxCount(),
	}).Info("Products service starting with semaphore")

	if err := http.ListenAndServe(":"+port, handlerWithLogging); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}
