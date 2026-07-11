package service

import (
	"errors"
	"time"

	"github.com/ikhdamw/e-wallet/auth-service/internal/model"
	"github.com/ikhdamw/e-wallet/auth-service/internal/repository"
	"github.com/ikhdamw/e-wallet/auth-service/pkg/config"
	"github.com/ikhdamw/e-wallet/auth-service/pkg/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *model.RegisterRequest) (*model.TokenResponse, error)
	Login(req *model.LoginRequest) (*model.TokenResponse, error)
	RefreshToken(refreshToken string) (*model.TokenResponse, error)
	GetMe(userID string) (*model.UserResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	redis      *redis.Client
	config     *config.Config
}

func NewAuthService(userRepo repository.UserRepository, redis *redis.Client, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		redis:    redis,
		config:   cfg,
	}
}

func (s *authService) Register(req *model.RegisterRequest) (*model.TokenResponse, error) {
	// Check if user exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &model.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Phone:        req.Phone,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	return s.generateTokens(user)
}

func (s *authService) Login(req *model.LoginRequest) (*model.TokenResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check status
	if user.Status != "active" {
		return nil, errors.New("account is not active")
	}

	// Generate tokens
	return s.generateTokens(user)
}

func (s *authService) RefreshToken(refreshToken string) (*model.TokenResponse, error) {
	// Parse refresh token
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Find user
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens
	return s.generateTokens(user)
}

func (s *authService) GetMe(userID string) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Phone:     user.Phone,
		AvatarURL: user.AvatarURL,
		Status:    user.Status,
	}, nil
}

func (s *authService) generateTokens(user *model.User) (*model.TokenResponse, error) {
	// Access token claims
	accessClaims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Refresh token claims
	refreshClaims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &model.TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    86400, // 24 hours in seconds
	}, nil
}
