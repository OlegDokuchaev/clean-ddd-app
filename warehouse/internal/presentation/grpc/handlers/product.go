package handlers

import (
	"context"
	productApplication "warehouse/internal/application/product"
	warehousev1 "warehouse/internal/presentation/grpc"
	"warehouse/internal/presentation/grpc/request"
	"warehouse/internal/presentation/grpc/response"
)

type ProductServiceHandler struct {
	warehousev1.UnimplementedProductServiceServer

	usecase productApplication.UseCase
}

func NewProductServiceHandler(usecase productApplication.UseCase) *ProductServiceHandler {
	return &ProductServiceHandler{
		usecase: usecase,
	}
}

func (h *ProductServiceHandler) CreateProduct(
	ctx context.Context,
	req *warehousev1.CreateProductRequest,
) (*warehousev1.CreateProductResponse, error) {
	data, err := request.ToCreateProductDto(req)
	if err != nil {
		return nil, err
	}

	productID, err := h.usecase.Create(ctx, data)
	if err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToCreateProductResponse(productID), nil
}

var _ warehousev1.ProductServiceServer = (*ProductServiceHandler)(nil)
