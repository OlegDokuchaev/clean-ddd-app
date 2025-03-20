package create_order

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type Reader interface {
	Start(ctx context.Context)
	Read(ctx context.Context) (*ResMessage, error)
	Stop()
}

type ReaderImpl struct {
	reader      *kafka.Reader
	messageChan chan *ResMessage
	errorChan   chan error

	cancelCtx  context.Context
	cancelFunc context.CancelFunc

	wg      sync.WaitGroup
	mu      sync.Mutex
	started bool
}

func NewReader(reader *kafka.Reader) *ReaderImpl {
	return &ReaderImpl{
		reader: reader,
	}
}

func (r *ReaderImpl) Start(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.started {
		log.Println("Reader is already started, no need to start again.")
		return
	}

	r.messageChan = make(chan *ResMessage, 1)
	r.errorChan = make(chan error, 1)

	r.cancelCtx, r.cancelFunc = context.WithCancel(ctx)

	r.started = true
	log.Println("Starting reader...")

	r.wg.Add(1)

	go func() {
		defer r.wg.Done()
		r.consumeMessages(r.cancelCtx, r.messageChan, r.errorChan)
	}()
}

func (r *ReaderImpl) consumeMessages(ctx context.Context, msgCh chan<- *ResMessage, errCh chan<- error) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Context canceled, stopping consumer")
			return
		default:
			msg, err := r.reader.ReadMessage(ctx)
			if ctx.Err() != nil {
				log.Printf("Context is done, stopping consumer")
				return
			}

			if err != nil {
				select {
				case errCh <- fmt.Errorf("error reading message: %w", err):
				default:
					log.Printf("Error channel is full, dropping error: %v", err)
				}
				continue
			}

			resMsg, err := r.parseMessage(msg.Value)
			if err != nil {
				select {
				case errCh <- fmt.Errorf("error parsing message: %w", err):
				default:
					log.Printf("Error channel is full, dropping error: %v", err)
				}
				continue
			}

			select {
			case msgCh <- resMsg:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (r *ReaderImpl) Read(ctx context.Context) (*ResMessage, error) {
	select {
	case msg, ok := <-r.messageChan:
		if !ok {
			return nil, fmt.Errorf("message channel closed")
		}
		return msg, nil
	case err, ok := <-r.errorChan:
		if !ok {
			return nil, fmt.Errorf("error channel closed")
		}
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *ReaderImpl) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.started {
		log.Printf("Reader is already stopped or was not started.")
		return
	}

	log.Printf("Stopping reader...")

	if r.cancelFunc != nil {
		r.cancelFunc()
	}

	r.wg.Wait()

	close(r.messageChan)
	close(r.errorChan)

	r.started = false
	log.Printf("Reader has been stopped.")
}

func (r *ReaderImpl) parseMessage(data []byte) (*ResMessage, error) {
	var msg ResMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("error deserializing message: %w", err)
	}
	return &msg, nil
}
