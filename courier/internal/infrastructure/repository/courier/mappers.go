package courier

import (
	courierDomain "courier/internal/domain/courier"
	"courier/internal/infrastructure/db/tables"
)

func ToDomain(model *tables.Courier) *courierDomain.Courier {
	return &courierDomain.Courier{
		ID:       model.ID,
		Name:     model.Name,
		Phone:    model.Phone,
		Password: model.Password,
		Created:  model.Created,
	}
}

func ToDomains(models []*tables.Courier) []*courierDomain.Courier {
	domains := make([]*courierDomain.Courier, 0, len(models))
	for _, model := range models {
		domains = append(domains, ToDomain(model))
	}
	return domains
}

func ToModel(domain *courierDomain.Courier) *tables.Courier {
	return &tables.Courier{
		ID:       domain.ID,
		Name:     domain.Name,
		Phone:    domain.Phone,
		Password: domain.Password,
		Created:  domain.Created,
	}
}
