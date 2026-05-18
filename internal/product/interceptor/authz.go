package interceptor

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/skhanal5/payflow/internal/shared/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthzRules map[string][]string

func NewAuthzRules() AuthzRules {
	return AuthzRules{}
}

func AuthzInterceptor(logger zerolog.Logger, rules AuthzRules) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Debug().Str("method", info.FullMethod).Msg("Intercepting gRPC call for authorization")

		requiredRoles, found := rules[info.FullMethod]
		if !found {
			logger.Debug().Str("method", info.FullMethod).Msg("Method not in authz rules. Allowing without authentication.")
			return handler(ctx, req)
		}

		claims, ok := auth.GetUserClaimsFromContext(ctx)
		if !ok {
			logger.Warn().Str("method", info.FullMethod).Msg("No user claims found in context. Denying access.")
			return nil, status.Error(codes.Unauthenticated, "Authentication required")
		}

		logger.Debug().Str("method", info.FullMethod).Str("userID", claims.UserID).Str("role", claims.Role).Msg("User claims found")

		if len(requiredRoles) == 0 {
			logger.Debug().Str("method", info.FullMethod).Msg("Method found in rules with no specific roles required. Allowing.")
			return handler(ctx, req)
		}

		isAuthorized := false
		for _, allowedRole := range requiredRoles {
			if claims.Role == allowedRole {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			logger.Warn().Str("method", info.FullMethod).Str("userID", claims.UserID).Str("role", claims.Role).
				Strs("requiredRoles", requiredRoles).Msg("User not authorized for this method")
			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("Access denied. Required roles: %s", strings.Join(requiredRoles, ", ")))
		}

		logger.Info().Str("method", info.FullMethod).Str("userID", claims.UserID).Msg("User authorized. Proceeding to handler.")
		return handler(ctx, req)
	}
}
