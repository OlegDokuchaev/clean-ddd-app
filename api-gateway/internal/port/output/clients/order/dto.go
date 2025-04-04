package order

import orderUseCase "api-gateway/internal/domain/usecases/order"

type CreateDto struct {
	CustomerID string
	Address    string
	Items      []orderUseCase.ItemDto
}
