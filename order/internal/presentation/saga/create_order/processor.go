package create_order

import (
	"context"
	"order/internal/infrastructure/logger"
	"sync"

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

func (p *Processor) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.started {
		return errors.New("processor is already running, no need to start again.")
	}

	p.cancelCtx, p.cancelFunc = context.WithCancel(ctx)
	p.started = true

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
			p.logger.Printf("%s processor stopping: context done", source)
			return
		default:
			msg, err := receiver.Read(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				p.logger.Printf("Error receiving message from %s: %v", source, err)
				continue
			}

			if err := p.handler.Handle(ctx, msg); err != nil {
				p.logger.Printf("Error handling %s message: %v", source, err)
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

	p.cancelFunc()

	p.wg.Wait()
	p.started = false

	return nil
}
