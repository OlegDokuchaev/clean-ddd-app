package di

import (
	customerAplication "customer/internal/application/customer"
	otpStore "customer/internal/infrastructure/otp_store"
	"go.uber.org/fx"
)

var OtpStoreModule = fx.Provide(
	// Config
	otpStore.NewConfig,

	// Client
	otpStore.NewClient,

	// Otp store
	fx.Annotate(
		otpStore.New,
		fx.As(new(customerAplication.OtpStore)),
	),
)
