package testutils

import (
	"context"
	"fmt"
	"order/internal/infrastructure/db/migrations"

	"github.com/golang-migrate/migrate/v4"
	migrateMongo "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	mongoContainer "github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	TestDbName         = "name"
	TestCollectionName = "orders"
)

type TestDB struct {
	Collection *mongo.Collection
	container  testcontainers.Container
}

func (d *TestDB) Close(ctx context.Context) error {
	return d.container.Terminate(ctx)
}

func setupDBContainer(ctx context.Context) (testcontainers.Container, error) {
	return mongoContainer.Run(ctx, "mongo:6")
}

func createDSN(ctx context.Context, container testcontainers.Container) (string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := container.MappedPort(ctx, "27017/tcp")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"mongodb://%s:%d",
		host, port.Int(),
	), nil
}

func createMongo(ctx context.Context, container testcontainers.Container) (*mongo.Client, error) {
	dsn, err := createDSN(ctx, container)
	if err != nil {
		return nil, err
	}

	clientOptions := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

func setupMigrations(db *mongo.Client, config *migrations.Config) (*migrate.Migrate, error) {
	driver, err := migrateMongo.WithInstance(db, &migrateMongo.Config{DatabaseName: TestDbName})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance(
		"file://"+config.MigrationsPath,
		TestDbName,
		driver,
	)
}

func createMigrations(db *mongo.Client, config *migrations.Config) error {
	migration, err := setupMigrations(db, config)
	if err != nil {
		return fmt.Errorf("failed to setup migrations: %w", err)
	}
	if err = migration.Up(); err != nil {
		return fmt.Errorf("failed to up migrations: %w", err)
	}
	return nil
}

func createDB(ctx context.Context, container testcontainers.Container, config *migrations.Config) (*mongo.Database, error) {
	client, err := createMongo(ctx, container)
	if err != nil {
		return nil, err
	}

	if err = createMigrations(client, config); err != nil {
		return nil, err
	}

	return client.Database(TestDbName), nil
}

func NewTestDB(ctx context.Context, config *migrations.Config) (*TestDB, error) {
	container, err := setupDBContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}

	db, err := createDB(ctx, container, config)
	if err != nil {
		if err := container.Terminate(ctx); err != nil {
			return nil, fmt.Errorf("failed to terminate mongo: %w", err)
		}
		return nil, err
	}

	collection := db.Collection(TestCollectionName)

	return &TestDB{
		Collection: collection,
		container:  container,
	}, nil
}
