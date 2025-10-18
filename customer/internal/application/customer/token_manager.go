package customer

import (
	"github.com/google/uuid"
)

type TokenManager interface {
	GenerateAccess(customerID uuid.UUID) (string, error)
	ParseAndValidateAccess(token string) (uuid.UUID, error)
	GenerateReset(customerID uuid.UUID) (string, error)
	ParseAndValidateReset(token string) (uuid.UUID, error)
}
