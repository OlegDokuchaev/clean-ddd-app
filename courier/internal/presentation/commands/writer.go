package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Writer interface {
	Write(ctx context.Context, res *ResMessage) error
	Close() error
}

type WriterImpl struct {
	writer *kafka.Writer
}

func NewWriter(writer *kafka.Writer) *WriterImpl {
	return &WriterImpl{writer: writer}
}

func (w *WriterImpl) Write(ctx context.Context, res *ResMessage) error {
	if res == nil {
		return nil
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

	log.Printf("Response sent: %s, type: %s", res.ID, res.Name)
	return nil
}

func (w *WriterImpl) Close() error {
	return w.writer.Close()
}

var _ Writer = (*WriterImpl)(nil)
