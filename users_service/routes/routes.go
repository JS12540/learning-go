package routes

import (
	"net/http"
	"users_service/handlers"
	"users_service/middleware"
	"users_service/utils"
)

func SetupRoutes(logger *utils.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/users", handlers.GetUsers)

	return middleware.LoggingMiddleware(logger, mux)
}
