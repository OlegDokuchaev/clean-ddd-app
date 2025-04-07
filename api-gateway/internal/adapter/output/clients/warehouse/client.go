package warehouse

import (
	"api-gateway/internal/adapter/output/clients/response"
	warehouseDto "api-gateway/internal/domain/dtos/warehouse"
	warehouseClient "api-gateway/internal/port/output/clients/warehouse"
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClientImpl struct {
	clients *GRPCClients
}

func NewClient(clients *GRPCClients) warehouseClient.Client {
	return &ClientImpl{clients: clients}
}

func (c *ClientImpl) ReserveItems(ctx context.Context, items []warehouseDto.ItemInfoDto) error {
	request := toReserveItemRequest(items)

	_, err := c.clients.Item.ReserveItem(ctx, request)
	if err != nil {
		return response.ParseGRPCError(err)
	}

	return nil
}

func (c *ClientImpl) ReleaseItems(ctx context.Context, items []warehouseDto.ItemInfoDto) error {
	request := toReleaseItemRequest(items)

	_, err := c.clients.Item.ReleaseItem(ctx, request)
	if err != nil {
		return response.ParseGRPCError(err)
	}

	return nil
}

func (c *ClientImpl) CreateProduct(ctx context.Context, data warehouseDto.CreateProductDto) (uuid.UUID, error) {
	request := toCreateProductRequest(data)

	resp, err := c.clients.Product.CreateProduct(ctx, request)
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
	resp, err := c.clients.Item.GetAllItems(ctx, &emptypb.Empty{})
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
