package commands

import (
	"context"
	"courier/internal/infrastructure/logger"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	otelkafkakonsumer "github.com/Trendyol/otel-kafka-konsumer"
	"github.com/segmentio/kafka-go"
)

type Reader interface {
	Start(ctx context.Context) error
	Read(ctx context.Context) (*CmdEnvelope, error)
	Stop() error
}

type ReaderImpl struct {
	reader      *otelkafkakonsumer.Reader
	commandChan chan *CmdEnvelope
	errorChan   chan error

	cancelCtx  context.Context
	cancelFunc context.CancelFunc

	wg      sync.WaitGroup
	mu      sync.Mutex
	started bool

	logger logger.Logger
}

func NewReader(reader *otelkafkakonsumer.Reader, logger logger.Logger) *ReaderImpl {
	return &ReaderImpl{
		reader: reader,
		logger: logger,
	}
}

func (r *ReaderImpl) log(level logger.Level, action, message string, extraFields map[string]any) {
	fields := map[string]any{
		"component": "command_reader",
		"action":    action,
		"topic":     r.reader.R.Config().Topic,
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
		return errors.New("command reader is already started")
	}

	r.commandChan = make(chan *CmdEnvelope, 1)
	r.errorChan = make(chan error, 1)

	r.cancelCtx, r.cancelFunc = context.WithCancel(ctx)
	r.started = true

	r.log(logger.Info, "start", "Starting command reader", nil)
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

	r.log(logger.Info, "stop_request", "Stopping command reader", nil)
	r.cancelFunc()
	r.wg.Wait()
	close(r.commandChan)
	close(r.errorChan)
	r.started = false

	r.log(logger.Info, "stopped", "Command reader stopped", nil)
	return nil
}

func (r *ReaderImpl) readCommands(ctx context.Context) {
	defer r.log(logger.Info, "goroutine_completed", "Command reader goroutine completed", nil)

	for {
		select {
		case <-ctx.Done():
			r.log(logger.Info, "stop", "Command reader stopping", map[string]any{"reason": ctx.Err().Error()})
			return

		default:
			r.readCommand(ctx)
		}
	}
}

func (r *ReaderImpl) readCommand(ctx context.Context) {
	// Read message
	msg, err := r.reader.ReadMessage(ctx)
	if ctx.Err() != nil {
		return
	}
	if err != nil {
		r.sendError(err, "read_message")
		return
	}

	// Parse the command message
	cmdEnv, err := r.parseCommandEnvelope(ctx, msg)
	if err != nil {
		r.log(logger.Error, "parse_error", "Failed to parse command message", map[string]any{
			"error":    err.Error(),
			"raw_data": msg.Value,
		})
		r.sendError(err, "parse_error")
		return
	}
	r.log(logger.Info, "command_parsed", "Command parsed successfully", map[string]any{
		"command":   cmdEnv.Cmd,
		"partition": msg.Partition,
		"offset":    msg.Offset,
	})

	// Send the command to the command channel
	select {
	case r.commandChan <- cmdEnv:
		r.log(logger.Info, "command_queued", "Command queued for processing", map[string]any{
			"command_id": cmdEnv.Cmd.ID,
		})
	case <-ctx.Done():
	}
}

func (r *ReaderImpl) Read(ctx context.Context) (*CmdEnvelope, error) {
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

func (r *ReaderImpl) parseCommandEnvelope(ctx context.Context, msg *kafka.Message) (*CmdEnvelope, error) {
	cmdMsg, err := parseCommandMessage(msg.Value)
	if err != nil {
		return nil, err
	}

	ctx = r.reader.TraceConfig.Propagator.Extract(ctx, otelkafkakonsumer.NewMessageCarrier(msg))

	return &CmdEnvelope{
		Ctx: ctx,
		Cmd: cmdMsg,
	}, nil
}

var _ Reader = (*ReaderImpl)(nil)

func parseCommandMessage(data []byte) (*CmdMessage, error) {
	var msg CmdMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("error deserializing message: %w", err)
	}
	return &msg, nil
}
