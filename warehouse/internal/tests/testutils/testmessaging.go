package testutils

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
	kafkaContainer "github.com/testcontainers/testcontainers-go/modules/kafka"
	"net"
	"strconv"
)

const (
	TestTopic = "test"
)

type TestMessaging struct {
	Container testcontainers.Container
	Writer    *kafka.Writer
	Reader    *kafka.Reader
}

func SetUpContainer(ctx context.Context) (testcontainers.Container, error) {
	return kafkaContainer.Run(ctx,
		"confluentinc/confluent-local:7.5.0",
		kafkaContainer.WithClusterID("test-cluster"),
	)
}

func CreateUrl(ctx context.Context, container testcontainers.Container) (string, error) {
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

func createTopic(ctx context.Context, url string) error {
	conn, err := kafka.DialContext(ctx, "tcp", url)
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

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             TestTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	return controllerConn.CreateTopics(topicConfigs...)
}

func CreateWriter(url string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(url),
		Topic:    TestTopic,
		Balancer: &kafka.LeastBytes{},
	}
}

func CreateReader(url string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{url},
		GroupID: "consumer-group-id",
		Topic:   TestTopic,
	})
}

func NewTestMessaging(ctx context.Context) (*TestMessaging, error) {
	container, err := SetUpContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup container: %w", err)
	}
	defer func() {
		if err != nil {
			_ = container.Terminate(ctx)
		}
	}()

	url, err := CreateUrl(ctx, container)
	if err != nil {
		return nil, fmt.Errorf("failed to create url: %w", err)
	}

	if err = createTopic(ctx, url); err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	writer := CreateWriter(url)
	reader := CreateReader(url)

	return &TestMessaging{
		Container: container,
		Writer:    writer,
		Reader:    reader,
	}, nil
}
