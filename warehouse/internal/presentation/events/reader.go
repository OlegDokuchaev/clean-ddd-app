package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"warehouse/internal/infrastructure/logger"

	"github.com/segmentio/kafka-go"
)

type Reader interface {
	Start(ctx context.Context) error
	Read(ctx context.Context) (*Event, error)
	Stop() error
}

type ReaderImpl struct {
	reader    *kafka.Reader
	eventChan chan *Event
	errorChan chan error

	cancelCtx  context.Context
	cancelFunc context.CancelFunc

	wg      sync.WaitGroup
	mu      sync.Mutex
	started bool

	logger logger.Logger
}

func NewReader(reader *kafka.Reader, logger logger.Logger) *ReaderImpl {
	return &ReaderImpl{
		reader: reader,
		logger: logger,
	}
}

func (r *ReaderImpl) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.started {
		return errors.New("event reader is already started")
	}

	r.eventChan = make(chan *Event, 1)
	r.errorChan = make(chan error, 1)

	r.cancelCtx, r.cancelFunc = context.WithCancel(ctx)

	r.started = true
	r.logger.Println("Starting event reader...")

	r.wg.Add(1)

	go func() {
		defer r.wg.Done()
		r.readEvents(r.cancelCtx)
	}()

	return nil
}

func (r *ReaderImpl) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.started {
		return errors.New("event reader is already stopped or was not started")
	}

	r.logger.Printf("Stopping event reader...")

	if r.cancelFunc != nil {
		r.cancelFunc()
	}

	r.wg.Wait()

	close(r.eventChan)
	close(r.errorChan)

	r.started = false
	r.logger.Printf("Event reader has been stopped.")

	return nil
}

func (r *ReaderImpl) readEvents(ctx context.Context) {
	defer func() {
		r.logger.Printf("Event reader goroutine completed")
	}()

	for {
		select {
		case <-ctx.Done():
			r.logger.Printf("Context canceled, stopping event reader")
			return
		default:
			msg, err := r.reader.ReadMessage(ctx)

			if ctx.Err() != nil {
				continue
			}

			if err != nil {
				select {
				case r.errorChan <- fmt.Errorf("error reading message: %w", err):
				case <-ctx.Done():
				}
				continue
			}

			event, err := parseEvent(msg.Value)
			if err != nil {
				select {
				case r.errorChan <- fmt.Errorf("error parsing message: %w", err):
				case <-ctx.Done():
				}
				continue
			}

			select {
			case r.eventChan <- event:
				r.logger.Printf("Event received: %s, type: %s", event.ID, event.Name)
			case <-ctx.Done():
			}
		}
	}
}

func (r *ReaderImpl) Read(ctx context.Context) (*Event, error) {
	select {
	case cmd, ok := <-r.eventChan:
		if !ok {
			return nil, fmt.Errorf("event channel closed")
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

var _ Reader = (*ReaderImpl)(nil)

func parseEvent(data []byte) (*Event, error) {
	var msg Event
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("error deserializing message: %w", err)
	}
	return &msg, nil
}
