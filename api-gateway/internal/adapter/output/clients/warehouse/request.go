package warehouse

import (
	warehouseGRPC "api-gateway/gen/warehouse/v1"
	warehouseDto "api-gateway/internal/domain/dtos/warehouse"
)

func toItemInfo(item warehouseDto.ItemInfoDto) *warehouseGRPC.ItemInfo {
	return &warehouseGRPC.ItemInfo{
		ProductId: item.ProductID.String(),
		Count:     int32(item.Count),
	}
}

func toItemsInfo(items []warehouseDto.ItemInfoDto) []*warehouseGRPC.ItemInfo {
	result := make([]*warehouseGRPC.ItemInfo, 0, len(items))
	for _, item := range items {
		result = append(result, toItemInfo(item))
	}
	return result
}

func toReserveItemRequest(items []warehouseDto.ItemInfoDto) *warehouseGRPC.ReserveItemRequest {
	return &warehouseGRPC.ReserveItemRequest{
		Items: toItemsInfo(items),
	}
}

func toReleaseItemRequest(items []warehouseDto.ItemInfoDto) *warehouseGRPC.ReleaseItemRequest {
	return &warehouseGRPC.ReleaseItemRequest{
		Items: toItemsInfo(items),
	}
}

func toCreateProductRequest(data warehouseDto.CreateProductDto) *warehouseGRPC.CreateProductRequest {
	return &warehouseGRPC.CreateProductRequest{
		Name:  data.Name,
		Price: data.Price.InexactFloat64(),
	}
}
