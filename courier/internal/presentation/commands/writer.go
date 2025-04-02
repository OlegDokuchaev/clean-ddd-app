package commands

import (
	"context"
	"courier/internal/infrastructure/logger"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Writer interface {
	Write(ctx context.Context, res *ResMessage) error
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

func (w *WriterImpl) log(level logger.Level, action, message string, extraFields map[string]any) {
	fields := map[string]any{
		"component": "command_writer",
		"action":    action,
		"topic":     w.writer.Topic,
	}
	for k, v := range extraFields {
		fields[k] = v
	}

	w.logger.Log(level, message, fields)
}

func (w *WriterImpl) Write(ctx context.Context, res *ResMessage) error {
	if res == nil {
		return nil
	}

	// Serialize the response
	msg, err := json.Marshal(res)
	if err != nil {
		w.log(logger.Error, "serialize_error", "Failed to serialize response", map[string]any{
			"response_id": res.ID,
			"error":       err.Error(),
		})
		return fmt.Errorf("error serializing response: %w", err)
	}

	kafkaMsg := kafka.Message{
		Value: msg,
	}

	// Write the message to Kafka
	if err = w.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		w.log(logger.Error, "kafka_write_error", "Failed to send response to Kafka", map[string]any{
			"response_id": res.ID,
			"error":       err.Error(),
		})
		return fmt.Errorf("error sending message: %w", err)
	}

	w.log(logger.Info, "response_sent", "Command response sent to Kafka", map[string]any{
		"response": res,
	})
	return nil
}

var _ Writer = (*WriterImpl)(nil)
