package customer

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrCustomerAlreadyExists      = errors.New("customer already exists")
	ErrCustomerPhoneAlreadyExists = errors.New("customer phone already exists")
	ErrCustomerNotFound           = errors.New("customer not found")
)

func ParseError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrCustomerNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return parsePgError(pgErr)
	}

	return err
}

func parsePgError(err *pgconn.PgError) error {
	switch err.ConstraintName {
	case "customers_pkey":
		return ErrCustomerAlreadyExists

	case "customers_phone_key":
		return ErrCustomerPhoneAlreadyExists

	default:
		return fmt.Errorf("customer not saved: %v", err)
	}
}
