package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	createOrderPublisher "order/internal/infrastructure/publisher/saga/create_order"

	"github.com/segmentio/kafka-go"
)

type Reader interface {
	Read(ctx context.Context) (*createOrderPublisher.CmdMessage, error)
	Close() error
}

type ReaderImpl struct {
	reader       *kafka.Reader
	commandChan  chan *createOrderPublisher.CmdMessage
	errorChan    chan error
	shutdownChan chan struct{}
}

func NewReader(reader *kafka.Reader) *ReaderImpl {
	return &ReaderImpl{
		reader:       reader,
		commandChan:  make(chan *createOrderPublisher.CmdMessage, 1),
		errorChan:    make(chan error, 1),
		shutdownChan: make(chan struct{}),
	}
}

func (r *ReaderImpl) Start(ctx context.Context) {
	go r.readCommands(ctx)
}

func (r *ReaderImpl) readCommands(ctx context.Context) {
	defer close(r.commandChan)
	defer close(r.errorChan)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Context canceled, stopping command reader")
			return
		case <-r.shutdownChan:
			log.Printf("Shutdown requested, stopping command reader")
			return
		default:
			msg, err := r.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					continue
				}

				select {
				case r.errorChan <- fmt.Errorf("error reading message: %w", err):
				default:
					log.Printf("Error channel full, dropping error: %v", err)
				}
				continue
			}

			cmdMsg, err := parseCommandMessage(msg.Value)
			if err != nil {
				select {
				case r.errorChan <- fmt.Errorf("error parsing message: %w", err):
				default:
					log.Printf("Error channel full, dropping error: %v", err)
				}
				continue
			}

			select {
			case r.commandChan <- cmdMsg:
				log.Printf("Command received: %s, type: %s", cmdMsg.ID, cmdMsg.Name)
			case <-ctx.Done():
				return
			case <-r.shutdownChan:
				return
			}
		}
	}
}

func (r *ReaderImpl) Read(ctx context.Context) (*createOrderPublisher.CmdMessage, error) {
	select {
	case cmd, ok := <-r.commandChan:
		if !ok {
			return nil, fmt.Errorf("command channel closed")
		}
		return cmd, nil
	case err, ok := <-r.errorChan:
		if !ok {
			return nil, fmt.Errorf("error channel closed")
		}
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *ReaderImpl) Close() error {
	close(r.shutdownChan)
	return r.reader.Close()
}

var _ Reader = (*ReaderImpl)(nil)

func parseCommandMessage(data []byte) (*createOrderPublisher.CmdMessage, error) {
	var msg createOrderPublisher.CmdMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("error deserializing message: %w", err)
	}
	return &msg, nil
}
