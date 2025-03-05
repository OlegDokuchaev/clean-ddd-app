package auth

import "github.com/google/uuid"

type TokenService interface {
	Generate(courierID uuid.UUID) (string, error)
	Decode(token string) (uuid.UUID, error)
}
