package item

import (
	"errors"
	"fmt"
	productRepository "warehouse/internal/infrastructure/repository/product"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrItemAlreadyExists = errors.New("item already exists")
	ErrItemNotFound      = errors.New("item not found")
	ErrItemsNotFound     = errors.New("items not found")
)

func ParseError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrItemNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return parsePgError(pgErr)
	}

	return err
}

func parsePgError(err *pgconn.PgError) error {
	switch err.ConstraintName {
	case "items_pkey":
		return ErrItemAlreadyExists

	case "items_product_id_fkey":
		return productRepository.ErrProductNotFound

	default:
		return fmt.Errorf("item not saved: %v", err)
	}
}
