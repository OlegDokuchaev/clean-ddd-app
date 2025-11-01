package customer_response

import "github.com/google/uuid"

type RegisterResponse struct {
	CustomerID uuid.UUID `json:"customer_id"`
}

type LoginResponse struct {
	ChallengeID string `json:"challenge_id"`
}

type VerifyOtpResponse struct {
	Token string `json:"token"`
}
