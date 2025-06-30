package warehouse

import (
	warehouseGRPC "api-gateway/gen/warehouse/v1"
	"fmt"

	"google.golang.org/grpc"
)

type productImageGrpcStreamReader struct {
	stream grpc.ServerStreamingClient[warehouseGRPC.GetImageResponse]
	buffer []byte
}

func (r *productImageGrpcStreamReader) Read(p []byte) (int, error) {
	for len(r.buffer) == 0 {
		msg, err := r.stream.Recv()
		if err != nil {
			return 0, err
		}

		switch x := msg.Data.(type) {
		case *warehouseGRPC.GetImageResponse_ChunkData:
			r.buffer = x.ChunkData
		default:
			return 0, fmt.Errorf("unexpected message type in stream")
		}
	}

	n := copy(p, r.buffer)
	r.buffer = r.buffer[n:]
	return n, nil
}
