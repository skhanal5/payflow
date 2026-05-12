package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserClaimsContextKey contextKey = "userClaims"
)

type UserClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"` // e.g., "admin", "user", "guest"
	jwt.RegisteredClaims
}

func GetUserClaimsFromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(UserClaimsContextKey).(*UserClaims)
	return claims, ok
}
