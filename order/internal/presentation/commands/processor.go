package commands

import (
	"context"
	"log"
	"sync"
)

type Processor struct {
	handler      Handler
	reader       Reader
	writer       Writer
	wg           sync.WaitGroup
	shutdownChan chan struct{}
}

func NewProcessor(handler Handler, reader Reader, writer Writer) *Processor {
	return &Processor{
		handler:      handler,
		reader:       reader,
		writer:       writer,
		shutdownChan: make(chan struct{}),
	}
}

func (p *Processor) Start(ctx context.Context) {
	p.wg.Add(1)
	go p.processCommands(ctx)
}

func (p *Processor) processCommands(ctx context.Context) {
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

	log.Println("Starting command processor")

	for {
		select {
		case <-processorCtx.Done():
			log.Println("Stopping command processor: context done")
			return
		default:
			cmd, err := p.reader.Read(processorCtx)
			if err != nil {
				if processorCtx.Err() != nil {
					return
				}
				log.Printf("Error reading command: %v", err)
				continue
			}

			res, err := p.handler.Handle(processorCtx, cmd)
			if err != nil {
				log.Printf("Error handling command %s: %v", cmd.ID, err)
				continue
			}

			if res != nil {
				if err := p.writer.Write(processorCtx, res); err != nil {
					log.Printf("Error sending response: %v", err)
				}
			}

			log.Printf("Command %s successfully processed", cmd.ID)
		}
	}
}

func (p *Processor) Stop() {
	close(p.shutdownChan)
	p.wg.Wait()
}

func (p *Processor) Close() error {
	var lastErr error

	p.Stop()

	if err := p.reader.Close(); err != nil {
		log.Printf("Error closing command reader: %v", err)
		lastErr = err
	}

	if err := p.writer.Close(); err != nil {
		log.Printf("Error closing response writer: %v", err)
		lastErr = err
	}

	return lastErr
}
