package customer

import (
	"context"
	customerApplication "customer/internal/application/customer"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type OtpStoreMock struct {
	mock.Mock
}

func (m *OtpStoreMock) Issue(
	ctx context.Context,
	challengeID string,
	consumerID uuid.UUID,
	code string,
	otpPolicy customerApplication.OtpPolicy,
) error {
	args := m.Called(ctx, challengeID, consumerID, code, otpPolicy)
	return args.Error(0)
}

func (m *OtpStoreMock) VerifyAndConsume(
	ctx context.Context,
	challengeID string,
	code string,
) (bool, int, bool, uuid.UUID, error) {
	args := m.Called(ctx, challengeID, code)
	return args.Bool(0), args.Int(1), args.Bool(2), args.Get(3).(uuid.UUID), args.Error(4)
}

func (m *OtpStoreMock) Invalidate(ctx context.Context, challengeID string) error {
	args := m.Called(ctx, challengeID)
	return args.Error(0)
}

var _ customerApplication.OtpStore = (*OtpStoreMock)(nil)
