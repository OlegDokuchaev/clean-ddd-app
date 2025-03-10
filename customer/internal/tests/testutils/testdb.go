package testutils

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"customer/internal/infrastructure/db"
	"customer/internal/infrastructure/db/migrations"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestDB struct {
	DB         *gorm.DB
	Container  testcontainers.Container
	Migrations *migrate.Migrate
}

func SetupTestDatabase(ctx context.Context, config *db.Config) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14",
		ExposedPorts: []string{config.Port},
		Env: map[string]string{
			"POSTGRES_USER":     config.Username,
			"POSTGRES_PASSWORD": config.Password,
			"POSTGRES_DB":       config.Database,
		},
		WaitingFor: wait.ForListeningPort(nat.Port(config.Port)),
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func CreateDSN(ctx context.Context, container testcontainers.Container, config *db.Config) (string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := container.MappedPort(ctx, nat.Port(config.Port))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port.Int(), config.Username, config.Password, config.Database,
	), nil
}

func InitGORM(dsn string) (*gorm.DB, error) {
	return gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
}

func SetupMigrations(gormDB *gorm.DB, config *migrations.Config) (*migrate.Migrate, error) {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{})
	if err != nil {
		return nil, err
	}

	return migrate.NewWithDatabaseInstance(
		"file://"+config.MigrationsPath,
		"postgres",
		driver,
	)
}

func NewTestDB(ctx context.Context, config *db.Config, mConfig *migrations.Config) (*TestDB, error) {
	container, err := SetupTestDatabase(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}

	dsn, err := CreateDSN(ctx, container, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create DSN: %w", err)
	}

	gormDB, err := InitGORM(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init GORM: %w", err)
	}

	setupMigrations, err := SetupMigrations(gormDB, mConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to setup migrations: %w", err)
	}

	return &TestDB{
		DB:         gormDB,
		Container:  container,
		Migrations: setupMigrations,
	}, nil
}
