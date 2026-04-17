package services

import (
	"capuchin/internal/config"
	"capuchin/internal/database"
	"capuchin/internal/models"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrDatabase           = errors.New("database error")
)

type AuthService interface {
	Signup(email, password string) (*models.User, error)
	Login(email, password string) (string, error)
	Logout(tokenStr string) error
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Signup(email, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	u := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
	}

	_, err = database.GetDB().Exec("INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)", u.ID, u.Email, u.PasswordHash)
	if err != nil {
		errStr := err.Error()
		// Convert storage-specific duplicate key errors into a stable domain error for handlers.
		if strings.Contains(errStr, "unique constraint") || strings.Contains(errStr, "duplicate key value") {
			return nil, ErrUserExists
		}
		return nil, ErrDatabase
	}

	return u, nil
}

func (s *authService) Login(email, password string) (string, error) {
	var u models.User
	err := database.GetDB().QueryRow("SELECT id, email, password_hash FROM users WHERE email=$1", email).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		// Use one response for unknown user and wrong password to avoid account enumeration.
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		// Keep the same error shape to avoid leaking which check failed.
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID.String(),
		"email":   u.Email,
		// Short-lived tokens reduce blast radius if a token is leaked.
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(config.JWTKey)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

func (s *authService) Logout(tokenStr string) error {
	if tokenStr == "" {
		return ErrInvalidToken
	}

	// Accept either raw JWTs or Authorization header values for caller flexibility.
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}

	// Parse to extract expiry so blacklist rows can be garbage-collected safely.
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return config.JWTKey, nil
	})

	if token == nil {
		return ErrInvalidToken
	}

	var expTime time.Time
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			expTime = time.Unix(int64(exp), 0)
		} else {
			// Fail-safe TTL keeps blacklist entries finite even for malformed claim types.
			expTime = time.Now().Add(72 * time.Hour)
		}
	} else {
		return ErrInvalidToken
	}

	// Idempotent logout avoids surfacing harmless duplicate requests as server errors.
	_, err := database.GetDB().Exec("INSERT INTO blacklisted_tokens (token, expired_at) VALUES ($1, $2) ON CONFLICT (token) DO NOTHING", tokenStr, expTime)
	if err != nil {
		return ErrDatabase
	}

	return nil
}
