package testutils

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
	kafkaContainer "github.com/testcontainers/testcontainers-go/modules/kafka"
)

type TestMessaging struct {
	Container testcontainers.Container
	url       string
}

func (m *TestMessaging) CreateTopics(ctx context.Context, topics ...string) error {
	conn, err := kafka.DialContext(ctx, "tcp", m.url)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.DialContext(ctx, "tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer func() {
		_ = controllerConn.Close()
	}()

	topicConfigs := make([]kafka.TopicConfig, 0, len(topics))
	for _, topic := range topics {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}
	return controllerConn.CreateTopics(topicConfigs...)
}

func (m *TestMessaging) CreateWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(m.url),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func (m *TestMessaging) CreateReader(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{m.url},
		Topic:   topic,
	})
}

func (m *TestMessaging) Close(ctx context.Context) error {
	return m.Container.Terminate(ctx)
}

func setUpContainer(ctx context.Context) (testcontainers.Container, error) {
	return kafkaContainer.Run(ctx,
		"confluentinc/confluent-local:7.5.0",
		kafkaContainer.WithClusterID("test-cluster"),
	)
}

func createUrl(ctx context.Context, container testcontainers.Container) (string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := container.MappedPort(ctx, "9093/tcp")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", host, port.Port()), nil
}

func NewTestMessaging(ctx context.Context) (*TestMessaging, error) {
	container, err := setUpContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup container: %w", err)
	}
	defer func() {
		if err != nil {
			_ = container.Terminate(ctx)
		}
	}()

	url, err := createUrl(ctx, container)
	if err != nil {
		return nil, fmt.Errorf("failed to create url: %w", err)
	}

	return &TestMessaging{
		Container: container,
		url:       url,
	}, nil
}
