package request

import (
	"order/internal/presentation/grpc/response"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
