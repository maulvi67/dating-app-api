package middleware

import (
	"dating-apps/helper/logger"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func Adapt(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func ServeHTTP(next http.Handler, log logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Serve the request by calling the next handler.
		next.ServeHTTP(w, r)

		// Log details about the request.
		duration := time.Since(start)
		log.Info("request served",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", duration,
			"remote", r.RemoteAddr,
		)
	})
}
