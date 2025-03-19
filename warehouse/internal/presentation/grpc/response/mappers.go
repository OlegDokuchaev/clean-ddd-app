package response

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	itemDomain "warehouse/internal/domain/item"
	productDomain "warehouse/internal/domain/product"
	warehousev1 "warehouse/internal/presentation/grpc"
)

func toProductResponse(product *productDomain.Product) *warehousev1.Product {
	return &warehousev1.Product{
		ProductId: product.ID.String(),
		Name:      product.Name,
		Price:     product.Price.InexactFloat64(),
		Created:   timestamppb.New(product.Created),
	}
}

func toItemResponse(item *itemDomain.Item) *warehousev1.Item {
	return &warehousev1.Item{
		ItemId:  item.ID.String(),
		Count:   int32(item.Count),
		Product: toProductResponse(item.Product),
	}
}

func toItemsResponse(items []*itemDomain.Item) []*warehousev1.Item {
	itemsResponse := make([]*warehousev1.Item, 0, len(items))
	for _, item := range items {
		itemsResponse = append(itemsResponse, toItemResponse(item))
	}
	return itemsResponse
}

func ToCreateProductResponse(productID uuid.UUID) *warehousev1.CreateProductResponse {
	return &warehousev1.CreateProductResponse{
		ProductId: productID.String(),
	}
}

func ToCreateItemResponse(itemID uuid.UUID) *warehousev1.CreateItemResponse {
	return &warehousev1.CreateItemResponse{
		ItemId: itemID.String(),
	}
}

func ToGetAllItemsResponse(items []*itemDomain.Item) *warehousev1.GetAllItemsResponse {
	return &warehousev1.GetAllItemsResponse{
		Items: toItemsResponse(items),
	}
}
