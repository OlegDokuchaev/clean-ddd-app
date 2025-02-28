package create_order

import orderDomain "order/internal/domain/order"

type Manager interface {
	Create(order *orderDomain.Order)
}
