package response

import (
	itemDomain "warehouse/internal/domain/item"
	productDomain "warehouse/internal/domain/product"
	warehousev1 "warehouse/internal/presentation/grpc"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		Version: item.Version.String(),
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

func ToGetAllItemsResponse(items []*itemDomain.Item) *warehousev1.GetAllItemsResponse {
	return &warehousev1.GetAllItemsResponse{
		Items: toItemsResponse(items),
	}
}
