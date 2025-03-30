package events

import (
	"context"
	"errors"
	"sync"
	"warehouse/internal/infrastructure/logger"
)

type Processor struct {
	handler       Handler
	productReader Reader

	cancelCtx  context.Context
	cancelFunc context.CancelFunc

	wg      sync.WaitGroup
	mu      sync.Mutex
	started bool

	logger logger.Logger
}

func NewProcessor(handler Handler, productReader Reader, logger logger.Logger) *Processor {
	return &Processor{
		handler:       handler,
		productReader: productReader,
		logger:        logger,
	}
}

func (p *Processor) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.started {
		return errors.New("processor is already running, no need to start again")
	}

	p.cancelCtx, p.cancelFunc = context.WithCancel(ctx)
	p.started = true

	p.wg.Add(1)
	go p.processEvents(p.cancelCtx, "product", p.productReader)

	return nil
}

func (p *Processor) processEvents(ctx context.Context, source string, reader Reader) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			p.logger.Printf("%s processor stopping: context done", source)
			return
		default:
			event, err := reader.Read(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				p.logger.Printf("Error reading event from %s: %v", source, err)
				continue
			}

			err = p.handler.Handle(ctx, event)
			if err != nil {
				p.logger.Printf("Error handling event %s from %s: %v", event.ID, source, err)
				continue
			}

			p.logger.Printf("Event %s successfully processed from %s", event.ID, source)
		}
	}
}

func (p *Processor) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.started {
		return errors.New("processor is not running or already stopped")
	}

	p.cancelFunc()

	p.wg.Wait()
	p.started = false

	return nil
}
