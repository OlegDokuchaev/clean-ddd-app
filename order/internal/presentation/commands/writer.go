package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"order/internal/infrastructure/logger"
	createOrderConsumer "order/internal/presentation/saga/create_order"

	"github.com/segmentio/kafka-go"
)

type Writer interface {
	Write(ctx context.Context, res *createOrderConsumer.ResMessage) error
}

type WriterImpl struct {
	writer *kafka.Writer
	logger logger.Logger
}

func NewWriter(writer *kafka.Writer, logger logger.Logger) *WriterImpl {
	return &WriterImpl{
		writer: writer,
		logger: logger,
	}
}

func (w *WriterImpl) Write(ctx context.Context, res *createOrderConsumer.ResMessage) error {
	if res == nil {
		return nil // Если ответ не требуется, просто выходим
	}

	msg, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("error serializing response: %w", err)
	}

	kafkaMsg := kafka.Message{
		Value: msg,
	}
	if err = w.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	w.logger.Printf("Response sent: %s, type: %s", res.ID, res.Name)
	return nil
}

var _ Writer = (*WriterImpl)(nil)
