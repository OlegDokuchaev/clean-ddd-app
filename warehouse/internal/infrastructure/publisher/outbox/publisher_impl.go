package outbox

import (
	"context"
	"encoding/json"
	outboxDomain "warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"

	"github.com/segmentio/kafka-go"
)

type PublisherImpl struct {
	productWriter *kafka.Writer
}

func NewPublisher(productWriter *kafka.Writer) *PublisherImpl {
	return &PublisherImpl{productWriter: productWriter}
}

func (p *PublisherImpl) Publish(ctx context.Context, message *outboxDomain.Message) error {
	writer, err := p.getWriterByMessage(message)
	if err != nil {
		return err
	}
	return publishMessage(ctx, writer, message)
}

func (p *PublisherImpl) getWriterByMessage(message *outboxDomain.Message) (*kafka.Writer, error) {
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

func publishMessage(ctx context.Context, writer *kafka.Writer, message *outboxDomain.Message) error {
	value, err := encodeMessage(message)
	if err != nil {
		return err
	}

	kafkaMsg := kafka.Message{
		Value: value,
	}
	err = writer.WriteMessages(ctx, kafkaMsg)
	return parseError(err)
}

var _ outboxDomain.Publisher = (*PublisherImpl)(nil)
