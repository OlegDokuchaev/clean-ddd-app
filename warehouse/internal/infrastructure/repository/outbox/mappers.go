package outbox

import (
	outboxDomain "warehouse/internal/domain/outbox"
	"warehouse/internal/infrastructure/db/tables"
)

func ToDomain(message *tables.OutboxMessage) *outboxDomain.Message {
	return &outboxDomain.Message{
		ID:      message.ID,
		Type:    message.Type,
		Payload: message.Payload,
	}
}

func ToDomains(messages []*tables.OutboxMessage) []*outboxDomain.Message {
	domains := make([]*outboxDomain.Message, 0, len(messages))
	for _, message := range messages {
		domains = append(domains, ToDomain(message))
	}
	return domains
}

func ToModel(message *outboxDomain.Message) *tables.OutboxMessage {
	return &tables.OutboxMessage{
		ID:      message.ID,
		Type:    message.Type,
		Payload: message.Payload,
	}
}
