package response

import (
	customerv1 "customer/internal/presentation/grpc"

	"github.com/google/uuid"
)

func ToRegisterResponse(customerID uuid.UUID) *customerv1.RegisterResponse {
	return &customerv1.RegisterResponse{
		CustomerId: customerID.String(),
	}
}

func ToLoginResponse(challengeID string) *customerv1.LoginResponse {
	return &customerv1.LoginResponse{
		ChallengeId: challengeID,
	}
}

func ToVerifyOtpResponse(token string) *customerv1.VerifyOtpResponse {
	return &customerv1.VerifyOtpResponse{
		Token: token,
	}
}

func ToAuthenticateResponse(customerID uuid.UUID) *customerv1.AuthenticateResponse {
	return &customerv1.AuthenticateResponse{
		CustomerId: customerID.String(),
	}
}
