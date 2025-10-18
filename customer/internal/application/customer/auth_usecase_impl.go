package customer

import (
	"context"
	customerDomain "customer/internal/domain/customer"
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/google/uuid"
)

type AuthUseCaseImpl struct {
	repo         customerDomain.Repository
	tokenManager TokenManager
	mailSender   MailSender
	otpStore     OtpStore
	lockout      customerDomain.LockoutPolicy
	otpPolicy    OtpPolicy
}

func NewAuthUseCase(
	repo customerDomain.Repository,
	tokenManager TokenManager,
	mailSender MailSender,
	otpStore OtpStore,
	lockout customerDomain.LockoutPolicy,
	otpPolicy OtpPolicy,
) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		repo:         repo,
		tokenManager: tokenManager,
		mailSender:   mailSender,
		otpStore:     otpStore,
		lockout:      lockout,
		otpPolicy:    otpPolicy,
	}
}

func (u *AuthUseCaseImpl) Register(ctx context.Context, data RegisterDto) (uuid.UUID, error) {
	customer, err := customerDomain.Create(data.Name, data.Phone, data.Email, data.Password)
	if err != nil {
		return uuid.Nil, err
	}

	if err = u.repo.Create(ctx, customer); err != nil {
		return uuid.Nil, err
	}

	return customer.ID, nil
}

func (u *AuthUseCaseImpl) Login(ctx context.Context, data LoginDto) (string, error) {
	customer, err := u.repo.GetByPhone(ctx, data.Phone)
	if err != nil {
		return "", err
	}
	if customer.IsLocked() {
		return "", customerDomain.ErrLocked
	}

	if ok := customer.CheckPassword(data.Password); !ok {
		customer.RegisterFailedAttempt(u.lockout)
		if saveErr := u.repo.Save(ctx, customer); saveErr != nil {
			return "", saveErr
		}
		return "", customerDomain.ErrInvalidCustomerPassword
	}

	customer.ResetFailedAttempts()

	challengeID := uuid.NewString()
	code := GenerateOtpCode(6)

	if err := u.otpStore.Issue(ctx, challengeID, customer.ID, code, u.otpPolicy); err != nil {
		return "", err
	}
	if err := u.mailSender.SendOtp(ctx, customer.Email, code); err != nil {
		_ = u.otpStore.Invalidate(ctx, challengeID)
		return "", err
	}

	if err := u.repo.Save(ctx, customer); err != nil {
		return "", err
	}
	return challengeID, nil
}

func (u *AuthUseCaseImpl) VerifyOtp(ctx context.Context, data VerifyOtpDto) (string, error) {
	ok, attemptsLeft, expired, consumerID, err := u.otpStore.VerifyAndConsume(ctx, data.ChallengeID, data.Code)
	if err != nil {
		return "", err
	}
	if expired {
		return "", ErrOtpExpired
	}
	if !ok {
		if attemptsLeft <= 0 {
			_ = u.otpStore.Invalidate(ctx, data.ChallengeID)
			return "", ErrOtpAttemptsExceeded
		}
		return "", ErrOtpInvalid
	}

	_ = u.otpStore.Invalidate(ctx, data.ChallengeID)

	consumer, err := u.repo.GetByID(ctx, consumerID)
	if err != nil {
		return "", err
	}

	token, err := u.tokenManager.GenerateAccess(consumer.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *AuthUseCaseImpl) RequestPasswordReset(ctx context.Context, email string) error {
	consumer, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil
	}

	resetToken, err := u.tokenManager.GenerateReset(consumer.ID)
	if err != nil {
		return err
	}

	return u.mailSender.SendPasswordResetLink(ctx, consumer.Email, resetToken)
}

func (u *AuthUseCaseImpl) CompletePasswordReset(ctx context.Context, token string, newPassword string) error {
	consumerID, err := u.tokenManager.ParseAndValidateReset(token)
	if err != nil {
		return err
	}

	consumer, err := u.repo.GetByID(ctx, consumerID)
	if err != nil {
		return err
	}
	if err := consumer.SetPassword(newPassword); err != nil {
		return err
	}

	consumer.ResetFailedAttempts()

	return u.repo.Save(ctx, consumer)
}

func (u *AuthUseCaseImpl) Authenticate(_ context.Context, token string) (uuid.UUID, error) {
	return u.tokenManager.ParseAndValidateAccess(token)
}

func GenerateOtpCode(digits int) string {
	n := int(math.Pow10(digits))
	v := rand.IntN(n)
	return fmt.Sprintf("%0*d", digits, v)
}

var _ AuthUseCase = (*AuthUseCaseImpl)(nil)
