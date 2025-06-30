package auth

import (
	customerApplication "customer/internal/application/customer"
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

func (m *TokenManagerImpl) Generate(customerID uuid.UUID) (string, error) {
	claims := m.createTokenClaims(customerID)
	return m.createToken(claims)
}

func (m *TokenManagerImpl) Decode(tokenString string) (uuid.UUID, error) {
	token, err := m.parseToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	claims, err := m.parseTokenClaims(token)
	if err != nil {
		return uuid.Nil, err
	}

	return claims.CustomerID, nil
}

func (m *TokenManagerImpl) createTokenClaims(customerID uuid.UUID) *tokenClaims {
	expirationTime := time.Now().Add(m.config.TokenTTL)
	return &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		CustomerID: customerID,
	}
}

func (m *TokenManagerImpl) createToken(claims *tokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.SigningKey))
}

func (m *TokenManagerImpl) parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidSigningMethod
			}
			return []byte(m.config.SigningKey), nil
		},
	)

	if err != nil {
		if jwt.ErrTokenExpired.Error() == err.Error() {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	return token, nil
}

func (m *TokenManagerImpl) parseTokenClaims(token *jwt.Token) (*tokenClaims, error) {
	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

var _ customerApplication.TokenManager = (*TokenManagerImpl)(nil)
