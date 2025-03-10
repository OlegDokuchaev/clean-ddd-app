package order

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderNotFound      = errors.New("order not found")
)

func ParseError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrOrderNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return parsePgError(pgErr)
	}

	return err
}

func parsePgError(err *pgconn.PgError) error {
	switch err.ConstraintName {
	case "orders_pkey":
		return ErrOrderAlreadyExists

	default:
		return fmt.Errorf("order not saved: %v", err)
	}
}
