package outbox

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
	"warehouse/internal/domain/outbox"
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
}

func NewProcessor(repository outbox.Repository, publisher outbox.Publisher) *Processor {
	return &Processor{
		repository: repository,
		publisher:  publisher,
		pollDelay:  DefaultPollDelay,
	}
}

func (p *Processor) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.started {
		return errors.New("outbox processor is already running; no need to start again")
	}

	p.cancelCtx, p.cancelFunc = context.WithCancel(ctx)
	p.started = true

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

	return nil
}

func (p *Processor) processOutbox(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Outbox processor stopping due to context cancellation.")
			return
		default:
			if err := p.processBatch(ctx); err != nil {
				log.Printf("Error processing outbox batch: %v", err)
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
	return nil
}
