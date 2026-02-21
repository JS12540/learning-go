package main

import (
	"fmt"
	"net/http"
	"users_service/config"
	"users_service/routes"
	"users_service/utils"
)

func main() {
	cfg := config.LoadConfig()
	logger := utils.NewLogger()

	router := routes.SetupRoutes(logger)

	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Info("Starting server on " + addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("Server failed: " + err.Error())
	}
}
