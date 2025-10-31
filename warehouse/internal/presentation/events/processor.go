package events

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"sync"
	"time"
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

func (p *Processor) log(level logger.Level, action, message string, extraFields map[string]any) {
	fields := map[string]any{
		"component": "event_processor",
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

	p.log(logger.Info, "start", "Starting event processor", nil)
	p.wg.Add(1)
	go p.processEvents(p.cancelCtx, p.productReader)
	return nil
}

func (p *Processor) processEvents(ctx context.Context, reader Reader) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			p.log(logger.Info, "stop", "EventMessage processor stopping", map[string]any{"reason": ctx.Err().Error()})
			return

		default:
			// Read the event
			event, err := reader.Read(ctx)
			if ctx.Err() != nil {
				continue
			}
			if err != nil {
				p.log(logger.Error, "read", "Error reading event", map[string]any{"error": err.Error()})
				continue
			}

			// Handle the event
			sCtx, span := startProcessSpan(event)
			startTime := time.Now()

			err = p.handler.Handle(sCtx, event.Msg)

			duration := time.Since(startTime)
			span.End()

			if err != nil {
				p.log(logger.Error, "process_error", "EventMessage processing failed", map[string]any{
					"event_id":    event.Msg.ID,
					"error":       err.Error(),
					"duration_ms": duration.Milliseconds(),
				})
				continue
			}

			p.log(logger.Info, "process_success", "EventMessage processed successfully", map[string]any{
				"event_id":    event.Msg.ID,
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

	p.log(logger.Info, "stop_request", "Stopping event processor", nil)
	p.cancelFunc()
	p.wg.Wait()
	p.started = false

	p.log(logger.Info, "stopped", "EventMessage processor stopped", nil)
	return nil
}

func startProcessSpan(event *EventEnvelope) (context.Context, trace.Span) {
	return otel.Tracer("warehouse-service.events").Start(
		event.Ctx,
		"kafka.process",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			semconv.MessagingSystemKey.String("kafka"),
			semconv.MessagingDestinationName(event.Topic),
			semconv.MessagingKafkaDestinationPartition(event.Partition),
			semconv.MessagingOperationKey.String("process"),
			attribute.String("event.id", event.Msg.ID.String()),
		),
	)
}
