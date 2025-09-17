package testutils

import (
	"context"
	"fmt"
	"order/internal/infrastructure/db"
	"order/internal/infrastructure/db/migrations"

	"github.com/golang-migrate/migrate/v4"
	migrateMongo "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	mongoContainer "github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	TestDbName              = "name"
	TestOrderCollectionName = "order"
)

type TestDB struct {
	DB  *mongo.Database
	Cfg *db.Config

	container testcontainers.Container
}

func (d *TestDB) Clear(ctx context.Context, collections ...string) error {
	if d.DB == nil {
		return nil
	}

	if len(collections) == 0 && d.Cfg != nil {
		collections = []string{d.Cfg.OrderCollection}
	}
	for _, col := range collections {
		if col == "" {
			continue
		}
		if _, err := d.DB.Collection(col).DeleteMany(ctx, bson.M{}); err != nil {
			return err
		}
	}
	return nil
}

func (d *TestDB) Close(ctx context.Context) error {
	if d.container != nil {
		return d.container.Terminate(ctx)
	}
	if _, err := d.DB.Collection(d.Cfg.OrderCollection).DeleteMany(ctx, bson.M{}); err != nil {
		return err
	}
	return d.DB.Client().Disconnect(ctx)
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

func setupMigrations(db *mongo.Client, dbName string, config *migrations.Config) (*migrate.Migrate, error) {
	driver, err := migrateMongo.WithInstance(db, &migrateMongo.Config{DatabaseName: dbName})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance(
		"file://"+config.MigrationsPath,
		dbName,
		driver,
	)
}

func createMigrations(db *mongo.Client, dbName string, config *migrations.Config) error {
	migration, err := setupMigrations(db, dbName, config)
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

	if err = createMigrations(client, TestDbName, config); err != nil {
		return nil, err
	}

	return client.Database(TestDbName), nil
}

func NewTestDB(ctx context.Context, tCfg *Config, mCfg *migrations.Config) (*TestDB, error) {
	switch tCfg.Mode {
	case ModeReal:
		dbCfg, err := db.NewConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to load db config: %w", err)
		}

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbCfg.URI))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to real mongo: %w", err)
		}
		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			return nil, fmt.Errorf("failed to ping real mongo: %w", err)
		}

		return &TestDB{DB: client.Database(dbCfg.Database), Cfg: dbCfg}, nil
	default:
		container, err := setupDBContainer(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to setup database: %w", err)
		}

		mongoDB, err := createDB(ctx, container, mCfg)
		if err != nil {
			if err := container.Terminate(ctx); err != nil {
				return nil, fmt.Errorf("failed to terminate mongo: %w", err)
			}
			return nil, err
		}

		dsn, err := createDSN(ctx, container)
		if err != nil {
			return nil, err
		}

		genCfg := &db.Config{
			URI:             dsn,
			Database:        TestDbName,
			OrderCollection: TestOrderCollectionName,
		}

		return &TestDB{
			DB:        mongoDB,
			Cfg:       genCfg,
			container: container,
		}, nil
	}
}
