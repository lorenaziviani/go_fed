package main

import (
	"net/http"

	"users/graph"
	"users/handlers"
	"users/logger"
	"users/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	logger := logger.SetupLogger()

	// Create GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// Create mux for routing
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", handlers.HealthHandler(logger))
	mux.Handle("/query", srv)
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))

	handlerWithLogging := middleware.LoggingMiddleware(logger)(mux)

	logger.WithFields(map[string]interface{}{
		"port": "8081",
		"endpoints": []string{
			"http://localhost:8081/ (GraphQL Playground)",
			"http://localhost:8081/query (GraphQL Query)",
			"http://localhost:8081/healthz (Health Check)",
		},
	}).Info("Users service starting")

	if err := http.ListenAndServe(":8081", handlerWithLogging); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}
