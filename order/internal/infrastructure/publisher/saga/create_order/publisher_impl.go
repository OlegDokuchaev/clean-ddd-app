package create_order

import (
	"context"
	"encoding/json"
	createOrder "order/internal/application/order/saga/create_order"

	otelkafkakonsumer "github.com/Trendyol/otel-kafka-konsumer"

	"github.com/segmentio/kafka-go"
)

type PublisherImpl struct {
	warehouseWriter *otelkafkakonsumer.Writer
	orderWriter     *otelkafkakonsumer.Writer
	courierWriter   *otelkafkakonsumer.Writer
}

func NewPublisher(
	warehouseWriter *otelkafkakonsumer.Writer,
	orderWriter *otelkafkakonsumer.Writer,
	courierWriter *otelkafkakonsumer.Writer,
) *PublisherImpl {
	return &PublisherImpl{
		warehouseWriter: warehouseWriter,
		orderWriter:     orderWriter,
		courierWriter:   courierWriter,
	}
}

func (p *PublisherImpl) PublishReserveItemsCmd(ctx context.Context, cmd createOrder.ReserveItemsCmd) error {
	cmdMsg := NewCmdMessage(ReserveItemsCmdName, cmd)
	return publishMessage(ctx, p.warehouseWriter, cmdMsg)
}

func (p *PublisherImpl) PublishReleaseItemsCmd(ctx context.Context, cmd createOrder.ReleaseItemsCmd) error {
	cmdMsg := NewCmdMessage(ReleaseItemsCmdName, cmd)
	return publishMessage(ctx, p.warehouseWriter, cmdMsg)
}

func (p *PublisherImpl) PublishCancelOutOfStockCmd(ctx context.Context, cmd createOrder.CancelOutOfStockCmd) error {
	cmdMsg := NewCmdMessage(CancelOutOfStockCmdName, cmd)
	return publishMessage(ctx, p.orderWriter, cmdMsg)
}

func (p *PublisherImpl) PublishAssignCourierCmd(ctx context.Context, cmd createOrder.AssignCourierCmd) error {
	cmdMsg := NewCmdMessage(AssignCourierCmdName, cmd)
	return publishMessage(ctx, p.courierWriter, cmdMsg)
}

func (p *PublisherImpl) PublishBeginDeliveryCmd(ctx context.Context, cmd createOrder.BeginDeliveryCmd) error {
	cmdMsg := NewCmdMessage(BeginDeliveryCmdName, cmd)
	return publishMessage(ctx, p.orderWriter, cmdMsg)
}

func (p *PublisherImpl) PublishCancelCourierNotFoundCmd(ctx context.Context, cmd createOrder.CancelCourierNotFoundCmd) error {
	cmdMsg := NewCmdMessage(CancelCourierNotFoundCmdName, cmd)
	return publishMessage(ctx, p.orderWriter, cmdMsg)
}

func encodeMessage(msg CmdMessage) ([]byte, error) {
	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, parseError(err)
	}
	return buf, nil
}

func publishMessage(ctx context.Context, writer *otelkafkakonsumer.Writer, msg CmdMessage) error {
	value, err := encodeMessage(msg)
	if err != nil {
		return err
	}

	kafkaMsg := kafka.Message{Value: value}

	ctx = writer.TraceConfig.Propagator.Extract(ctx, otelkafkakonsumer.NewMessageCarrier(&kafkaMsg))

	err = writer.WriteMessage(ctx, kafkaMsg)
	return parseError(err)
}

var _ createOrder.Publisher = (*PublisherImpl)(nil)
