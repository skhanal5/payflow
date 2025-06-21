package middleware

import (
	"github.com/rs/zerolog"
	"net/http"
)

func LoggingMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info().
				Str("method", r.Method).
				Str("url", r.URL.Path).
				Msg("Incoming HTTP request")

			next.ServeHTTP(w, r)
		})
	}
}
