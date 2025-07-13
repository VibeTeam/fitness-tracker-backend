package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType distinguishes between access and refresh tokens.
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims represents the JWT payload used across the system.
type Claims struct {
	UserID int32     `json:"user_id"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

// Manager handles creation and verification of JWT access/refresh tokens.
type Manager struct {
	accessSecret    []byte
	refreshSecret   []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewManager builds a Manager supplied with secrets and TTL values.
// Using different secrets for access and refresh tokens simplifies rotation.
func NewManager(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{
		accessSecret:    []byte(accessSecret),
		refreshSecret:   []byte(refreshSecret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

// NewTokens returns freshly minted access and refresh tokens for the given user.
func (m *Manager) NewTokens(userID int32) (string, string, error) {
	access, err := m.newToken(userID, AccessToken, m.accessTokenTTL, m.accessSecret)
	if err != nil {
		return "", "", err
	}
	refresh, err := m.newToken(userID, RefreshToken, m.refreshTokenTTL, m.refreshSecret)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// ValidateAccessToken parses and validates access token, returning embedded user ID.
func (m *Manager) ValidateAccessToken(tokenStr string) (int32, error) {
	claims, err := m.parseToken(tokenStr, m.accessSecret, AccessToken)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// RefreshTokens validates provided refresh token and returns a new token pair.
func (m *Manager) RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := m.parseToken(refreshToken, m.refreshSecret, RefreshToken)
	if err != nil {
		return "", "", err
	}
	return m.NewTokens(claims.UserID)
}

// --- internals ---

func (m *Manager) newToken(userID int32, ttype TokenType, ttl time.Duration, secret []byte) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Type:   ttype,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ErrInvalidToken is returned when token verification fails or token kind mismatches.
var ErrInvalidToken = errors.New("invalid token")

func (m *Manager) parseToken(tokenStr string, secret []byte, expectedType TokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	cl, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || cl.Type != expectedType {
		return nil, ErrInvalidToken
	}
	return cl, nil
}
