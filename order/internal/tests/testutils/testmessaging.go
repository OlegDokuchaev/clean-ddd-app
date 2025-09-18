package testutils

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"order/internal/infrastructure/messaging"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
	kafkaContainer "github.com/testcontainers/testcontainers-go/modules/kafka"
)

type TestMessaging struct {
	Cfg       *messaging.Config
	container testcontainers.Container
}

const (
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

func (m *TestMessaging) PurgeTopic(ctx context.Context, topic string) error {
	admin, err := confluentKafka.NewAdminClient(&confluentKafka.ConfigMap{
		"bootstrap.servers": m.Cfg.Address,
	})
	if err != nil {
		return err
	}
	defer admin.Close()

	md, err := admin.GetMetadata(&topic, false, 10_000)
	if err != nil {
		return fmt.Errorf("get metadata: %w", err)
	}
	tmd := md.Topics[topic]
	if tmd.Error.Code() != confluentKafka.ErrNoError {
		return fmt.Errorf("topic metadata error: %v", tmd.Error)
	}

	var tps []confluentKafka.TopicPartition
	for _, p := range tmd.Partitions {
		tps = append(tps, confluentKafka.TopicPartition{
			Topic:     &topic,
			Partition: p.ID,
			Offset:    confluentKafka.OffsetEnd,
		})
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	res, err := admin.DeleteRecords(ctx, tps)
	if err != nil {
		return fmt.Errorf("DeleteRecords request failed: %w", err)
	}

	for _, r := range res.DeleteRecordsResults {
		if r.TopicPartition.Error != nil {
			return fmt.Errorf("partition %s: %v", r.TopicPartition, r.TopicPartition.Error)
		}
	}
	return nil
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

	for _, topic := range topics {
		if err := m.PurgeTopic(ctx, topic); err != nil {
			return err
		}
	}

	return nil
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

func NewTestMessaging(ctx context.Context, tCfg *Config) (*TestMessaging, error) {
	switch tCfg.Mode {
	case ModeReal:
		mCfg, err := messaging.NewConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to create new messaging config: %w", err)
		}
		return &TestMessaging{Cfg: mCfg}, nil
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

		mCfg := &messaging.Config{
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

		return &TestMessaging{Cfg: mCfg, container: container}, nil
	}
}
