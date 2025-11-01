package customer_request

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyOtpRequest struct {
	Code string `json:"code" binding:"required"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required"`
}

type CompletePasswordReset struct {
	NewPassword string `json:"new_password" binding:"required"`
}
