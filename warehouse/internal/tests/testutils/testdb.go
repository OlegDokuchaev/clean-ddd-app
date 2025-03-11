package testutils

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"warehouse/internal/infrastructure/db/migrations"

	"github.com/testcontainers/testcontainers-go"
	postgresContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	TestDbName = "name"
	TestDbUser = "user"
	TestDbPass = "password"
)

type TestDB struct {
	DB        *gorm.DB
	Container testcontainers.Container
}

func SetupTestDatabase(ctx context.Context) (testcontainers.Container, error) {
	return postgresContainer.Run(ctx,
		"postgres:16-alpine",
		postgresContainer.WithDatabase(TestDbName),
		postgresContainer.WithUsername(TestDbUser),
		postgresContainer.WithPassword(TestDbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
}

func CreateDSN(ctx context.Context, container testcontainers.Container) (string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port.Int(), TestDbUser, TestDbPass, TestDbName,
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
		TestDbPass,
		driver,
	)
}

func NewTestDB(ctx context.Context, mConfig *migrations.Config) (*TestDB, error) {
	container, err := SetupTestDatabase(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}
	defer func() {
		if err != nil {
			_ = container.Terminate(ctx)
		}
	}()

	dsn, err := CreateDSN(ctx, container)
	if err != nil {
		return nil, fmt.Errorf("failed to create DSN: %w", err)
	}

	db, err := InitGORM(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init GORM: %w", err)
	}

	migration, err := SetupMigrations(db, mConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to setup migrations: %w", err)
	}
	if err = migration.Up(); err != nil {
		return nil, fmt.Errorf("failed to up migrations: %w", err)
	}

	return &TestDB{
		DB:        db,
		Container: container,
	}, nil
}
