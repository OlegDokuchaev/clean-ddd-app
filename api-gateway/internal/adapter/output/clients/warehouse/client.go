package warehouse

import (
	warehouseGRPC "api-gateway/gen/warehouse/v1"
	"api-gateway/internal/adapter/output/clients/response"
	warehouseDto "api-gateway/internal/domain/dtos/warehouse"
	"api-gateway/internal/infrastructure/config"
	warehouseClient "api-gateway/internal/port/output/clients/warehouse"
	"context"
	"io"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClientImpl struct {
	clients *GRPCClients
	sConfig *config.StreamingConfig
}

func NewClient(clients *GRPCClients, sConfig *config.StreamingConfig) warehouseClient.Client {
	return &ClientImpl{
		clients: clients,
		sConfig: sConfig,
	}
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

func (c *ClientImpl) UpdateProductImage(
	ctx context.Context,
	productID uuid.UUID,
	fileReader io.Reader,
	contentType string,
) error {
	// Create stream
	stream, err := c.clients.ProductImage.UpdateImage(ctx)
	if err != nil {
		return response.ParseGRPCError(err)
	}

	// Send product ID and content type
	infoMsg := &warehouseGRPC.UpdateImageRequest{
		Data: &warehouseGRPC.UpdateImageRequest_Info{
			Info: &warehouseGRPC.UpdateImageInfo{
				ProductId:   productID.String(),
				ContentType: contentType,
			},
		},
	}
	if err = stream.Send(infoMsg); err != nil {
		return response.ParseGRPCError(err)
	}

	// Send file data
	buf := make([]byte, c.sConfig.FileChunkSizeBytes)
	for {
		n, err := fileReader.Read(buf)
		if err != nil && err != io.EOF {
			return response.ParseGRPCError(err)
		}
		if n > 0 {
			chunkMsg := &warehouseGRPC.UpdateImageRequest{
				Data: &warehouseGRPC.UpdateImageRequest_ChunkData{
					ChunkData: buf[:n],
				},
			}
			if err = stream.Send(chunkMsg); err != nil {
				return response.ParseGRPCError(err)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return response.ParseGRPCError(err)
		}
	}

	// Close stream
	_, err = stream.CloseAndRecv()
	if err != nil {
		return response.ParseGRPCError(err)
	}

	return nil
}

func (c *ClientImpl) GetProductImage(ctx context.Context, productID uuid.UUID) (io.Reader, string, error) {
	// Create stream
	stream, err := c.clients.ProductImage.GetImage(ctx, &warehouseGRPC.GetImageRequest{
		ProductId: productID.String(),
	})
	if err != nil {
		return nil, "", response.ParseGRPCError(err)
	}

	// Get file stats
	msg, err := stream.Recv()
	if err != nil {
		return nil, "", response.ParseGRPCError(err)
	}

	var contentType string
	switch x := msg.Data.(type) {
	case *warehouseGRPC.GetImageResponse_Info:
		contentType = x.Info.GetContentType()
	default:
		return nil, "", response.ErrInternalServerError
	}

	// Get file data
	fileReader := &productImageGrpcStreamReader{
		stream: stream,
		buffer: nil,
	}

	return fileReader, contentType, nil
}

var _ warehouseClient.Client = (*ClientImpl)(nil)
