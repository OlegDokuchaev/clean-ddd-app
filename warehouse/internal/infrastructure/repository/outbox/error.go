package outbox

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrOutboxMessageAlreadyExists = errors.New("outbox message already exists")
	ErrOutboxMessageNotFound      = errors.New("outbox message not found")
)

func ParseError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrOutboxMessageNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return parsePgError(pgErr)
	}

	return err
}

func parsePgError(err *pgconn.PgError) error {
	switch err.ConstraintName {
	case "outbox_messages_pkey":
		return ErrOutboxMessageAlreadyExists

	default:
		return fmt.Errorf("outbox message not saved: %v", err)
	}
}
