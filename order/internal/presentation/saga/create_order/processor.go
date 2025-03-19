package create_order

import (
	"context"
	"log"
	"sync"
)

type Processor struct {
	handler         Handler
	warehouseReader Reader
	courierReader   Reader
	wg              sync.WaitGroup
	shutdownChan    chan struct{}
}

func NewProcessor(handler Handler, warehouseReader Reader, courierReader Reader) *Processor {
	return &Processor{
		handler:         handler,
		warehouseReader: warehouseReader,
		courierReader:   courierReader,
		shutdownChan:    make(chan struct{}),
	}
}

func (p *Processor) Start(ctx context.Context) {
	p.wg.Add(2)
	go p.processMessages(ctx, "warehouse", p.warehouseReader)
	go p.processMessages(ctx, "courier", p.courierReader)
}

func (p *Processor) processMessages(ctx context.Context, source string, receiver Reader) {
	defer p.wg.Done()

	processorCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		select {
		case <-p.shutdownChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	for {
		select {
		case <-processorCtx.Done():
			log.Printf("%s processor stopping: context done", source)
			return
		default:
			msg, err := receiver.Read(processorCtx)
			if err != nil {
				if processorCtx.Err() != nil {
					return
				}
				log.Printf("Error receiving message from %s: %v", source, err)
				continue
			}

			if err := p.handler.Handle(processorCtx, msg); err != nil {
				log.Printf("Error handling %s message: %v", source, err)
			}
		}
	}
}

func (p *Processor) Stop() {
	close(p.shutdownChan)
	p.wg.Wait()
	log.Println("Processor gracefully stopped")
}
