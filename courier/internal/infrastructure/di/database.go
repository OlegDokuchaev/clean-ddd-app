package di

import (
	"courier/internal/infrastructure/db"

	"go.uber.org/fx"
)

var DatabaseModule = fx.Provide(
	// Database configuration
	db.NewConfig,

	// Database connection
	db.NewDB,
)
