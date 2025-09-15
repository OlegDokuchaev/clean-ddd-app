package testutils

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"order/internal/infrastructure/messaging"

	"github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
	kafkaContainer "github.com/testcontainers/testcontainers-go/modules/kafka"
)

type TestMessaging struct {
	Cfg       *messaging.Config
	container testcontainers.Container
}

const (
	EnvMsgMode               = "E2E_MSG_MODE"
	TestWarehouseTopic       = "warehouse-topic"
	TestWarehouseResTopic    = "warehouse-topic-res"
	TestWarehouseResGroupID  = "warehouse-res-consumer"
	TestOrderTopic           = "order-topic"
	TestOrderResTopic        = "order-topic-res"
	TestOrderConsumerGroupID = "order-consumer"
	TestCourierTopic         = "courier-topic"
	TestCourierResTopic      = "courier-topic-res"
	TestCourierResGroupID    = "courier-res-consumer"
)

func (m *TestMessaging) CreateWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(m.Cfg.Address),
		Topic: topic,
	}
}

func (m *TestMessaging) CreateReader(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{m.Cfg.Address},
		Topic:   topic,
	})
}

func (m *TestMessaging) Close(ctx context.Context) error {
	if m.container == nil {
		return nil
	}
	return m.container.Terminate(ctx)
}

func (m *TestMessaging) Clear(ctx context.Context) error {
	topics := []string{
		m.Cfg.CourierCmdTopic,
		m.Cfg.CourierCmdResTopic,
		m.Cfg.WarehouseCmdTopic,
		m.Cfg.WarehouseCmdResTopic,
		m.Cfg.OrderCmdTopic,
		m.Cfg.OrderCmdResTopic,
	}

	conn, err := kafka.DialContext(ctx, "tcp", m.Cfg.Address)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	admin, err := kafka.DialContext(ctx, "tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer func() { _ = admin.Close() }()

	_ = admin.DeleteTopics(topics...)

	topicConfigs := make([]kafka.TopicConfig, 0, len(topics))
	for _, t := range topics {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             t,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}
	return admin.CreateTopics(topicConfigs...)
}

func setupKafkaContainer(ctx context.Context) (testcontainers.Container, error) {
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

func createTopics(ctx context.Context, url string, topics ...string) error {
	conn, err := kafka.DialContext(ctx, "tcp", url)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.DialContext(ctx, "tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer func() { _ = controllerConn.Close() }()

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

func NewTestMessaging(ctx context.Context, cfg *messaging.Config) (*TestMessaging, error) {
	switch strings.ToLower(os.Getenv(EnvMsgMode)) {
	case "real":
		if cfg == nil {
			return nil, fmt.Errorf("messaging config Address must be set for real mode")
		}
		return &TestMessaging{Cfg: cfg}, nil
	default:
		container, err := setupKafkaContainer(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to setup container: %w", err)
		}

		url, err := createUrl(ctx, container)
		if err != nil {
			if err := container.Terminate(ctx); err != nil {
				return nil, fmt.Errorf("failed to terminate container: %w", err)
			}
			return nil, fmt.Errorf("failed to create url: %w", err)
		}

		if err = createTopics(ctx, url, TestWarehouseTopic, TestWarehouseResTopic, TestOrderTopic, TestOrderResTopic,
			TestCourierTopic, TestCourierResTopic); err != nil {
			return nil, fmt.Errorf("failed to create topics: %w", err)
		}

		cfg = &messaging.Config{
			Address: url,

			OrderCmdTopic:           TestOrderTopic,
			OrderCmdResTopic:        TestOrderResTopic,
			OrderCmdConsumerGroupID: TestOrderConsumerGroupID,

			WarehouseCmdTopic:              TestWarehouseTopic,
			WarehouseCmdResTopic:           TestWarehouseResTopic,
			WarehouseCmdResConsumerGroupID: TestWarehouseResGroupID,

			CourierCmdTopic:              TestCourierTopic,
			CourierCmdResTopic:           TestCourierResTopic,
			CourierCmdResConsumerGroupID: TestCourierResGroupID,
		}

		return &TestMessaging{Cfg: cfg, container: container}, nil
	}
}
