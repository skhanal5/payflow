package auth

import (
	"context"
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const (
	UserClaimsContextKey contextKey = "userClaims"
)

const ClaimsMetadataKey = "x-user-claims"

type UserClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GetUserClaimsFromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(UserClaimsContextKey).(*UserClaims)
	if ok {
		return claims, true
	}
	claims, ok = getUserClaimsFromMetadata(ctx)
	return claims, ok
}

func getUserClaimsFromMetadata(ctx context.Context) (*UserClaims, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, false
	}
	vals := md.Get(ClaimsMetadataKey)
	if len(vals) == 0 {
		return nil, false
	}
	var claims UserClaims
	if err := json.Unmarshal([]byte(vals[0]), &claims); err != nil {
		return nil, false
	}
	return &claims, true
}

func ClaimsToMetadata(claims *UserClaims) (metadata.MD, error) {
	data, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}
	return metadata.Pairs(ClaimsMetadataKey, string(data)), nil
}
