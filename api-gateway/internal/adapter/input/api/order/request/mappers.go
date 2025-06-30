package order_request

import orderDto "api-gateway/internal/domain/dtos/order"

func ToOrderCreateDto(request *CreateRequest) orderDto.CreateDto {
	return orderDto.CreateDto{
		Address: request.Address,
		Items:   ToItemDtoList(request.Items),
	}
}

func ToItemDtoList(schemas []*ItemSchema) []orderDto.ItemDto {
	items := make([]orderDto.ItemDto, 0, len(schemas))
	for _, schema := range schemas {
		items = append(items, ToItemDto(schema))
	}
	return items
}

func ToItemDto(schema *ItemSchema) orderDto.ItemDto {
	return orderDto.ItemDto{
		ProductID: schema.ProductID,
		Price:     schema.Price,
		Count:     schema.Count,
	}
}
