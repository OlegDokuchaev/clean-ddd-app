package auth

import "github.com/google/uuid"

type TokenManager interface {
	Generate(courierID uuid.UUID) (string, error)
	Decode(token string) (uuid.UUID, error)
}
