package messaging

import (
	otelkafkakonsumer "github.com/Trendyol/otel-kafka-konsumer"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func NewWarehouseCmdReader(config *Config, tp *sdktrace.TracerProvider) (*otelkafkakonsumer.Reader, error) {
	return otelkafkakonsumer.NewReader(
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{config.Address},
			GroupID: config.WarehouseCmdConsumerGroupID,
			Topic:   config.WarehouseCmdTopic,
		}),
		otelkafkakonsumer.WithTracerProvider(tp),
		otelkafkakonsumer.WithPropagator(propagation.TraceContext{}),
		otelkafkakonsumer.WithAttributes(
			[]attribute.KeyValue{
				semconv.MessagingDestinationKindTopic,
				semconv.MessagingKafkaClientIDKey.String(config.WarehouseCmdTopic),
			},
		),
	)
}

func NewProductEventReader(config *Config, tp *sdktrace.TracerProvider) (*otelkafkakonsumer.Reader, error) {
	return otelkafkakonsumer.NewReader(
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{config.Address},
			GroupID: config.ProductEventConsumerGroupID,
			Topic:   config.ProductEventTopic,
		}),
		otelkafkakonsumer.WithTracerProvider(tp),
		otelkafkakonsumer.WithPropagator(propagation.TraceContext{}),
		otelkafkakonsumer.WithAttributes(
			[]attribute.KeyValue{
				semconv.MessagingDestinationKindTopic,
				semconv.MessagingKafkaClientIDKey.String(config.ProductEventTopic),
			},
		),
	)
}
