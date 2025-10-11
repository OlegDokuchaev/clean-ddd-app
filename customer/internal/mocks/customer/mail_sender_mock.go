package customer

import (
	"context"
	customerApplication "customer/internal/application/customer"
	"github.com/stretchr/testify/mock"
)

type MailSenderMock struct {
	mock.Mock
}

func (m *MailSenderMock) SendOtp(ctx context.Context, toEmail string, code string) error {
	args := m.Called(ctx, toEmail, code)
	return args.Error(0)
}

func (m *MailSenderMock) SendPasswordResetLink(ctx context.Context, toEmail string, token string) error {
	args := m.Called(ctx, toEmail, token)
	return args.Error(0)
}

var _ customerApplication.MailSender = (*MailSenderMock)(nil)
