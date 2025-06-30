package request

import (
	orderUsecase "order/internal/application/order/usecase"
	orderDomain "order/internal/domain/order"
	orderv1 "order/internal/presentation/grpc"
)

func ToItem(item *orderv1.OrderItem) (orderDomain.Item, error) {
	var data orderDomain.Item

	productID, err := ParseUUID(item.ProductId)
	if err != nil {
		return data, err
	}

	data.ProductID = productID
	data.Price = ParseDecimal(item.Price)
	data.Count = int(item.Count)

	return data, nil
}

func ToItems(items []*orderv1.OrderItem) ([]orderDomain.Item, error) {
	orderItems := make([]orderDomain.Item, 0, len(items))
	for _, item := range items {
		orderItem, err := ToItem(item)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}

func ToCreateDto(req *orderv1.CreateOrderRequest) (orderUsecase.CreateDto, error) {
	var data orderUsecase.CreateDto

	customerID, err := ParseUUID(req.CustomerId)
	if err != nil {
		return data, err
	}

	items, err := ToItems(req.Items)
	if err != nil {
		return data, err
	}

	data.CustomerID = customerID
	data.Address = req.Address
	data.Items = items

	return data, nil
}
