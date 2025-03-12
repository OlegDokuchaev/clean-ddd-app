package outbox

import (
	"context"
	"encoding/json"
	"time"
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
	switch message.Type {
	case productDomain.CreatedEventName:
		return p.productWriter, nil

	default:
		return nil, ErrInvalidOutboxMessage
	}
}

func encodeMessage(message *outboxDomain.Message) ([]byte, error) {
	buf, err := json.Marshal(message)
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
		Time:  time.Now(),
	}
	err = writer.WriteMessages(ctx, kafkaMsg)
	return parseError(err)
}

var _ outboxDomain.Publisher = (*PublisherImpl)(nil)
