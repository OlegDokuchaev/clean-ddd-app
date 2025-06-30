package customer

import (
	"github.com/google/uuid"
)

type TokenManager interface {
	Generate(customerID uuid.UUID) (string, error)
	Decode(token string) (uuid.UUID, error)
}
