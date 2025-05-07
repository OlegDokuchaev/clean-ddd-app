package di

import (
	"order/internal/infrastructure/db"

	"go.uber.org/fx"
)

var DatabaseModule = fx.Provide(
	// Database configuration
	db.NewConfig,

	// Database connection
	db.NewConnection,

	// Database
	db.NewDB,

	// Order collection
	db.NewOrderCollection,
)
