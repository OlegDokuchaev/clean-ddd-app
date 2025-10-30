package messaging

import (
	otelkafkakonsumer "github.com/Trendyol/otel-kafka-konsumer"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func NewCourierCmdResWriter(config *Config, tp *sdktrace.TracerProvider) (*otelkafkakonsumer.Writer, error) {
	return otelkafkakonsumer.NewWriter(
		&kafka.Writer{
			Addr:  kafka.TCP(config.Address),
			Topic: config.CourierCmdResTopic,
		},
		otelkafkakonsumer.WithTracerProvider(tp),
		otelkafkakonsumer.WithPropagator(propagation.TraceContext{}),
		otelkafkakonsumer.WithAttributes(
			[]attribute.KeyValue{
				semconv.MessagingDestinationKindTopic,
				semconv.MessagingKafkaClientIDKey.String("CourierCmdTopic"),
			},
		),
	)
}
