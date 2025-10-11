package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type (
	TokenType string
)

const (
	AccessToken TokenType = "access"
	ResetToken  TokenType = "reset"
)

type tokenClaims struct {
	jwt.RegisteredClaims

	Type       TokenType `json:"type"`
	CustomerID uuid.UUID `json:"customer_id"`
}
