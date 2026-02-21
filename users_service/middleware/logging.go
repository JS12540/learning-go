package middleware

import (
	"net/http"
	"time"
	"users_service/utils"
)

func LoggingMiddleware(logger *utils.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		logger.Info(r.Method + " " + r.URL.Path + " " + time.Since(start).String())
	})
}
