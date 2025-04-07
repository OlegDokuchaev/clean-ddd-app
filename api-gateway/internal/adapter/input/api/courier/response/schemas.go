package response

import "github.com/google/uuid"

type RegisterResponse struct {
	CourierID uuid.UUID `json:"courier_id"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
