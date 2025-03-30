package commands

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
	Read(ctx context.Context) (*CmdMessage, error)
	Stop() error
}

type ReaderImpl struct {
	reader      *kafka.Reader
	commandChan chan *CmdMessage
	errorChan   chan error

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
		r.logger.Println("Command reader is already started, no need to start again.")
		return errors.New("command reader is already started")
	}

	r.commandChan = make(chan *CmdMessage, 1)
	r.errorChan = make(chan error, 1)

	r.cancelCtx, r.cancelFunc = context.WithCancel(ctx)

	r.started = true
	r.logger.Println("Starting command reader...")

	r.wg.Add(1)

	go func() {
		defer r.wg.Done()
		r.readCommands(r.cancelCtx)
	}()

	return nil
}

func (r *ReaderImpl) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.started {
		return errors.New("command reader is already stopped or was not started")
	}

	r.logger.Printf("Stopping command reader...")

	if r.cancelFunc != nil {
		r.cancelFunc()
	}

	r.wg.Wait()

	close(r.commandChan)
	close(r.errorChan)

	r.started = false
	r.logger.Printf("Command reader has been stopped.")

	return nil
}

func (r *ReaderImpl) readCommands(ctx context.Context) {
	defer func() {
		r.logger.Printf("Command reader goroutine completed")
	}()

	for {
		select {
		case <-ctx.Done():
			r.logger.Printf("Context canceled, stopping command reader")
			return
		default:
			msg, err := r.reader.ReadMessage(ctx)

			if ctx.Err() != nil {
				r.logger.Printf("Context is done, stopping command reader")
				return
			}

			if err != nil {
				select {
				case r.errorChan <- fmt.Errorf("error reading message: %w", err):
				default:
					r.logger.Printf("Error channel full, dropping error: %v", err)
				}
				continue
			}

			cmdMsg, err := parseCommandMessage(msg.Value)
			if err != nil {
				select {
				case r.errorChan <- fmt.Errorf("error parsing message: %w", err):
				default:
					r.logger.Printf("Error channel full, dropping error: %v", err)
				}
				continue
			}

			select {
			case r.commandChan <- cmdMsg:
				r.logger.Printf("Command received: %s, type: %s", cmdMsg.ID, cmdMsg.Name)
			case <-ctx.Done():
				return
			}
		}
	}
}

func (r *ReaderImpl) Read(ctx context.Context) (*CmdMessage, error) {
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

var _ Reader = (*ReaderImpl)(nil)

func parseCommandMessage(data []byte) (*CmdMessage, error) {
	var msg CmdMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("error deserializing message: %w", err)
	}
	return &msg, nil
}
