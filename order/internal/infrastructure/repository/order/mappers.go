package order

import (
	orderDomain "order/internal/domain/order"
	"order/internal/infrastructure/db/documents"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func toDoc(o *orderDomain.Order) *documents.Order {
	return &documents.Order{
		ID:         o.ID.String(),
		CustomerID: o.CustomerID.String(),
		Status:     o.Status,
		Created:    o.Created,
		Version:    o.Version.String(),
		Delivery:   toDeliveryDoc(o.Delivery),
		Items:      toItemsDoc(o.Items),
	}
}

func toItemDoc(domain orderDomain.Item) documents.OrderItem {
	return documents.OrderItem{
		ProductID: domain.ProductID.String(),
		Price:     domain.Price.String(),
		Count:     domain.Count,
	}
}

func toDeliveryDoc(domain orderDomain.Delivery) documents.Delivery {
	var courierID *string
	if domain.CourierID != nil {
		id := domain.CourierID.String()
		courierID = &id
	}

	return documents.Delivery{
		CourierID: courierID,
		Address:   domain.Address,
		Arrived:   domain.Arrived,
	}
}

func toItemsDoc(domains []orderDomain.Item) []documents.OrderItem {
	items := make([]documents.OrderItem, 0, len(domains))
	for _, domain := range domains {
		items = append(items, toItemDoc(domain))
	}
	return items
}

func toDomain(doc *documents.Order) (*orderDomain.Order, error) {
	id, err := uuid.Parse(doc.ID)
	if err != nil {
		return nil, err
	}
	customerID, err := uuid.Parse(doc.CustomerID)
	if err != nil {
		return nil, err
	}
	version, err := uuid.Parse(doc.Version)
	if err != nil {
		return nil, err
	}

	items, err := toItemsDomain(doc.Items)
	if err != nil {
		return nil, err
	}

	delivery, err := toDeliveryDomain(doc.Delivery)
	if err != nil {
		return nil, err
	}

	return &orderDomain.Order{
		ID:         id,
		CustomerID: customerID,
		Status:     doc.Status,
		Created:    doc.Created,
		Version:    version,
		Delivery:   delivery,
		Items:      items,
	}, nil
}

func toItemDomain(doc documents.OrderItem) (orderDomain.Item, error) {
	prodID, err := uuid.Parse(doc.ProductID)
	if err != nil {
		return orderDomain.Item{}, err
	}
	price, err := decimal.NewFromString(doc.Price)
	if err != nil {
		return orderDomain.Item{}, err
	}

	return orderDomain.Item{
		ProductID: prodID,
		Price:     price,
		Count:     doc.Count,
	}, nil
}

func toDeliveryDomain(doc documents.Delivery) (orderDomain.Delivery, error) {
	var courierID *uuid.UUID
	if doc.CourierID != nil {
		tmp, err := uuid.Parse(*doc.CourierID)
		if err != nil {
			return orderDomain.Delivery{}, err
		}
		courierID = &tmp
	}

	return orderDomain.Delivery{
		CourierID: courierID,
		Address:   doc.Address,
		Arrived:   doc.Arrived,
	}, nil
}

func toDomains(docs []documents.Order) ([]*orderDomain.Order, error) {
	orders := make([]*orderDomain.Order, 0, len(docs))
	for _, doc := range docs {
		o, err := toDomain(&doc)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func toItemsDomain(docs []documents.OrderItem) ([]orderDomain.Item, error) {
	items := make([]orderDomain.Item, 0, len(docs))
	for _, model := range docs {
		item, err := toItemDomain(model)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
