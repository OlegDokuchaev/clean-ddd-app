package order

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

func ParseError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrOrderNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.ConstraintName {
		default:
			return fmt.Errorf("order not saved: %v", pgErr)
		}
	}

	return err
}
