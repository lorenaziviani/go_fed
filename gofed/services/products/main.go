package main

import (
	"net/http"

	"products/graph"
	"products/handlers"
	"products/logger"
	"products/middleware"

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
		"port": "8082",
		"endpoints": []string{
			"http://localhost:8082/ (GraphQL Playground)",
			"http://localhost:8082/query (GraphQL Query)",
			"http://localhost:8082/healthz (Health Check)",
		},
	}).Info("Products service starting")

	if err := http.ListenAndServe(":8082", handlerWithLogging); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}
