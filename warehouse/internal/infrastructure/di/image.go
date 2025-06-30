package di

import (
	productApplication "warehouse/internal/application/product"
	productImage "warehouse/internal/infrastructure/image/product"

	"go.uber.org/fx"
)

var ImageModule = fx.Provide(
	// Config
	productImage.NewConfig,

	// Client
	productImage.NewClient,

	// Service
	fx.Annotate(
		productImage.NewImageService,
		fx.As(new(productApplication.ImageService)),
	),
)
