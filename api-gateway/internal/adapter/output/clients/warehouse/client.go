package warehouse

import (
	warehouseGRPC "api-gateway/gen/warehouse/v1"
	"api-gateway/internal/adapter/output/clients/response"
	warehouseDto "api-gateway/internal/domain/dtos/warehouse"
	warehouseClient "api-gateway/internal/port/output/clients/warehouse"
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClientImpl struct {
	itemClient    warehouseGRPC.ItemServiceClient
	productClient warehouseGRPC.ProductServiceClient
}

func NewClient(
	itemClient warehouseGRPC.ItemServiceClient,
	productClient warehouseGRPC.ProductServiceClient,
) warehouseClient.Client {
	return &ClientImpl{
		itemClient:    itemClient,
		productClient: productClient,
	}
}

func (c *ClientImpl) ReserveItems(ctx context.Context, items []warehouseDto.ItemInfoDto) error {
	request := toReserveItemRequest(items)

	_, err := c.itemClient.ReserveItem(ctx, request)
	if err != nil {
		return response.ParseGRPCError(err)
	}

	return nil
}

func (c *ClientImpl) ReleaseItems(ctx context.Context, items []warehouseDto.ItemInfoDto) error {
	request := toReleaseItemRequest(items)

	_, err := c.itemClient.ReleaseItem(ctx, request)
	if err != nil {
		return response.ParseGRPCError(err)
	}

	return nil
}

func (c *ClientImpl) CreateProduct(ctx context.Context, data warehouseDto.CreateProductDto) (uuid.UUID, error) {
	request := toCreateProductRequest(data)

	resp, err := c.productClient.CreateProduct(ctx, request)
	if err != nil {
		return uuid.Nil, response.ParseGRPCError(err)
	}

	productID, err := response.ToUUID(resp.ProductId)
	if err != nil {
		return uuid.Nil, err
	}

	return productID, nil
}

func (c *ClientImpl) GetAllItems(ctx context.Context) ([]*warehouseDto.ItemDto, error) {
	resp, err := c.itemClient.GetAllItems(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, response.ParseGRPCError(err)
	}

	items, err := toItems(resp.Items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

var _ warehouseClient.Client = (*ClientImpl)(nil)
