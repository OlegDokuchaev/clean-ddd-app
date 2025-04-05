package warehouse

import (
	warehouseGRPC "api-gateway/gen/warehouse/v1"
	"api-gateway/internal/adapter/output/clients/response"
	warehouseDto "api-gateway/internal/domain/dtos/warehouse"
)

func toItems(protoItems []*warehouseGRPC.Item) ([]*warehouseDto.ItemDto, error) {
	items := make([]*warehouseDto.ItemDto, 0, len(protoItems))
	for _, protoItem := range protoItems {
		item, err := toItem(protoItem)
		if err == nil {
			items = append(items, item)
		}
	}
	return items, nil
}

func toItem(protoItem *warehouseGRPC.Item) (*warehouseDto.ItemDto, error) {
	itemID, err := response.ToUUID(protoItem.ItemId)
	if err != nil {
		return nil, err
	}

	versionID, err := response.ToUUID(protoItem.Version)
	if err != nil {
		return nil, err
	}

	product, err := toProduct(protoItem.Product)
	if err != nil {
		return nil, err
	}

	return &warehouseDto.ItemDto{
		ItemID:  itemID,
		Count:   int(protoItem.Count),
		Product: *product,
		Version: versionID,
	}, nil
}

func toProduct(protoProduct *warehouseGRPC.Product) (*warehouseDto.ProductDto, error) {
	productID, err := response.ToUUID(protoProduct.ProductId)
	if err != nil {
		return nil, err
	}

	return &warehouseDto.ProductDto{
		ProductID: productID,
		Name:      protoProduct.Name,
		Price:     response.ToDecimal(protoProduct.Price),
		Created:   protoProduct.Created.AsTime(),
	}, nil
}
