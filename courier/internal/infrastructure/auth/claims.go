package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type tokenClaims struct {
	jwt.RegisteredClaims

	CourierID uuid.UUID `json:"courier_id"`
}
