package product

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrProductNotFound = errors.New("product not found")

func ParseError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrProductNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fmt.Errorf("product not saved: %v", pgErr)
	}

	return err
}
