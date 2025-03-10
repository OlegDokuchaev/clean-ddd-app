package courier

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrCourierAlreadyExists      = errors.New("courier already exists")
	ErrCourierPhoneAlreadyExists = errors.New("courier phone already exists")
	ErrCourierNotFound           = errors.New("courier not found")
)

func ParseError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrCourierNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return parsePgError(pgErr)
	}

	return err
}

func parsePgError(err *pgconn.PgError) error {
	switch err.ConstraintName {
	case "couriers_pkey":
		return ErrCourierAlreadyExists

	case "couriers_phone_key":
		return ErrCourierPhoneAlreadyExists

	default:
		return fmt.Errorf("courier not saved: %v", err)
	}
}
