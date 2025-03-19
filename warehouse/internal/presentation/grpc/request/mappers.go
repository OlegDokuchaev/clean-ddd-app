package request

import (
	itemApplication "warehouse/internal/application/item"
	productApplication "warehouse/internal/application/product"
	warehousev1 "warehouse/internal/presentation/grpc"
)

func ToCreateProductDto(req *warehousev1.CreateProductRequest) (productApplication.CreateDto, error) {
	return productApplication.CreateDto{
		Name:  req.Name,
		Price: parseDecimal(req.Price),
	}, nil
}

func ToCreateItemDto(req *warehousev1.CreateItemRequest) (itemApplication.CreateDto, error) {
	var data itemApplication.CreateDto

	productID, err := parseUUID(req.ProductId)
	if err != nil {
		return data, err
	}

	data.ProductID = productID
	data.Count = int(req.Count)

	return data, nil
}

func ToReserveItemDto(req *warehousev1.ReserveItemRequest) (itemApplication.ReserveDto, error) {
	var data itemApplication.ReserveDto

	itemID, err := parseUUID(req.ItemId)
	if err != nil {
		return data, err
	}

	data.ItemID = itemID
	data.Count = int(req.Count)

	return data, nil
}

func ToReleaseItemDto(req *warehousev1.ReleaseItemRequest) (itemApplication.ReleaseDto, error) {
	var data itemApplication.ReleaseDto

	itemID, err := parseUUID(req.ItemId)
	if err != nil {
		return data, err
	}

	data.ItemID = itemID
	data.Count = int(req.Count)

	return data, nil
}
