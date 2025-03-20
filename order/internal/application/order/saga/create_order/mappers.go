package create_order

import (
	orderDomain "order/internal/domain/order"
)

func domainItemToOrderItem(domainItem orderDomain.Item) OrderItem {
	return OrderItem{
		ProductID: domainItem.ProductID,
		Count:     domainItem.Count,
	}
}

func domainItemsToOrderItems(domainItems []orderDomain.Item) []OrderItem {
	orderItems := make([]OrderItem, len(domainItems))
	for i, item := range domainItems {
		orderItems[i] = domainItemToOrderItem(item)
	}
	return orderItems
}
