package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	createOrderPublisher "order/internal/infrastructure/publisher/saga/create_order"

	"github.com/segmentio/kafka-go"
)

type Processor struct {
	handler *Handler
	reader  *kafka.Reader
	writer  *kafka.Writer
}

func NewProcessor(handler *Handler, reader *kafka.Reader, writer *kafka.Writer) *Processor {
	return &Processor{
		handler: handler,
		reader:  reader,
		writer:  writer,
	}
}

func (p *Processor) readCmdMessage(ctx context.Context) (*createOrderPublisher.CmdMessage, error) {
	msg, err := p.reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading message: %w", err)
	}

	cmdMsg, err := parseCmdMessage(msg.Value)
	if err != nil {
		return nil, fmt.Errorf("error parsing message: %w", err)
	}

	return cmdMsg, nil
}

func (p *Processor) writeResMessage(ctx context.Context, resMessage *ResMessage) error {
	msg, err := json.Marshal(resMessage)
	if err != nil {
		return fmt.Errorf("error serializing response: %w", err)
	}

	kafkaMsg := kafka.Message{
		Value: msg,
	}
	if err = p.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

func (p *Processor) handleMessage(ctx context.Context, cmdMsg *createOrderPublisher.CmdMessage) error {
	log.Printf("Processing message %s of type: %s", cmdMsg.ID, cmdMsg.Name)

	resMsg, err := p.handler.Handle(ctx, cmdMsg)
	if err != nil {
		return fmt.Errorf("error handling message: %w", err)
	}

	if err = p.writeResMessage(ctx, resMsg); err != nil {
		return fmt.Errorf("error sending response: %w", err)
	}

	log.Printf("Message %s successfully processed", cmdMsg.ID)
	return nil
}

func (p *Processor) Process(ctx context.Context) {
	log.Println("Starting command processor")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping command processor")
			return
		default:
			p.processNextMessage(ctx)
		}
	}
}

func (p *Processor) processNextMessage(ctx context.Context) {
	cmdMsg, err := p.readCmdMessage(ctx)
	if err != nil {
		log.Printf("Error reading message: %v", err)
		return
	}

	log.Printf("Received message %s of type: %s", cmdMsg.ID, cmdMsg.Name)

	if err = p.handleMessage(ctx, cmdMsg); err != nil {
		log.Printf("Error processing message %s: %v", cmdMsg.ID, err)
	}
}

func (p *Processor) Close() error {
	log.Println("Closing Kafka connections")

	if err := p.reader.Close(); err != nil {
		return fmt.Errorf("error closing Kafka reader: %w", err)
	}

	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("error closing Kafka writer: %w", err)
	}

	log.Println("Kafka connections successfully closed")
	return nil
}

func parseCmdMessage(kafkaMsg []byte) (*createOrderPublisher.CmdMessage, error) {
	var msg createOrderPublisher.CmdMessage
	if err := json.Unmarshal(kafkaMsg, &msg); err != nil {
		return nil, fmt.Errorf("error deserializing message: %w", err)
	}
	return &msg, nil
}
