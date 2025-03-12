package outbox

import (
	"context"
	"time"
	outboxDomain "warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"

	"github.com/segmentio/kafka-go"
)

type PublisherImpl struct {
	config *Config
	writer *kafka.Writer
}

func NewPublisher(config *Config, writer *kafka.Writer) *PublisherImpl {
	return &PublisherImpl{
		config: config,
		writer: writer,
	}
}

func (p *PublisherImpl) Publish(ctx context.Context, message *outboxDomain.Message) error {
	topic, err := p.getTopicByMessage(message)
	if err != nil {
		return err
	}

	kafkaMsg := kafka.Message{
		Topic: topic,
		Value: message.Payload,
		Time:  time.Now(),
	}

	err = p.writer.WriteMessages(ctx, kafkaMsg)
	return parseError(err)
}

func (p *PublisherImpl) getTopicByMessage(message *outboxDomain.Message) (string, error) {
	switch message.Type {
	case productDomain.CreatedEventName:
		return p.config.ProductTopic, nil

	default:
		return "", ErrInvalidOutboxMessage
	}
}

var _ outboxDomain.Publisher = (*PublisherImpl)(nil)
