package outbox

import (
	outboxDomain "warehouse/internal/domain/outbox"
	"warehouse/internal/infrastructure/db/tables"
)

func ToDomain(model *tables.OutboxMessage) *outboxDomain.Message {
	return &outboxDomain.Message{
		ID:      model.ID,
		Type:    model.Type,
		Payload: model.Payload,
	}
}

func ToDomains(models []*tables.OutboxMessage) []*outboxDomain.Message {
	domains := make([]*outboxDomain.Message, 0, len(models))
	for _, model := range models {
		domains = append(domains, ToDomain(model))
	}
	return domains
}

func ToModel(domain *outboxDomain.Message) *tables.OutboxMessage {
	return &tables.OutboxMessage{
		ID:      domain.ID,
		Type:    domain.Type,
		Payload: domain.Payload,
	}
}
