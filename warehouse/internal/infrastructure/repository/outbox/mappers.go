package outbox

import (
	"encoding/json"
	outboxDomain "warehouse/internal/domain/outbox"
	"warehouse/internal/infrastructure/db/tables"
)

func ToDomain(model *tables.OutboxMessage) (*outboxDomain.Message, error) {
	var metadata map[string]any
	if err := json.Unmarshal(model.Metadata, &metadata); err != nil {
		return nil, err
	}

	return &outboxDomain.Message{
		ID:       model.ID,
		Name:     model.Name,
		Payload:  model.Payload,
		Metadata: metadata,
	}, nil
}

func ToDomains(models []*tables.OutboxMessage) ([]*outboxDomain.Message, error) {
	domains := make([]*outboxDomain.Message, 0, len(models))
	for _, model := range models {
		domain, err := ToDomain(model)
		if err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}
	return domains, nil
}

func ToModel(domain *outboxDomain.Message) (*tables.OutboxMessage, error) {
	bytes, err := json.Marshal(domain.Metadata)
	if err != nil {
		return nil, err
	}

	return &tables.OutboxMessage{
		ID:       domain.ID,
		Name:     domain.Name,
		Payload:  domain.Payload,
		Metadata: bytes,
	}, nil
}
