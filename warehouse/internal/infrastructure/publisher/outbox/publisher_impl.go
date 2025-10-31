package outbox

import (
	"context"
	"encoding/json"
	otelkafkakonsumer "github.com/Trendyol/otel-kafka-konsumer"
	outboxDomain "warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"

	"github.com/segmentio/kafka-go"
)

type PublisherImpl struct {
	productWriter *otelkafkakonsumer.Writer
}

func NewPublisher(productWriter *otelkafkakonsumer.Writer) *PublisherImpl {
	return &PublisherImpl{productWriter: productWriter}
}

func (p *PublisherImpl) Publish(ctx context.Context, message *outboxDomain.Message) error {
	writer, err := p.getWriterByMessage(message)
	if err != nil {
		return err
	}
	return publishMessage(ctx, writer, message)
}

func (p *PublisherImpl) getWriterByMessage(message *outboxDomain.Message) (*otelkafkakonsumer.Writer, error) {
	switch message.Name {
	case productDomain.CreatedEventName:
		return p.productWriter, nil

	default:
		return nil, ErrInvalidOutboxMessage
	}
}

func encodeMessage(message *outboxDomain.Message) ([]byte, error) {
	value := KafkaMessageValue{
		ID:      message.ID,
		Name:    message.Name,
		Payload: json.RawMessage(message.Payload),
	}
	buf, err := json.Marshal(value)
	if err != nil {
		return nil, parseError(err)
	}
	return buf, nil
}

func publishMessage(ctx context.Context, writer *otelkafkakonsumer.Writer, message *outboxDomain.Message) error {
	value, err := encodeMessage(message)
	if err != nil {
		return err
	}

	kafkaMsg := kafka.Message{Value: value}

	ctx = writer.TraceConfig.Propagator.Extract(ctx, otelkafkakonsumer.NewMessageCarrier(&kafkaMsg))

	err = writer.WriteMessage(ctx, kafkaMsg)
	return parseError(err)
}

var _ outboxDomain.Publisher = (*PublisherImpl)(nil)
