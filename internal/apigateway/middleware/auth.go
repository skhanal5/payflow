package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"

	"github.com/skhanal5/payflow/internal/shared/auth"
)

func AuthMiddleware(logger zerolog.Logger, jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn().Str("path", r.URL.Path).Msg("Missing Authorization header")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			tokenParts := strings.SplitN(authHeader, " ", 2)
			if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
				logger.Warn().Str("header", authHeader).Msg("Invalid Authorization header format")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			tokenString := tokenParts[1]
			claims := &auth.UserClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenSignatureInvalid
				}
				return []byte(jwtSecret), nil
			}, jwt.WithLeeway(5*time.Second))

			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					logger.Warn().Msg("Token expired")
					http.Error(w, "Token Expired", http.StatusUnauthorized)
					return
				case errors.Is(err, jwt.ErrTokenSignatureInvalid):
					logger.Warn().Msg("Invalid token signature")
				default:
					logger.Warn().Err(err).Msg("Failed to parse token")
				}
				http.Error(w, "Invalid Token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				logger.Warn().Msg("Token is not valid")
				http.Error(w, "Invalid Token", http.StatusUnauthorized)
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), auth.UserClaimsContextKey, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
