package commands

import (
	"context"
	"courier/internal/infrastructure/logger"
	"errors"
	"sync"
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

func (p *Processor) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.started {
		return errors.New("processor is already running, no need to start again")
	}

	p.cancelCtx, p.cancelFunc = context.WithCancel(ctx)
	p.started = true

	p.wg.Add(1)
	go p.processCommands(p.cancelCtx)

	return nil
}

func (p *Processor) processCommands(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			p.logger.Printf("%s processor stopping: context done", ctx)
			return
		default:
			cmd, err := p.reader.Read(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				p.logger.Printf("Error reading command: %v", err)
				continue
			}

			res, err := p.handler.Handle(ctx, cmd)
			if err != nil {
				p.logger.Printf("Error handling command %s: %v", cmd.ID, err)
				continue
			}

			if res != nil {
				if err := p.writer.Write(ctx, res); err != nil {
					p.logger.Printf("Error sending response: %v", err)
				}
			}

			p.logger.Printf("Command %s successfully processed", cmd.ID)
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
