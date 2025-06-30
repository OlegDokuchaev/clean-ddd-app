package request

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"warehouse/internal/presentation/grpc/response"
)

func ParseUUID(key string) (uuid.UUID, error) {
	if id, err := uuid.Parse(key); err != nil {
		return uuid.Nil, response.ErrInvalidID
	} else {
		return id, nil
	}
}

func ParseDecimal(key float64) decimal.Decimal {
	return decimal.NewFromFloat(key)
}
