package create_order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"order/internal/infrastructure/logger"
	"sync"

	otelkafkakonsumer "github.com/Trendyol/otel-kafka-konsumer"
	"github.com/segmentio/kafka-go"
)

type Reader interface {
	Start(ctx context.Context) error
	Read(ctx context.Context) (*ResEnvelope, error)
	Stop() error
}

type ReaderImpl struct {
	reader     *otelkafkakonsumer.Reader
	resultChan chan *ResEnvelope
	errorChan  chan error

	cancelCtx  context.Context
	cancelFunc context.CancelFunc

	wg      sync.WaitGroup
	mu      sync.Mutex
	started bool

	logger logger.Logger
}

func NewReader(reader *otelkafkakonsumer.Reader, logger logger.Logger) *ReaderImpl {
	return &ReaderImpl{
		reader: reader,
		logger: logger,
	}
}

func (r *ReaderImpl) log(level logger.Level, action, message string, extraFields map[string]any) {
	fields := map[string]any{
		"component": "create_order_saga_reader",
		"action":    action,
		"topic":     r.reader.R.Config().Topic,
	}
	for k, v := range extraFields {
		fields[k] = v
	}

	r.logger.Log(level, message, fields)
}

func (r *ReaderImpl) sendError(err error, action string) {
	r.log(logger.Error, action, err.Error(), nil)

	select {
	case r.errorChan <- fmt.Errorf("error reading message: %w", err):
	default:
		r.log(logger.Error, "channel_full", "Error channel full", map[string]any{"error": err.Error()})
	}
}

func (r *ReaderImpl) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.started {
		return errors.New("reader is already started")
	}

	r.resultChan = make(chan *ResEnvelope, 1)
	r.errorChan = make(chan error, 1)

	r.cancelCtx, r.cancelFunc = context.WithCancel(ctx)
	r.started = true

	r.log(logger.Info, "start", "Starting create order saga reader", nil)
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.readResults(r.cancelCtx)
	}()
	return nil
}

func (r *ReaderImpl) readResults(ctx context.Context) {
	defer r.log(logger.Info, "goroutine_completed", "Create order saga reader goroutine completed", nil)

	for {
		select {
		case <-ctx.Done():
			r.log(logger.Info, "stop", "Create order saga reader stopping", map[string]any{
				"reason": ctx.Err().Error(),
			})
			return

		default:
			// Read message
			msg, err := r.reader.ReadMessage(ctx)
			if ctx.Err() != nil {
				continue
			}
			if err != nil {
				r.sendError(err, "read_message")
				continue
			}

			// Parse the result message
			res, err := r.parseResultEnvelope(ctx, msg)
			if err != nil {
				r.log(logger.Error, "parse_error", "Failed to parse result message", map[string]any{
					"error":    err.Error(),
					"raw_data": msg.Value,
				})
				r.sendError(err, "parse_error")
				continue
			}
			r.log(logger.Info, "result_parsed", "Result parsed successfully", map[string]any{
				"result":    res.Msg,
				"partition": msg.Partition,
				"offset":    msg.Offset,
			})

			// Send the result to the result channel
			select {
			case r.resultChan <- res:
				r.log(logger.Info, "result_queued", "Result queued for processing", map[string]any{
					"result_id": res.Msg.ID,
				})
			case <-ctx.Done():
				continue
			}
		}
	}
}

func (r *ReaderImpl) Read(ctx context.Context) (*ResEnvelope, error) {
	select {
	case msg, ok := <-r.resultChan:
		if !ok {
			return nil, fmt.Errorf("message channel closed")
		}
		return msg, nil
	case err, ok := <-r.errorChan:
		if !ok {
			return nil, fmt.Errorf("error channel closed")
		}
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *ReaderImpl) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.started {
		return errors.New("create order saga reader is already stopped or was not started")
	}

	r.log(logger.Info, "stop_request", "Stopping create order saga reader", nil)
	r.cancelFunc()
	r.wg.Wait()
	close(r.resultChan)
	close(r.errorChan)
	r.started = false

	r.log(logger.Info, "stopped", "Create order saga reader stopped", nil)
	return nil
}

func (r *ReaderImpl) parseResultEnvelope(ctx context.Context, msg *kafka.Message) (*ResEnvelope, error) {
	cmdMsg, err := parseResultMessage(msg.Value)
	if err != nil {
		return nil, err
	}

	ctx = r.reader.TraceConfig.Propagator.Extract(ctx, otelkafkakonsumer.NewMessageCarrier(msg))

	return &ResEnvelope{
		Ctx: ctx,
		Msg: cmdMsg,
	}, nil
}

func parseResultMessage(data []byte) (*ResMessage, error) {
	var msg ResMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("error deserializing message: %w", err)
	}
	return &msg, nil
}
