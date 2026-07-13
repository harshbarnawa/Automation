package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken is returned when a token is malformed, expired or signed
// with an unexpected key.
var ErrInvalidToken = errors.New("invalid token")

// Claims is the JWT payload carried by an access token.
type Claims struct {
	UserID string `json:"uid"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Clock returns the current time; it is injectable for deterministic tests.
type Clock func() time.Time

// TokenManager issues and verifies signed JWT access tokens.
type TokenManager struct {
	secret    []byte
	accessTTL time.Duration
	issuer    string
	now       Clock
}

// NewTokenManager builds a TokenManager. A non-positive TTL falls back to 15
// minutes so tokens always expire.
func NewTokenManager(secret string, accessTTL time.Duration, issuer string) *TokenManager {
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}
	return &TokenManager{
		secret:    []byte(secret),
		accessTTL: accessTTL,
		issuer:    issuer,
		now:       time.Now,
	}
}

// WithClock overrides the internal clock and returns the manager for chaining.
func (m *TokenManager) WithClock(clock Clock) *TokenManager {
	m.now = clock
	return m
}

// AccessTTL reports the configured access-token lifetime.
func (m *TokenManager) AccessTTL() time.Duration {
	return m.accessTTL
}

// Generate signs an access token for the given user.
func (m *TokenManager) Generate(userID, email string) (string, error) {
	issuedAt := m.now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(m.accessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signed, nil
}

// Verify parses and validates an access token, returning its claims.
func (m *TokenManager) Verify(tokenString string) (Claims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithTimeFunc(m.now),
	)

	var claims Claims
	_, err := parser.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return m.secret, nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if claims.UserID == "" {
		return Claims{}, ErrInvalidToken
	}

	return claims, nil
}
