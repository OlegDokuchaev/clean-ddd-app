package customer

import (
	"context"
	customerApplication "customer/internal/application/customer"
	customerDomain "customer/internal/domain/customer"
	customerMock "customer/internal/mocks/customer"
	"customer/internal/tests/testutils/mothers"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *AuthUseCaseTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *AuthUseCaseTestSuite) TestRegister() {
	tests := []struct {
		name        string
		data        customerApplication.RegisterDto
		setup       func(repo *customerMock.RepositoryMock)
		expectedErr error
	}{
		{
			name: "Success",
			data: customerApplication.RegisterDto{
				Name:     "test",
				Phone:    "+79032895555",
				Email:    "test@example.com",
				Password: "password",
			},
			setup: func(repo *customerMock.RepositoryMock) {
				repo.On("Create", s.ctx, mock.Anything).Return(nil)
			},
		},
		{
			name: "Failure: invalid phone",
			data: customerApplication.RegisterDto{
				Name:     "test",
				Phone:    "invalid phone",
				Email:    "test@example.com",
				Password: "password",
			},
			setup:       func(repo *customerMock.RepositoryMock) {},
			expectedErr: customerDomain.ErrInvalidCustomerPhone,
		},
		{
			name: "Failure: repository create error",
			data: customerApplication.RegisterDto{
				Name:     "test",
				Phone:    "+79032895555",
				Email:    "test@example.com",
				Password: "password",
			},
			setup: func(repo *customerMock.RepositoryMock) {
				repo.On("Create", s.ctx, mock.Anything).
					Return(errors.New("create customer error"))
			},
			expectedErr: errors.New("create customer error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(customerMock.RepositoryMock)
			tokenManager := new(customerMock.TokenManagerMock)
			mailSender := new(customerMock.MailSenderMock)
			otpStore := new(customerMock.OtpStoreMock)

			lockout := customerDomain.LockoutPolicy{MaxFailed: 3, LockFor: time.Minute * 30}
			otpPolicy := customerApplication.OtpPolicy{TTL: 5 * time.Minute, MaxAttempts: 3}
			uc := customerApplication.NewAuthUseCase(
				repo, tokenManager, mailSender, otpStore,
				lockout, otpPolicy,
			)
			tc.setup(repo)

			customerID, err := uc.Register(s.ctx, tc.data)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.NotEqual(s.T(), uuid.Nil, customerID)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			repo.AssertExpectations(s.T())
		})
	}
}

func (s *AuthUseCaseTestSuite) TestLogin() {
	tests := []struct {
		name        string
		data        customerApplication.LoginDto
		setup       func(repo *customerMock.RepositoryMock, mail *customerMock.MailSenderMock, otp *customerMock.OtpStoreMock)
		expectedErr error
	}{
		{
			name: "Success",
			data: customerApplication.LoginDto{
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *customerMock.RepositoryMock, mail *customerMock.MailSenderMock, otp *customerMock.OtpStoreMock) {
				user := mothers.CustomerWithLogin("+79032895555", "password")

				repo.On("GetByPhone", s.ctx, user.Phone).Return(user, nil)
				otp.On("Issue", s.ctx, mock.Anything, user.ID, mock.Anything, mock.Anything).Return(nil)
				mail.On("SendOtp", s.ctx, user.Email, mock.Anything, mock.Anything).Return(nil)
				repo.On("Save", s.ctx, mock.Anything).Return(nil)
			},
		},
		{
			name: "Failure: repo.GetByPhone error",
			data: customerApplication.LoginDto{
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *customerMock.RepositoryMock, _ *customerMock.MailSenderMock, _ *customerMock.OtpStoreMock) {
				repo.On("GetByPhone", s.ctx, "+79032895555").
					Return((*customerDomain.Customer)(nil), errors.New("get by phone error"))
			},
			expectedErr: errors.New("get by phone error"),
		},
		{
			name: "Failure: user locked",
			data: customerApplication.LoginDto{
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *customerMock.RepositoryMock, _ *customerMock.MailSenderMock, _ *customerMock.OtpStoreMock) {
				user := mothers.LockedCustomerWithLogin("+79032895555", "password")
				repo.On("GetByPhone", s.ctx, user.Phone).Return(user, nil)
			},
			expectedErr: customerDomain.ErrLocked,
		},
		{
			name: "Failure: wrong password (increments failed & saves)",
			data: customerApplication.LoginDto{
				Phone:    "+79032895555",
				Password: "wrong",
			},
			setup: func(repo *customerMock.RepositoryMock, _ *customerMock.MailSenderMock, _ *customerMock.OtpStoreMock) {
				user := mothers.CustomerWithLogin("+79032895555", "password")

				repo.On("GetByPhone", s.ctx, user.Phone).Return(user, nil)
				repo.On("Save", s.ctx, mock.Anything).Return(nil)
			},
			expectedErr: customerDomain.ErrInvalidCustomerPassword,
		},
		{
			name: "Failure: otpStore.Issue error",
			data: customerApplication.LoginDto{
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *customerMock.RepositoryMock, _ *customerMock.MailSenderMock, otp *customerMock.OtpStoreMock) {
				user := mothers.CustomerWithLogin("+79032895555", "password")

				repo.On("GetByPhone", s.ctx, user.Phone).Return(user, nil)
				otp.On("Issue", s.ctx, mock.Anything, user.ID, mock.Anything, mock.Anything).
					Return(errors.New("issue error"))
			},
			expectedErr: errors.New("issue error"),
		},
		{
			name: "Failure: mailSender.SendOtp error",
			data: customerApplication.LoginDto{
				Phone:    "+79032895555",
				Password: "password",
			},
			setup: func(repo *customerMock.RepositoryMock, mail *customerMock.MailSenderMock, otp *customerMock.OtpStoreMock) {
				user := mothers.CustomerWithLogin("+79032895555", "password")

				repo.On("GetByPhone", s.ctx, user.Phone).Return(user, nil)
				otp.On("Issue", s.ctx, mock.Anything, user.ID, mock.Anything, mock.Anything).Return(nil)
				mail.On("SendOtp", s.ctx, user.Email, mock.Anything, mock.Anything).
					Return(errors.New("send error"))
				otp.On("Invalidate", s.ctx, mock.Anything).Return(nil)
			},
			expectedErr: errors.New("send error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()

			repo := new(customerMock.RepositoryMock)
			tokenManager := new(customerMock.TokenManagerMock)
			mailSender := new(customerMock.MailSenderMock)
			otpStore := new(customerMock.OtpStoreMock)

			lockout := customerDomain.LockoutPolicy{MaxFailed: 3, LockFor: time.Minute * 30}
			otpPolicy := customerApplication.OtpPolicy{TTL: 5 * time.Minute, MaxAttempts: 3}
			uc := customerApplication.NewAuthUseCase(
				repo, tokenManager, mailSender, otpStore,
				lockout, otpPolicy,
			)

			tc.setup(repo, mailSender, otpStore)

			challengeID, err := uc.Login(s.ctx, tc.data)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.NotEmpty(s.T(), challengeID)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			repo.AssertExpectations(s.T())
			tokenManager.AssertExpectations(s.T())
			mailSender.AssertExpectations(s.T())
			otpStore.AssertExpectations(s.T())
		})
	}
}

func (s *AuthUseCaseTestSuite) TestVerifyOtp() {
	tests := []struct {
		name        string
		data        customerApplication.VerifyOtpDto
		setup       func(otp *customerMock.OtpStoreMock, repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock)
		expectedErr error
	}{
		{
			name: "Success",
			data: customerApplication.VerifyOtpDto{
				ChallengeID: "challenge",
				Code:        "123456",
			},
			setup: func(otp *customerMock.OtpStoreMock, repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock) {
				customer := mothers.DefaultCustomer()

				otp.On("VerifyAndConsume", s.ctx, "challenge", "123456").
					Return(true, 2, false, customer.ID, nil)
				otp.On("Invalidate", s.ctx, "challenge").Return(nil)
				repo.On("GetByID", s.ctx, customer.ID).Return(customer, nil)
				token.On("GenerateAccess", customer.ID).
					Return("access-token", nil)
			},
		},
		{
			name: "Failure: OTP expired",
			data: customerApplication.VerifyOtpDto{ChallengeID: "c", Code: "1"},
			setup: func(otp *customerMock.OtpStoreMock, _ *customerMock.RepositoryMock, _ *customerMock.TokenManagerMock) {
				customerID := uuid.New()
				otp.On("VerifyAndConsume", s.ctx, "c", "1").
					Return(false, 0, true, customerID, nil)
			},
			expectedErr: customerApplication.ErrOtpExpired,
		},
		{
			name: "Failure: OTP invalid (attempts left)",
			data: customerApplication.VerifyOtpDto{ChallengeID: "c2", Code: "000000"},
			setup: func(otp *customerMock.OtpStoreMock, _ *customerMock.RepositoryMock, _ *customerMock.TokenManagerMock) {
				customerID := uuid.New()
				otp.On("VerifyAndConsume", s.ctx, "c2", "000000").
					Return(false, 2, false, customerID, nil)
			},
			expectedErr: customerApplication.ErrOtpInvalid,
		},
		{
			name: "Failure: attempts exceeded -> invalidated",
			data: customerApplication.VerifyOtpDto{ChallengeID: "c3", Code: "000000"},
			setup: func(otp *customerMock.OtpStoreMock, _ *customerMock.RepositoryMock, _ *customerMock.TokenManagerMock) {
				customerID := uuid.New()
				otp.On("VerifyAndConsume", s.ctx, "c3", "000000").
					Return(false, 0, false, customerID, nil)
				otp.On("Invalidate", s.ctx, "c3").Return(nil)
			},
			expectedErr: customerApplication.ErrOtpAttemptsExceeded,
		},
		{
			name: "Failure: get customer by id error",
			data: customerApplication.VerifyOtpDto{ChallengeID: "c4", Code: "123456"},
			setup: func(otp *customerMock.OtpStoreMock, repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock) {
				customerID := uuid.New()

				otp.On("VerifyAndConsume", s.ctx, "c4", "123456").
					Return(true, 3, false, customerID, nil)
				otp.On("Invalidate", s.ctx, "c4").Return(nil)
				repo.On("GetByID", s.ctx, customerID).Return((*customerDomain.Customer)(nil), errors.New("get by id error"))
			},
			expectedErr: errors.New("get by id error"),
		},
		{
			name: "Failure: token generation error",
			data: customerApplication.VerifyOtpDto{ChallengeID: "c4", Code: "123456"},
			setup: func(otp *customerMock.OtpStoreMock, repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock) {
				customer := mothers.DefaultCustomer()

				otp.On("VerifyAndConsume", s.ctx, "c4", "123456").
					Return(true, 3, false, customer.ID, nil)
				otp.On("Invalidate", s.ctx, "c4").Return(nil)
				repo.On("GetByID", s.ctx, customer.ID).Return(customer, nil)
				token.On("GenerateAccess", customer.ID).
					Return("", errors.New("gen token error"))
			},
			expectedErr: errors.New("gen token error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(customerMock.RepositoryMock)
			tokenManager := new(customerMock.TokenManagerMock)
			mailSender := new(customerMock.MailSenderMock)
			otpStore := new(customerMock.OtpStoreMock)

			lockout := customerDomain.LockoutPolicy{MaxFailed: 3, LockFor: time.Minute * 30}
			otpPolicy := customerApplication.OtpPolicy{TTL: 5 * time.Minute, MaxAttempts: 3}
			uc := customerApplication.NewAuthUseCase(
				repo, tokenManager, mailSender, otpStore,
				lockout, otpPolicy,
			)

			tc.setup(otpStore, repo, tokenManager)

			token, err := uc.VerifyOtp(s.ctx, tc.data)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.Equal(s.T(), "access-token", token)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			otpStore.AssertExpectations(s.T())
			tokenManager.AssertExpectations(s.T())
		})
	}
}

func (s *AuthUseCaseTestSuite) TestAuthenticate() {
	tests := []struct {
		name        string
		data        string
		setup       func(token *customerMock.TokenManagerMock) uuid.UUID
		expectedErr error
	}{
		{
			name: "Success",
			data: "token",
			setup: func(token *customerMock.TokenManagerMock) uuid.UUID {
				uid := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				token.On("ParseAndValidateAccess", "token").Return(uid, nil)
				return uid
			},
		},
		{
			name: "Failure: ParseAndValidateAccess error",
			data: "token",
			setup: func(token *customerMock.TokenManagerMock) uuid.UUID {
				token.On("ParseAndValidateAccess", "token").Return(uuid.Nil, errors.New("decode token error"))
				return uuid.Nil
			},
			expectedErr: errors.New("decode token error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(customerMock.RepositoryMock)
			tokenManager := new(customerMock.TokenManagerMock)
			mailSender := new(customerMock.MailSenderMock)
			otpStore := new(customerMock.OtpStoreMock)

			lockout := customerDomain.LockoutPolicy{MaxFailed: 3, LockFor: time.Minute * 30}
			otpPolicy := customerApplication.OtpPolicy{TTL: 5 * time.Minute, MaxAttempts: 3}
			uc := customerApplication.NewAuthUseCase(
				repo, tokenManager, mailSender, otpStore,
				lockout, otpPolicy,
			)

			expected := tc.setup(tokenManager)

			got, err := uc.Authenticate(s.ctx, tc.data)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.Equal(s.T(), expected, got)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			tokenManager.AssertExpectations(s.T())
		})
	}
}

func (s *AuthUseCaseTestSuite) TestRequestPasswordReset() {
	tests := []struct {
		name        string
		email       string
		setup       func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock, mail *customerMock.MailSenderMock)
		expectedErr error
	}{
		{
			name:  "Success",
			email: "user@example.com",
			setup: func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock, mail *customerMock.MailSenderMock) {
				user := mothers.CustomerWithEmail("user@example.com")
				repo.On("GetByEmail", s.ctx, "user@example.com").Return(user, nil)
				token.On("GenerateReset", user.ID).Return("reset-token", nil)
				mail.On("SendPasswordResetLink", s.ctx, "user@example.com", "reset-token").Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:  "Email not found -> still nil (no enumeration)",
			email: "absent@example.com",
			setup: func(repo *customerMock.RepositoryMock, _ *customerMock.TokenManagerMock, _ *customerMock.MailSenderMock) {
				repo.On("GetByEmail", s.ctx, "absent@example.com").Return((*customerDomain.Customer)(nil), errors.New("not found"))
			},
			expectedErr: nil,
		},
		{
			name:  "Failure: GenerateResetToken error",
			email: "user@example.com",
			setup: func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock, mail *customerMock.MailSenderMock) {
				user := mothers.CustomerWithEmail("user@example.com")
				repo.On("GetByEmail", s.ctx, "user@example.com").Return(user, nil)
				token.On("GenerateReset", user.ID).Return("", errors.New("gen reset token error"))
			},
			expectedErr: errors.New("gen reset token error"),
		},
		{
			name:  "Failure: SendPasswordResetLink error",
			email: "user@example.com",
			setup: func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock, mail *customerMock.MailSenderMock) {
				user := mothers.CustomerWithEmail("user@example.com")
				repo.On("GetByEmail", s.ctx, user.Email).Return(user, nil)
				token.On("GenerateReset", user.ID).Return("reset-token", nil)
				mail.On("SendPasswordResetLink", s.ctx, user.Email, "reset-token").
					Return(errors.New("send mail error"))
			},
			expectedErr: errors.New("send mail error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()

			repo := new(customerMock.RepositoryMock)
			tokenManager := new(customerMock.TokenManagerMock)
			mailSender := new(customerMock.MailSenderMock)
			otpStore := new(customerMock.OtpStoreMock)

			lockout := customerDomain.LockoutPolicy{MaxFailed: 3, LockFor: time.Minute * 30}
			otpPolicy := customerApplication.OtpPolicy{TTL: 5 * time.Minute, MaxAttempts: 3}
			uc := customerApplication.NewAuthUseCase(
				repo, tokenManager, mailSender, otpStore,
				lockout, otpPolicy,
			)

			tc.setup(repo, tokenManager, mailSender)

			err := uc.RequestPasswordReset(s.ctx, tc.email)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			repo.AssertExpectations(s.T())
			tokenManager.AssertExpectations(s.T())
			mailSender.AssertExpectations(s.T())
		})
	}
}

func (s *AuthUseCaseTestSuite) TestCompletePasswordReset() {
	tests := []struct {
		name        string
		token       string
		newPassword string
		setup       func(repo *customerMock.RepositoryMock, token *customerMock.TokenManagerMock)
		expectedErr error
	}{
		{
			name:        "Success",
			token:       "reset-token",
			newPassword: "N3wPass!",
			setup: func(repo *customerMock.RepositoryMock, tm *customerMock.TokenManagerMock) {
				user := mothers.LockedCustomer()

				tm.On("ParseAndValidateReset", "reset-token").Return(user.ID, nil)
				repo.On("GetByID", s.ctx, user.ID).Return(user, nil)
				repo.On("Save", s.ctx, mock.Anything).Return(nil)
			},
		},
		{
			name:        "Failure: invalid/expired reset token",
			token:       "bad",
			newPassword: "whatever",
			setup: func(_ *customerMock.RepositoryMock, tm *customerMock.TokenManagerMock) {
				tm.On("ParseAndValidateReset", "bad").Return(uuid.Nil, errors.New("invalid"))
			},
			expectedErr: errors.New("invalid"),
		},
		{
			name:        "Failure: repo.GetByID error",
			token:       "reset-token",
			newPassword: "N3wPass!",
			setup: func(repo *customerMock.RepositoryMock, tm *customerMock.TokenManagerMock) {
				uid := uuid.New()
				tm.On("ParseAndValidateReset", "reset-token").Return(uid, nil)
				repo.On("GetByID", s.ctx, uid).Return((*customerDomain.Customer)(nil), errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
		{
			name:        "Failure: invalid new password",
			token:       "reset-token",
			newPassword: "",
			setup: func(repo *customerMock.RepositoryMock, tm *customerMock.TokenManagerMock) {
				uid := uuid.New()
				tm.On("ParseAndValidateReset", "reset-token").Return(uid, nil)
				user := mothers.DefaultCustomer()
				user.ID = uid
				repo.On("GetByID", s.ctx, uid).Return(user, nil)
			},
			expectedErr: customerDomain.ErrInvalidCustomerPassword,
		},
		{
			name:        "Failure: repo.Save error",
			token:       "reset-token",
			newPassword: "N3wPass!",
			setup: func(repo *customerMock.RepositoryMock, tm *customerMock.TokenManagerMock) {
				uid := uuid.New()
				tm.On("ParseAndValidateReset", "reset-token").Return(uid, nil)
				user := mothers.DefaultCustomer()
				user.ID = uid
				repo.On("GetByID", s.ctx, uid).Return(user, nil)
				repo.On("Save", s.ctx, mock.Anything).Return(errors.New("save error"))
			},
			expectedErr: errors.New("save error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()

			repo := new(customerMock.RepositoryMock)
			tokenManager := new(customerMock.TokenManagerMock)
			mailSender := new(customerMock.MailSenderMock)
			otpStore := new(customerMock.OtpStoreMock)

			lockout := customerDomain.LockoutPolicy{MaxFailed: 3, LockFor: time.Minute * 30}
			otpPolicy := customerApplication.OtpPolicy{TTL: 5 * time.Minute, MaxAttempts: 3}
			uc := customerApplication.NewAuthUseCase(
				repo, tokenManager, mailSender, otpStore,
				lockout, otpPolicy,
			)

			tc.setup(repo, tokenManager)

			err := uc.CompletePasswordReset(s.ctx, tc.token, tc.newPassword)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			repo.AssertExpectations(s.T())
			tokenManager.AssertExpectations(s.T())
		})
	}
}

func TestAuthUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AuthUseCaseTestSuite))
}
