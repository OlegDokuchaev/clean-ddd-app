package messaging

import (
	otelkafkakonsumer "github.com/Trendyol/otel-kafka-konsumer"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func NewOrderCommandReader(config *Config, tp *sdktrace.TracerProvider) (*otelkafkakonsumer.Reader, error) {
	return otelkafkakonsumer.NewReader(
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{config.Address},
			GroupID: config.OrderCmdConsumerGroupID,
			Topic:   config.OrderCmdTopic,
		}),
		otelkafkakonsumer.WithTracerProvider(tp),
		otelkafkakonsumer.WithPropagator(propagation.TraceContext{}),
		otelkafkakonsumer.WithAttributes(
			[]attribute.KeyValue{
				semconv.MessagingDestinationKindTopic,
				semconv.MessagingKafkaClientIDKey.String(config.OrderCmdTopic),
			},
		),
	)
}

func NewWarehouseCommandResultReader(config *Config, tp *sdktrace.TracerProvider) (*otelkafkakonsumer.Reader, error) {
	return otelkafkakonsumer.NewReader(
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{config.Address},
			GroupID: config.WarehouseCmdResConsumerGroupID,
			Topic:   config.WarehouseCmdResTopic,
		}),
		otelkafkakonsumer.WithTracerProvider(tp),
		otelkafkakonsumer.WithPropagator(propagation.TraceContext{}),
		otelkafkakonsumer.WithAttributes(
			[]attribute.KeyValue{
				semconv.MessagingDestinationKindTopic,
				semconv.MessagingKafkaClientIDKey.String(config.WarehouseCmdResTopic),
			},
		),
	)
}

func NewCourierCommandResultReader(config *Config, tp *sdktrace.TracerProvider) (*otelkafkakonsumer.Reader, error) {
	return otelkafkakonsumer.NewReader(
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{config.Address},
			GroupID: config.CourierCmdResConsumerGroupID,
			Topic:   config.CourierCmdResTopic,
		}),
		otelkafkakonsumer.WithTracerProvider(tp),
		otelkafkakonsumer.WithPropagator(propagation.TraceContext{}),
		otelkafkakonsumer.WithAttributes(
			[]attribute.KeyValue{
				semconv.MessagingDestinationKindTopic,
				semconv.MessagingKafkaClientIDKey.String(config.CourierCmdResTopic),
			},
		),
	)
}
