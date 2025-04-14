package di

import (
	productImage "warehouse/internal/infrastructure/image/product"

	"go.uber.org/fx"
)

var ImageModule = fx.Provide(
	// Config
	productImage.NewConfig,

	// Client
	productImage.NewClient,

	// Service
	productImage.NewImageService,
)
