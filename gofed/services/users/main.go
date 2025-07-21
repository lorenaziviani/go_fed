package main

import (
	"net/http"
	"os"
	"users/graph"
	"users/handlers"
	"users/logger"
	"users/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8081"

func main() {
	port := os.Getenv("USERS_SERVICE_PORT")
	if port == "" {
		port = defaultPort
	}

	// Configure logger
	logger := logger.SetupLogger()

	// Create resolver with cache
	resolver := graph.NewResolver()

	// Configure GraphQL
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Configure mux
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
		"cache_max_size": resolver.Cache().Size(),
		"cache_ttl":      "5m",
	}).Info("Users service starting with cache")

	if err := http.ListenAndServe(":"+port, handlerWithLogging); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}
