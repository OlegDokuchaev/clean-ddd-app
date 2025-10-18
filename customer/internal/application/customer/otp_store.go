package customer

import (
	"context"
	"github.com/google/uuid"
)

type OtpStore interface {
	Issue(ctx context.Context, challengeID string, consumerID uuid.UUID, code string, otpPolicy OtpPolicy) error
	VerifyAndConsume(ctx context.Context, challengeID string, code string) (ok bool, attemptsLeft int, expired bool, consumerID uuid.UUID, err error)
	Invalidate(ctx context.Context, challengeID string) error
}
