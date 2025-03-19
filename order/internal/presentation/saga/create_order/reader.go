package create_order

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Reader interface {
	Start(ctx context.Context)
	Read(ctx context.Context) (*ResMessage, error)
	Stop()
}

type ReaderImpl struct {
	reader       *kafka.Reader
	messageChan  chan *ResMessage
	errorChan    chan error
	shutdownChan chan struct{}
}

func NewReader(reader *kafka.Reader) *ReaderImpl {
	return &ReaderImpl{
		reader:       reader,
		shutdownChan: make(chan struct{}),
		messageChan:  make(chan *ResMessage, 1),
		errorChan:    make(chan error, 1),
	}
}

func (c *ReaderImpl) Start(ctx context.Context) {
	go c.consumeMessages(ctx)
}

func (c *ReaderImpl) Stop() {
	log.Printf("Stopping reader...")
	close(c.shutdownChan)
}

func (c *ReaderImpl) consumeMessages(ctx context.Context) {
	defer close(c.messageChan)
	defer close(c.errorChan)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Context canceled, stopping consumer")
			return
		case <-c.shutdownChan:
			log.Printf("Shutdown requested, stopping consumer")
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)

			if err != nil {
				if ctx.Err() != nil {
					return
				}

				// Check if shutdown was requested before reporting error
				select {
				case <-c.shutdownChan:
					return
				default:
					select {
					case c.errorChan <- fmt.Errorf("error reading message: %w", err):
					default:
						// Error channel is full, logging error
						log.Printf("Error channel full, dropping error: %v", err)
					}
				}
				continue
			}

			resMsg, err := c.parseMessage(msg.Value)
			if err != nil {
				select {
				case <-c.shutdownChan:
					return
				default:
					select {
					case c.errorChan <- fmt.Errorf("error parsing message: %w", err):
					default:
						// Error channel is full, logging error
						log.Printf("Error channel full, dropping error: %v", err)
					}
				}
				continue
			}

			select {
			case c.messageChan <- resMsg:
			case <-ctx.Done():
				return
			case <-c.shutdownChan:
				return
			}
		}
	}
}

func (c *ReaderImpl) Read(ctx context.Context) (*ResMessage, error) {
	select {
	case msg, ok := <-c.messageChan:
		if !ok {
			return nil, fmt.Errorf("message channel closed")
		}
		return msg, nil

	case err, ok := <-c.errorChan:
		if !ok {
			return nil, fmt.Errorf("error channel closed")
		}
		return nil, err

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *ReaderImpl) parseMessage(data []byte) (*ResMessage, error) {
	var msg ResMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("error deserializing message: %w", err)
	}
	return &msg, nil
}
