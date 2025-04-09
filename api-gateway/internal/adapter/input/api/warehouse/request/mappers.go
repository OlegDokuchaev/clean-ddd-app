package warehouse_request

import (
	warehouseDto "api-gateway/internal/domain/dtos/warehouse"
)

func ToItemInfoDto(info *ItemInfoSchema) warehouseDto.ItemInfoDto {
	return warehouseDto.ItemInfoDto{
		ProductID: info.ProductID,
		Count:     info.Count,
	}
}

func ToItemInfoDtoList(items []*ItemInfoSchema) []warehouseDto.ItemInfoDto {
	result := make([]warehouseDto.ItemInfoDto, 0, len(items))
	for _, item := range items {
		result = append(result, ToItemInfoDto(item))
	}
	return result
}

func ToCreateProductDto(req *CreateProductRequest) warehouseDto.CreateProductDto {
	return warehouseDto.CreateProductDto{
		Name:  req.Name,
		Price: req.Price,
	}
}
