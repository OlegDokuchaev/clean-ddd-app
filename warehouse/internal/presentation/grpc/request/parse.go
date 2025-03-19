package request

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"warehouse/internal/presentation/grpc/response"
)

func parseUUID(key string) (uuid.UUID, error) {
	if id, err := uuid.Parse(key); err != nil {
		return uuid.Nil, response.ErrInvalidID
	} else {
		return id, nil
	}
}

func parseDecimal(key float64) decimal.Decimal {
	return decimal.NewFromFloat(key)
}
