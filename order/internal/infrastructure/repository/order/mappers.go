package order

import (
	orderDomain "order/internal/domain/order"
	"order/internal/infrastructure/db/tables"

	"github.com/google/uuid"
)

func ToDomain(model *tables.Order) *orderDomain.Order {
	return &orderDomain.Order{
		ID:         model.ID,
		CustomerID: model.CustomerID,
		Status:     model.Status,
		Created:    model.Created,
		Version:    model.Version,
		Items:      toItemDomains(model.Items),
		Delivery:   toDeliveryDomain(model.Delivery),
	}
}

func toItemDomain(model tables.OrderItem) orderDomain.Item {
	return orderDomain.Item{
		ProductID: model.ProductID,
		Price:     model.Price,
		Count:     model.Count,
	}
}

func toDeliveryDomain(model tables.Delivery) orderDomain.Delivery {
	return orderDomain.Delivery{
		CourierID: model.CourierID,
		Address:   model.Address,
		Arrived:   model.Arrived,
	}
}

func ToDomains(models []*tables.Order) []*orderDomain.Order {
	domains := make([]*orderDomain.Order, 0, len(models))
	for _, model := range models {
		domains = append(domains, ToDomain(model))
	}
	return domains
}

func toItemDomains(models []tables.OrderItem) []orderDomain.Item {
	domains := make([]orderDomain.Item, 0, len(models))
	for _, model := range models {
		domains = append(domains, toItemDomain(model))
	}
	return domains
}

func ToModel(domain *orderDomain.Order) *tables.Order {
	return &tables.Order{
		ID:         domain.ID,
		CustomerID: domain.CustomerID,
		Status:     domain.Status,
		Created:    domain.Created,
		Version:    domain.Version,
		Items:      toItemModels(domain.ID, domain.Items),
		Delivery:   toDeliveryModel(domain.ID, domain.Delivery),
	}
}

func toItemModels(orderID uuid.UUID, domains []orderDomain.Item) []tables.OrderItem {
	models := make([]tables.OrderItem, 0, len(domains))
	for _, domain := range domains {
		models = append(models, toItemModel(orderID, domain))
	}
	return models
}

func toItemModel(orderID uuid.UUID, domain orderDomain.Item) tables.OrderItem {
	return tables.OrderItem{
		ID:        uuid.New(),
		OrderID:   orderID,
		ProductID: domain.ProductID,
		Price:     domain.Price,
		Count:     domain.Count,
	}
}

func toDeliveryModel(orderID uuid.UUID, domain orderDomain.Delivery) tables.Delivery {
	return tables.Delivery{
		ID:        uuid.New(),
		OrderID:   orderID,
		CourierID: domain.CourierID,
		Address:   domain.Address,
		Arrived:   domain.Arrived,
	}
}
