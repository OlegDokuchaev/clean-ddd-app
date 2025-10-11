package customer

import (
	"github.com/google/uuid"
	"time"
)

type OtpPolicy struct {
	TTL         time.Duration
	MaxAttempts int
}

type RegisterDto struct {
	Name     string
	Email    string
	Phone    string
	Password string
}

type LoginDto struct {
	Phone    string
	Password string
}

type VerifyOtpDto struct {
	ChallengeID string
	Code        string
}

type ChangePasswordDto struct {
	UserID      uuid.UUID
	OldPassword string
	NewPassword string
}
