package warehouse

import (
	warehouseGRPC "api-gateway/gen/warehouse/v1"
	"api-gateway/internal/adapter/output/clients/response"
	warehouseDto "api-gateway/internal/domain/dtos/warehouse"
	warehouseClient "api-gateway/internal/port/output/clients/warehouse"
	"context"
	"io"

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

func (c *ClientImpl) UpdateProductImage(
	ctx context.Context,
	productID uuid.UUID,
	fileReader io.Reader,
	contentType string,
) error {
	// Create stream
	stream, err := c.clients.ProductImage.UploadImage(ctx)
	if err != nil {
		return response.ParseGRPCError(err)
	}

	// Send product ID and content type
	if err = stream.Send(&warehouseGRPC.UploadImageRequest{
		ProductId:   productID.String(),
		ContentType: contentType,
	}); err != nil {
		return response.ParseGRPCError(err)
	}

	// Send file data
	buf := make([]byte, 32*1024)
	for {
		n, readErr := fileReader.Read(buf)
		if readErr != nil && readErr != io.EOF {
			return response.ParseGRPCError(err)
		}
		if n > 0 {
			if err := stream.Send(&warehouseGRPC.UploadImageRequest{
				ChunkData: buf[:n],
			}); err != nil {
				return response.ParseGRPCError(err)
			}
		}
		if readErr == io.EOF {
			break
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
		return nil, "", response.ParseGRPCError(err)
	}

	// Get file data
	fileReader := &productImageGrpcStreamReader{
		stream: stream,
		buffer: nil,
	}

	return fileReader, contentType, nil
}

var _ warehouseClient.Client = (*ClientImpl)(nil)
