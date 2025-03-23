package response

import (
	courierv1 "courier/internal/presentation/grpc"

	"github.com/google/uuid"
)

func ToRegisterResponse(courierID uuid.UUID) *courierv1.RegisterResponse {
	return &courierv1.RegisterResponse{
		CourierId: courierID.String(),
	}
}

func ToLoginResponse(token string) *courierv1.LoginResponse {
	return &courierv1.LoginResponse{
		Token: token,
	}
}

func ToAuthenticateResponse(courierID uuid.UUID) *courierv1.AuthenticateResponse {
	return &courierv1.AuthenticateResponse{
		CourierId: courierID.String(),
	}
}
