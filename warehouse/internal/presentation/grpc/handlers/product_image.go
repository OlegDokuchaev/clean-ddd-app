package handlers

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	productApplication "warehouse/internal/application/product"
	warehousev1 "warehouse/internal/presentation/grpc"
	"warehouse/internal/presentation/grpc/request"
	"warehouse/internal/presentation/grpc/response"
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
	stream grpc.ClientStreamingServer[warehousev1.UploadImageRequest, emptypb.Empty],
) error {
	req, err := stream.Recv()
	if err != nil {
		return response.ParseError(err)
	}

	productID, err := request.ParseUUID(req.GetProductId())
	if err != nil {
		return response.ParseError(err)
	}
	contentType := req.GetContentType()

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

var _ warehousev1.ProductImageServiceServer = (*ProductImageServiceHandler)(nil)
