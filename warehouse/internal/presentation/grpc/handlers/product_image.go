package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	productApplication "warehouse/internal/application/product"
	warehousev1 "warehouse/internal/presentation/grpc"
	"warehouse/internal/presentation/grpc/request"
	"warehouse/internal/presentation/grpc/response"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductImageServiceHandler struct {
	warehousev1.UnimplementedProductImageServiceServer

	usecase productApplication.ImageUseCase
}

func NewProductImageServiceHandler(usecase productApplication.ImageUseCase) *ProductImageServiceHandler {
	return &ProductImageServiceHandler{
		usecase: usecase,
	}
}

func (h *ProductImageServiceHandler) UploadImage(
	stream grpc.ClientStreamingServer[warehousev1.UpdateImageRequest, emptypb.Empty],
) error {
	// Get file info
	var productID uuid.UUID
	var contentType string

	req, err := stream.Recv()
	if err != nil {
		return response.ParseError(err)
	}

	switch x := req.Data.(type) {
	case *warehousev1.UpdateImageRequest_Info:
		productID, err = request.ParseUUID(x.Info.GetProductId())
		if err != nil {
			return response.ParseError(err)
		}
		contentType = x.Info.GetContentType()
	default:
		return response.ErrInternalError
	}

	// Get file data
	pr, pw := io.Pipe()
	defer pr.Close()

	errChan := make(chan error, 1)

	go func() {
		defer pw.Close()

		for {
			req, err = stream.Recv()
			if err == io.EOF {
				errChan <- nil
				return
			}
			if err != nil {
				errChan <- fmt.Errorf("error receiving stream chunk: %w", err)
				return
			}

			if data := req.GetChunkData(); len(data) > 0 {
				if _, wErr := pw.Write(data); wErr != nil {
					errChan <- fmt.Errorf("error writing chunk to pipe: %w", wErr)
					return
				}
			}
		}
	}()

	err = h.usecase.UpdateByID(stream.Context(), productID, pr, contentType)
	if err != nil {
		return response.ParseError(err)
	}

	if goroutineErr := <-errChan; goroutineErr != nil {
		return response.ParseError(err)
	}

	return stream.SendAndClose(&emptypb.Empty{})
}

func (h *ProductImageServiceHandler) GetImage(
	req *warehousev1.GetImageRequest,
	stream grpc.ServerStreamingServer[warehousev1.GetImageResponse],
) error {
	productID, err := request.ParseUUID(req.GetProductId())
	if err != nil {
		return response.ParseError(err)
	}

	fileReader, contentType, err := h.usecase.GetByID(stream.Context(), productID)
	if err != nil {
		return response.ParseError(err)
	}

	// Send file info
	headerMsg := &warehousev1.GetImageResponse{
		Data: &warehousev1.GetImageResponse_Info{
			Info: &warehousev1.GetImageInfo{
				ContentType: contentType,
			},
		},
	}
	if err = stream.Send(headerMsg); err != nil {
		return response.ParseError(err)
	}

	// Send file data
	buf := make([]byte, 32*1024)
	for {
		n, readErr := fileReader.Read(buf)
		if n > 0 {
			chunkMsg := &warehousev1.GetImageResponse{
				Data: &warehousev1.GetImageResponse_ChunkData{
					ChunkData: buf[:n],
				},
			}
			if err = stream.Send(chunkMsg); err != nil {
				return response.ParseError(err)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return response.ParseError(err)
		}
	}

	return nil
}

var _ warehousev1.ProductImageServiceServer = (*ProductImageServiceHandler)(nil)
