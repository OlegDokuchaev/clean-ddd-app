package outbox

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"
	outboxDomain "warehouse/internal/domain/outbox"
	"warehouse/internal/infrastructure/db/tables"
)

type RepositoryImpl struct {
	db *gorm.DB
}

func New(db *gorm.DB) *RepositoryImpl {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Create(ctx context.Context, message *outboxDomain.Message) error {
	carrier := propagation.MapCarrier{}
	propagation.TraceContext{}.Inject(ctx, carrier)
	for k, v := range carrier {
		message.Metadata[k] = v
	}

	model, err := ToModel(message)
	if err != nil {
		return err
	}

	res := r.db.WithContext(ctx).Create(model)
	return ParseError(res.Error)
}

func (r *RepositoryImpl) GetByID(ctx context.Context, messageID uuid.UUID) (*outboxDomain.Message, error) {
	var model tables.OutboxMessage
	res := r.db.WithContext(ctx).First(&model, "id = ?", messageID)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return ToDomain(&model)
}

func (r *RepositoryImpl) GetAll(ctx context.Context) ([]*outboxDomain.Message, error) {
	var models []*tables.OutboxMessage
	res := r.db.WithContext(ctx).Find(&models)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return ToDomains(models)
}

func (r *RepositoryImpl) Delete(ctx context.Context, message *outboxDomain.Message) error {
	model, err := ToModel(message)
	if err != nil {
		return err
	}

	res := r.db.WithContext(ctx).Delete(model)
	if res.RowsAffected == 0 {
		return ErrOutboxMessageNotFound
	}

	return ParseError(res.Error)
}
