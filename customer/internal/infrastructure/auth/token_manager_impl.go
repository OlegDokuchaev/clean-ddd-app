package auth

import (
	customerApplication "customer/internal/application/customer"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManagerImpl struct {
	config *Config
}

func NewTokenManager(config *Config) *TokenManagerImpl {
	return &TokenManagerImpl{config: config}
}

func (m *TokenManagerImpl) GenerateAccess(customerID uuid.UUID) (string, error) {
	claims := m.createTokenClaims(AccessToken, customerID, m.config.AccessTTL)
	return m.createToken(claims)
}

func (m *TokenManagerImpl) ParseAndValidateAccess(token string) (uuid.UUID, error) {
	c, err := m.parse(token)
	if err != nil {
		return uuid.Nil, err
	}
	if c.Type != AccessToken {
		return uuid.Nil, ErrInvalidToken
	}
	return c.CustomerID, nil
}

func (m *TokenManagerImpl) GenerateReset(customerID uuid.UUID) (string, error) {
	claims := m.createTokenClaims(ResetToken, customerID, m.config.ResetTTL)
	return m.createToken(claims)
}

func (m *TokenManagerImpl) ParseAndValidateReset(token string) (uuid.UUID, error) {
	c, err := m.parse(token)
	if err != nil {
		return uuid.Nil, err
	}
	if c.Type != ResetToken {
		return uuid.Nil, ErrInvalidToken
	}
	return c.CustomerID, nil
}

func (m *TokenManagerImpl) createTokenClaims(type_ TokenType, customerID uuid.UUID, ttl time.Duration) *tokenClaims {
	expirationTime := time.Now().Add(ttl)
	return &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Type:       type_,
		CustomerID: customerID,
	}
}

func (m *TokenManagerImpl) createToken(claims *tokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.SigningKey))
}

func (m *TokenManagerImpl) parse(raw string) (*tokenClaims, error) {
	tok, err := jwt.ParseWithClaims(
		raw,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidSigningMethod
			}
			return []byte(m.config.SigningKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}
	claims, ok := tok.Claims.(*tokenClaims)
	if !ok || !tok.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

var _ customerApplication.TokenManager = (*TokenManagerImpl)(nil)
