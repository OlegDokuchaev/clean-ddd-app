package outbox

import (
	"context"
	"github.com/google/uuid"
	outboxDomain "warehouse/internal/domain/outbox"
	"warehouse/internal/infrastructure/db/tables"

	"gorm.io/gorm"
)

type RepositoryImpl struct {
	db *gorm.DB
}

func New(db *gorm.DB) *RepositoryImpl {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Create(ctx context.Context, message *outboxDomain.Message) error {
	model := ToModel(message)
	res := r.db.WithContext(ctx).Create(model)
	return ParseError(res.Error)
}

func (r *RepositoryImpl) GetByID(ctx context.Context, messageID uuid.UUID) (*outboxDomain.Message, error) {
	var model tables.OutboxMessage
	res := r.db.WithContext(ctx).First(&model, "id = ?", messageID)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return ToDomain(&model), nil
}

func (r *RepositoryImpl) GetAll(ctx context.Context) ([]*outboxDomain.Message, error) {
	var models []*tables.OutboxMessage
	res := r.db.WithContext(ctx).Find(&models)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return ToDomains(models), nil
}

func (r *RepositoryImpl) Delete(ctx context.Context, message *outboxDomain.Message) error {
	model := ToModel(message)
	res := r.db.WithContext(ctx).Delete(model)
	if res.RowsAffected == 0 {
		return ErrOutboxMessageNotFound
	}
	return ParseError(res.Error)
}
