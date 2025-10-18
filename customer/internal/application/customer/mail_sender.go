package customer

import (
	"context"
)

type MailSender interface {
	SendOtp(ctx context.Context, toEmail string, code string) error
	SendPasswordResetLink(ctx context.Context, toEmail string, token string) error
}
