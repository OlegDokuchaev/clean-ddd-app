package di

import (
	customerAplication "customer/internal/application/customer"
	mailSender "customer/internal/infrastructure/mail_sender"
	"go.uber.org/fx"
)

var MailSenderModule = fx.Provide(
	// Config
	mailSender.NewConfig,

	// Client
	mailSender.NewClient,

	// Otp store
	fx.Annotate(
		mailSender.New,
		fx.As(new(customerAplication.MailSender)),
	),
)
