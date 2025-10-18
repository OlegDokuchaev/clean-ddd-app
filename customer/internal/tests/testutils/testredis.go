package testutils

import (
	"context"
	otpStore "customer/internal/infrastructure/otp_store"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

const (
	testRedisPassword = "password"
	testOTPKeyPrefix  = "otp:"
)

type TestRedis struct {
	Client    *redis.Client
	Cfg       *otpStore.Config
	container testcontainers.Container
}

func (r *TestRedis) Close(ctx context.Context) error {
	if r.container == nil {
		return nil
	}
	return r.container.Terminate(ctx)
}

func (r *TestRedis) Clear(ctx context.Context) error {
	iter := r.Client.Scan(ctx, 0, r.Cfg.KeyPrefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		if err := r.Client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("failed to delete key %s: %w", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to iterate keys: %w", err)
	}
	return nil
}

func setupRedisContainer(ctx context.Context) (testcontainers.Container, error) {
	return redisContainer.Run(ctx, "redis:7-alpine")
}

func createRedisClientFromContainer(ctx context.Context, container testcontainers.Container, prefix string) (*redis.Client, *otpStore.Config, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	port, err := container.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return nil, nil, err
	}

	cfg := &otpStore.Config{
		Addr:         fmt.Sprintf("%s:%s", host, port.Port()),
		Password:     testRedisPassword,
		DB:           0,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		KeyPrefix:    prefix,
	}
	client := otpStore.NewClient(cfg)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to ping redis: %w", err)
	}
	return client, cfg, nil
}

func NewTestRedis(ctx context.Context, tCfg *Config) (*TestRedis, error) {
	switch tCfg.Mode {
	case ModeReal:
		cfg, err := otpStore.NewConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to load redis config: %w", err)
		}

		client := otpStore.NewClient(cfg)
		if err := client.Ping(ctx).Err(); err != nil {
			return nil, fmt.Errorf("failed to ping redis: %w", err)
		}

		return &TestRedis{Client: client, Cfg: cfg}, nil
	default:
		container, err := setupRedisContainer(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to setup redis: %w", err)
		}

		client, cfg, err := createRedisClientFromContainer(ctx, container, testOTPKeyPrefix)
		if err != nil {
			_ = container.Terminate(ctx)
			return nil, err
		}

		return &TestRedis{Client: client, Cfg: cfg, container: container}, nil
	}
}
