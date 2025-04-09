package warehouse_response

import (
	warehouseDto "api-gateway/internal/domain/dtos/warehouse"
)

func ToItemResponse(item *warehouseDto.ItemDto) ItemResponse {
	return ItemResponse{
		ItemID:  item.ItemID,
		Count:   item.Count,
		Product: ToProductSchema(item.Product),
		Version: item.Version.String(),
	}
}

func ToItemsResponse(items []*warehouseDto.ItemDto) ItemsResponse {
	result := make([]ItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, ToItemResponse(item))
	}
	return ItemsResponse{Items: result}
}

func ToProductSchema(product warehouseDto.ProductDto) ProductSchema {
	return ProductSchema{
		ProductID: product.ProductID,
		Name:      product.Name,
		Price:     product.Price,
		Created:   product.Created,
	}
}
