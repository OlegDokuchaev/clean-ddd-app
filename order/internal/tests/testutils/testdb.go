package testutils

import (
	"context"
	"fmt"
	"order/internal/infrastructure/db/migrations"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	Container testcontainers.Container
	DB        *gorm.DB
}

func setupDBContainer(ctx context.Context) (testcontainers.Container, error) {
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

func createDSN(ctx context.Context, container testcontainers.Container) (string, error) {
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

func initGORM(dsn string) (*gorm.DB, error) {
	return gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
}

func createGORM(ctx context.Context, container testcontainers.Container) (*gorm.DB, error) {
	dsn, err := createDSN(ctx, container)
	if err != nil {
		return nil, fmt.Errorf("failed to create DSN: %w", err)
	}

	db, err := initGORM(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init GORM: %w", err)
	}

	return db, nil
}

func setupMigrations(db *gorm.DB, config *migrations.Config) (*migrate.Migrate, error) {
	sqlDB, err := db.DB()
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

func createMigrations(db *gorm.DB, config *migrations.Config) error {
	migration, err := setupMigrations(db, config)
	if err != nil {
		return fmt.Errorf("failed to setup migrations: %w", err)
	}
	if err = migration.Up(); err != nil {
		return fmt.Errorf("failed to up migrations: %w", err)
	}
	return nil
}

func createDB(ctx context.Context, container testcontainers.Container, config *migrations.Config) (*gorm.DB, error) {
	db, err := createGORM(ctx, container)
	if err != nil {
		return nil, err
	}

	if err = createMigrations(db, config); err != nil {
		return nil, err
	}

	return db, nil
}

func NewTestDB(ctx context.Context, config *migrations.Config) (*TestDB, error) {
	container, err := setupDBContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}

	db, err := createDB(ctx, container, config)
	if err != nil {
		if err := container.Terminate(ctx); err != nil {
			return nil, fmt.Errorf("failed to terminate database: %w", err)
		}
		return nil, err
	}

	return &TestDB{
		DB:        db,
		Container: container,
	}, nil
}
