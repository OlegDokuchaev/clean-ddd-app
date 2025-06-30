package handlers

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	itemApplication "warehouse/internal/application/item"
	warehousev1 "warehouse/internal/presentation/grpc"
	"warehouse/internal/presentation/grpc/request"
	"warehouse/internal/presentation/grpc/response"
)

type ItemServiceHandler struct {
	warehousev1.UnimplementedItemServiceServer

	usecase itemApplication.UseCase
}

func NewItemServiceHandler(usecase itemApplication.UseCase) *ItemServiceHandler {
	return &ItemServiceHandler{
		usecase: usecase,
	}
}

func (h *ItemServiceHandler) ReserveItem(
	ctx context.Context,
	req *warehousev1.ReserveItemRequest,
) (*emptypb.Empty, error) {
	data, err := request.ToReserveItemDto(req)
	if err != nil {
		return nil, err
	}

	err = h.usecase.Reserve(ctx, data)
	if err != nil {
		return nil, response.ParseError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ItemServiceHandler) ReleaseItem(
	ctx context.Context,
	req *warehousev1.ReleaseItemRequest,
) (*emptypb.Empty, error) {
	data, err := request.ToReleaseItemDto(req)
	if err != nil {
		return nil, err
	}

	err = h.usecase.Release(ctx, data)
	if err != nil {
		return nil, response.ParseError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ItemServiceHandler) GetAllItems(
	ctx context.Context,
	_ *emptypb.Empty,
) (*warehousev1.GetAllItemsResponse, error) {
	data, err := h.usecase.GetAll(ctx)
	if err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToGetAllItemsResponse(data), nil
}

var _ warehousev1.ItemServiceServer = (*ItemServiceHandler)(nil)
