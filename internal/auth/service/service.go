package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/skhanal5/payflow/gen/go/auth"
	"github.com/skhanal5/payflow/internal/auth/repository"
	"github.com/skhanal5/payflow/internal/shared/auth"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	repo      repository.UserRepository
	jwtSecret string
	logger    *zerolog.Logger
}

func NewAuthService(repo repository.UserRepository, jwtSecret string, logger *zerolog.Logger) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	salt, err := generateSalt()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user")
	}

	hash := hashPassword(req.Password, salt)
	encoded := fmt.Sprintf("%s:%s", hex.EncodeToString(salt), hex.EncodeToString(hash))

	user := &repository.User{
		UserID:         req.UserId,
		HashedPassword: encoded,
	}

	if _, err := s.repo.InsertUser(ctx, user); err != nil {
		s.logger.Error().Err(err).Str("user_id", req.UserId).Msg("Failed to insert user")
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	token, err := s.signJWT(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}

	return &pb.AuthResponse{Token: token}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	user, err := s.repo.GetUserByUserID(ctx, req.UserId)
	if err != nil {
		s.logger.Warn().Str("user_id", req.UserId).Msg("Login attempt for non-existent user")
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	parts, err := decodeHash(user.HashedPassword)
	if err != nil || !verifyPassword(req.Password, parts.salt, parts.hash) {
		s.logger.Warn().Str("user_id", req.UserId).Msg("Invalid password")
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	token, err := s.signJWT(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}

	return &pb.AuthResponse{Token: token}, nil
}

func (s *AuthService) signJWT(userID string) (string, error) {
	claims := &auth.UserClaims{
		UserID: userID,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

type hashParts struct {
	salt []byte
	hash []byte
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	return salt, err
}

func hashPassword(password string, salt []byte) []byte {
	h := sha256.Sum256(append(salt, []byte(password)...))
	return h[:]
}

func verifyPassword(password string, salt, expectedHash []byte) bool {
	h := sha256.Sum256(append(salt, []byte(password)...))
	return hex.EncodeToString(h[:]) == hex.EncodeToString(expectedHash)
}

func decodeHash(encoded string) (hashParts, error) {
	parts := hashParts{}
	colon := -1
	for i, b := range encoded {
		if b == ':' {
			colon = i
			break
		}
	}
	if colon < 0 {
		return parts, fmt.Errorf("invalid hash format")
	}
	var err error
	parts.salt, err = hex.DecodeString(encoded[:colon])
	if err != nil {
		return parts, err
	}
	parts.hash, err = hex.DecodeString(encoded[colon+1:])
	return parts, err
}
