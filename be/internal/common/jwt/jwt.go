package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"be/internal/config"
)

type Claims struct {
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	RoleID      string   `json:"roleId"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret    []byte
	accessTTL time.Duration
}

func NewManager(cfg config.Config) *Manager {
	return &Manager{
		secret:    []byte(cfg.JWTSecret),
		accessTTL: cfg.JWTAccessTTL,
	}
}

func (m *Manager) SignAccessToken(userID, email, name, role, roleID string, permissions []string) (string, error) {
	now := time.Now()
	claims := Claims{
		Email:       email,
		Name:        name,
		Role:        role,
		RoleID:      roleID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}
	return signed, nil
}

func (m *Manager) VerifyAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("verify access token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token claims")
	}
	return claims, nil
}
