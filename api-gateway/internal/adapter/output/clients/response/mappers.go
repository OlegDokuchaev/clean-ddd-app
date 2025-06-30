package response

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func ToUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return id, ErrInternalServerError
	}
	return id, nil
}

func ToDecimal(key float64) decimal.Decimal {
	return decimal.NewFromFloat(key)
}
