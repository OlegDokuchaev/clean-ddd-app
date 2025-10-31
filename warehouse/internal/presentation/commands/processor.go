package commands

import (
	"context"
	"errors"
	"sync"
	"time"
	"warehouse/internal/infrastructure/logger"
)

type Processor struct {
	handler Handler
	reader  Reader
	writer  Writer

	cancelCtx  context.Context
	cancelFunc context.CancelFunc

	wg      sync.WaitGroup
	mu      sync.Mutex
	started bool

	logger logger.Logger
}

func NewProcessor(handler Handler, reader Reader, writer Writer, logger logger.Logger) *Processor {
	return &Processor{
		handler: handler,
		reader:  reader,
		writer:  writer,
		logger:  logger,
	}
}

func (p *Processor) log(level logger.Level, action, message string, extraFields map[string]any) {
	fields := map[string]any{
		"component": "command_processor",
		"action":    action,
	}
	for k, v := range extraFields {
		fields[k] = v
	}

	p.logger.Log(level, message, fields)
}

func (p *Processor) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.started {
		return errors.New("processor is already running, no need to start again")
	}

	p.cancelCtx, p.cancelFunc = context.WithCancel(ctx)
	p.started = true

	p.log(logger.Info, "start", "Starting command processor", nil)
	p.wg.Add(1)
	go p.processCommands(p.cancelCtx)
	return nil
}

func (p *Processor) processCommands(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			p.log(logger.Info, "stop", "Command processor stopping", map[string]any{"reason": ctx.Err().Error()})
			return

		default:
			// Read the command
			cmd, err := p.reader.Read(ctx)
			if ctx.Err() != nil {
				continue
			}
			if err != nil {
				p.log(logger.Error, "read", "Error reading command", map[string]any{"error": err.Error()})
				continue
			}

			// Handle the command
			startTime := time.Now()
			res, err := p.handler.Handle(cmd.Ctx, cmd.Msg)
			duration := time.Since(startTime)

			if err != nil {
				p.log(logger.Error, "process_error", "Command processing failed", map[string]any{
					"command_id":  cmd.Msg.ID,
					"error":       err.Error(),
					"duration_ms": duration.Milliseconds(),
				})
				continue
			}

			p.log(logger.Info, "process_success", "Command processed successfully", map[string]any{
				"command_id":   cmd.Msg.ID,
				"duration_ms":  duration.Milliseconds(),
				"has_response": res != nil,
			})

			// Write the response
			if res != nil {
				if err := p.writer.Write(cmd.Ctx, res); err != nil {
					p.log(logger.Error, "write_error", "Error sending response", map[string]any{
						"command_id":  cmd.Msg.ID,
						"response_id": res.ID,
						"error":       err.Error(),
					})
				}
			}
		}
	}
}

func (p *Processor) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.started {
		return errors.New("processor is not running or already stopped")
	}

	p.log(logger.Info, "stop_request", "Stopping command processor", nil)
	p.cancelFunc()
	p.wg.Wait()
	p.started = false

	p.log(logger.Info, "stopped", "Command processor stopped", nil)
	return nil
}
