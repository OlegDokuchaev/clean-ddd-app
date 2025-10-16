package testutils

import (
	"context"
	"customer/internal/infrastructure/db"
	"customer/internal/infrastructure/db/migrations"
	"customer/internal/infrastructure/db/tables"
	"fmt"
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
	DB  *gorm.DB
	Cfg *db.Config

	container testcontainers.Container
}

func (d *TestDB) Clear(ctx context.Context) error {
	return d.DB.
		WithContext(ctx).
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&tables.Customer{}).
		Error
}

func (d *TestDB) Close(ctx context.Context) error {
	if d.container != nil {
		return d.container.Terminate(ctx)
	}
	return nil
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

func createDSN(c *db.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.Database,
	)
}

func initGORM(dsn string) (*gorm.DB, error) {
	return gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
}

func createGORM(dbCfg *db.Config) (*gorm.DB, error) {
	dsn := createDSN(dbCfg)

	sqlDB, err := initGORM(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init GORM: %w", err)
	}

	return sqlDB, nil
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

func createDB(dbCfg *db.Config, mCfg *migrations.Config) (*gorm.DB, error) {
	sqlDB, err := createGORM(dbCfg)
	if err != nil {
		return nil, err
	}

	if err = createMigrations(sqlDB, mCfg); err != nil {
		return nil, err
	}

	return sqlDB, nil
}

func NewTestDB(ctx context.Context, tCfg *Config, mCfg *migrations.Config) (*TestDB, error) {
	switch tCfg.Mode {
	case ModeReal:
		dbCfg, err := db.NewConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to load db config: %w", err)
		}

		sqlDB, err := createDB(dbCfg, mCfg)
		if err != nil {
			return nil, err
		}

		return &TestDB{
			DB:  sqlDB,
			Cfg: dbCfg,
		}, nil
	default:
		container, err := setupDBContainer(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to setup database: %w", err)
		}

		host, err := container.Host(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get database host: %w", err)
		}

		port, err := container.MappedPort(ctx, "5432/tcp")
		if err != nil {
			return nil, fmt.Errorf("failed to get database port: %w", err)
		}

		dbCfg := &db.Config{
			Host:     host,
			Port:     port.Port(),
			Database: TestDbName,
			Username: TestDbUser,
			Password: TestDbPass,
		}

		sqlDB, err := createDB(dbCfg, mCfg)
		if err != nil {
			if err := container.Terminate(ctx); err != nil {
				return nil, fmt.Errorf("failed to terminate database: %w", err)
			}
			return nil, err
		}

		return &TestDB{
			DB:        sqlDB,
			Cfg:       dbCfg,
			container: container,
		}, nil
	}
}
