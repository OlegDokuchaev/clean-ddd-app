package courier

import (
	courierDomain "courier/internal/domain/courier"
	"courier/internal/infrastructure/db/tables"
)

func ToDomain(courier *tables.Courier) *courierDomain.Courier {
	return &courierDomain.Courier{
		ID:       courier.ID,
		Name:     courier.Name,
		Phone:    courier.Phone,
		Password: courier.Password,
		Created:  courier.Created,
	}
}

func ToDomains(couriers []*tables.Courier) []*courierDomain.Courier {
	domains := make([]*courierDomain.Courier, 0, len(couriers))
	for _, courier := range couriers {
		domains = append(domains, ToDomain(courier))
	}
	return domains
}

func ToModel(courier *courierDomain.Courier) *tables.Courier {
	return &tables.Courier{
		ID:       courier.ID,
		Name:     courier.Name,
		Phone:    courier.Phone,
		Password: courier.Password,
		Created:  courier.Created,
	}
}
