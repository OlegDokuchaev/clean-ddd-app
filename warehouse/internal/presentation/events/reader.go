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

func (r *ReaderImpl) log(level logger.Level, action, message string, extraFields map[string]any) {
	fields := map[string]any{
		"component": "event_reader",
		"action":    action,
		"topic":     r.reader.Config().Topic,
	}
	for k, v := range extraFields {
		fields[k] = v
	}

	r.logger.Log(level, message, fields)
}

func (r *ReaderImpl) sendError(err error, action string) {
	r.log(logger.Error, action, err.Error(), nil)

	select {
	case r.errorChan <- fmt.Errorf("error reading message: %w", err):
	default:
		r.log(logger.Error, "channel_full", "Error channel full", map[string]any{"error": err.Error()})
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

	r.log(logger.Info, "start", "Starting event reader", nil)
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

	r.log(logger.Info, "stop_request", "Stopping event reader", nil)
	r.cancelFunc()
	r.wg.Wait()
	close(r.eventChan)
	close(r.errorChan)
	r.started = false

	r.log(logger.Info, "stopped", "Event reader stopped", nil)
	return nil
}

func (r *ReaderImpl) readEvents(ctx context.Context) {
	defer r.log(logger.Info, "goroutine_completed", "Event reader goroutine completed", nil)

	for {
		select {
		case <-ctx.Done():
			r.log(logger.Info, "stop", "Event reader stopping", map[string]any{"reason": ctx.Err().Error()})
			return

		default:
			// Read message
			msg, err := r.reader.ReadMessage(ctx)
			if ctx.Err() != nil {
				continue
			}
			if err != nil {
				r.sendError(err, "read_message")
				continue
			}

			// Parse the event
			event, err := parseEvent(msg.Value)
			if err != nil {
				r.log(logger.Error, "parse_error", "Failed to parse event", map[string]any{
					"error":    err.Error(),
					"raw_data": msg.Value,
				})
				r.sendError(err, "parse_error")
				continue
			}
			r.log(logger.Info, "event_parsed", "Event parsed successfully", map[string]any{
				"event":     event,
				"partition": msg.Partition,
				"offset":    msg.Offset,
			})

			// Send the command to the command channel
			select {
			case r.eventChan <- event:
				r.log(logger.Info, "event_queued", "Event queued for processing", map[string]any{
					"event_id": event.ID,
				})
			case <-ctx.Done():
				continue
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
