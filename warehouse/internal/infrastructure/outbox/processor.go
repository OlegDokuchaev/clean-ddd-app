package outbox

import (
	"context"
	"errors"
	"sync"
	"time"
	"warehouse/internal/domain/outbox"
	"warehouse/internal/infrastructure/logger"
)

const DefaultPollDelay = 1 * time.Second

type Processor struct {
	repository outbox.Repository
	publisher  outbox.Publisher
	pollDelay  time.Duration

	cancelCtx  context.Context
	cancelFunc context.CancelFunc

	wg      sync.WaitGroup
	mu      sync.Mutex
	started bool

	logger logger.Logger
}

func NewProcessor(repository outbox.Repository, publisher outbox.Publisher, logger logger.Logger) *Processor {
	return &Processor{
		repository: repository,
		publisher:  publisher,
		pollDelay:  DefaultPollDelay,
		logger:     logger,
	}
}

func (p *Processor) log(level logger.Level, action, message string, extra map[string]any) {
	fields := map[string]any{
		"component": "outbox_processor",
		"action":    action,
	}
	for k, v := range extra {
		fields[k] = v
	}

	p.logger.Log(level, message, fields)
}

func (p *Processor) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.started {
		return errors.New("outbox processor is already running; no need to start again")
	}

	p.cancelCtx, p.cancelFunc = context.WithCancel(ctx)
	p.started = true

	p.log(logger.Info, "start", "Outbox processor started", nil)

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.processOutbox(p.cancelCtx)
	}()

	return nil
}

func (p *Processor) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.started {
		return errors.New("outbox processor is not running or already stopped")
	}

	p.cancelFunc()
	p.wg.Wait()
	p.started = false

	p.log(logger.Info, "stopped", "Outbox processor stopped", nil)

	return nil
}

func (p *Processor) processOutbox(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			p.log(logger.Info, "stopping", "Outbox processor stopping due to context cancellation", nil)
			return
		default:
			if err := p.processBatch(ctx); err != nil {
				p.log(logger.Error, "batch_error", "Error processing outbox batch", map[string]any{
					"error": err.Error(),
				})
			}
			time.Sleep(p.pollDelay)
		}
	}
}

func (p *Processor) processBatch(ctx context.Context) error {
	messages, err := p.repository.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, message := range messages {
		if err := p.processMessage(ctx, message); err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) processMessage(ctx context.Context, message *outbox.Message) error {
	if err := p.publisher.Publish(ctx, message); err != nil {
		return err
	}
	if err := p.repository.Delete(ctx, message); err != nil {
		return err
	}

	p.log(logger.Info, "processed", "Message processed", map[string]any{"message": message})
	return nil
}
