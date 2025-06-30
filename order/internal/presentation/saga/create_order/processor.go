package create_order

import (
	"context"
	"order/internal/infrastructure/logger"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Processor struct {
	handler         Handler
	warehouseReader Reader
	courierReader   Reader

	cancelCtx  context.Context
	cancelFunc context.CancelFunc

	wg      sync.WaitGroup
	mu      sync.Mutex
	started bool

	logger logger.Logger
}

func NewProcessor(handler Handler, warehouseReader Reader, courierReader Reader, logger logger.Logger) *Processor {
	return &Processor{
		handler:         handler,
		warehouseReader: warehouseReader,
		courierReader:   courierReader,
		logger:          logger,
	}
}

func (p *Processor) log(level logger.Level, action, message string, extraFields map[string]any) {
	fields := map[string]any{
		"component": "create_order_saga_processor",
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
		return errors.New("processor is already running, no need to start again.")
	}

	p.cancelCtx, p.cancelFunc = context.WithCancel(ctx)
	p.started = true

	p.log(logger.Info, "start", "Starting create order saga processor", nil)
	p.wg.Add(2)
	go p.processMessages(p.cancelCtx, "warehouse", p.warehouseReader)
	go p.processMessages(p.cancelCtx, "courier", p.courierReader)
	return nil
}

func (p *Processor) processMessages(ctx context.Context, source string, receiver Reader) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			p.log(logger.Info, "stop", "Create order saga processor stopping", map[string]any{
				"reason": ctx.Err().Error(),
			})
			return

		default:
			// Read the result
			msg, err := receiver.Read(ctx)
			if ctx.Err() != nil {
				continue
			}
			if err != nil {
				p.log(logger.Error, "read", "Error reading command", map[string]any{
					"source": source,
					"error":  err.Error(),
				})
				continue
			}

			// Handle the command
			startTime := time.Now()
			err = p.handler.Handle(ctx, msg)
			duration := time.Since(startTime)

			if err != nil {
				p.log(logger.Error, "process_error", "Result processing failed", map[string]any{
					"result_id":   msg.ID,
					"source":      source,
					"error":       err.Error(),
					"duration_ms": duration.Milliseconds(),
				})
				continue
			}

			p.log(logger.Info, "process_success", "Result processed successfully", map[string]any{
				"result_id":   msg.ID,
				"source":      source,
				"duration_ms": duration.Milliseconds(),
			})
		}
	}
}

func (p *Processor) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.started {
		return errors.New("processor is not running or already stopped")
	}

	p.log(logger.Info, "stop_request", "Stopping create order saga processor", nil)
	p.cancelFunc()
	p.wg.Wait()
	p.started = false

	p.log(logger.Info, "stopped", "Create order saga processor stopped", nil)
	return nil
}
