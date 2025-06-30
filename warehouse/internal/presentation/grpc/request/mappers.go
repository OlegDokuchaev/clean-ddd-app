package request

import (
	itemApplication "warehouse/internal/application/item"
	productApplication "warehouse/internal/application/product"
	warehousev1 "warehouse/internal/presentation/grpc"
)

func toItemInfoDto(req *warehousev1.ItemInfo) (itemApplication.ItemDto, error) {
	productID, err := ParseUUID(req.ProductId)
	if err != nil {
		return itemApplication.ItemDto{}, err
	}

	return itemApplication.ItemDto{
		ProductID: productID,
		Count:     int(req.Count),
	}, nil
}

func toItemsInfoDto(req []*warehousev1.ItemInfo) ([]itemApplication.ItemDto, error) {
	var items []itemApplication.ItemDto
	for _, item := range req {
		itemInfo, err := toItemInfoDto(item)
		if err != nil {
			return nil, err
		}
		items = append(items, itemInfo)
	}
	return items, nil
}

func ToCreateProductDto(req *warehousev1.CreateProductRequest) (productApplication.CreateDto, error) {
	return productApplication.CreateDto{
		Name:  req.Name,
		Price: ParseDecimal(req.Price),
	}, nil
}

func ToReserveItemDto(req *warehousev1.ReserveItemRequest) (itemApplication.ReserveDto, error) {
	var data itemApplication.ReserveDto

	items, err := toItemsInfoDto(req.Items)
	if err != nil {
		return data, err
	}
	data.Items = items

	return data, nil
}

func ToReleaseItemDto(req *warehousev1.ReleaseItemRequest) (itemApplication.ReleaseDto, error) {
	var data itemApplication.ReleaseDto

	items, err := toItemsInfoDto(req.Items)
	if err != nil {
		return data, err
	}
	data.Items = items

	return data, nil
}
