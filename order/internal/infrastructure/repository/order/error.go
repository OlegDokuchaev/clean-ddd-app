package order

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderNotFound      = errors.New("order not found")
)

func ParseError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrOrderNotFound
	}

	var we mongo.WriteException
	if errors.As(err, &we) {
		for _, writeErr := range we.WriteErrors {
			if writeErr.Code == 11000 {
				return ErrOrderAlreadyExists
			}
		}
		return fmt.Errorf("order not saved: %w", err)
	}

	return err
}
