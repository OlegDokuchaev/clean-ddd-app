package customer

import (
	request "api-gateway/internal/adapter/input/api/customer/request"
	response "api-gateway/internal/adapter/input/api/customer/response"
	commonRequest "api-gateway/internal/adapter/input/api/request"
	commonResponse "api-gateway/internal/adapter/input/api/response"
	customerUseCase "api-gateway/internal/domain/usecases/customer"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Handler struct {
	uc customerUseCase.UseCase
}

func NewHandler(customerUseCase customerUseCase.UseCase) *Handler {
	return &Handler{uc: customerUseCase}
}

// Register godoc
// @Summary Register new customer
// @Description Register a new customer with name, password and phone
// @Tags customers
// @Accept json
// @Produce json
// @Param request body customer_request.RegisterRequest true "Customer registration data"
// @Success 201 {object} customer_response.RegisterResponse "Customer created successfully"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request format"
// @Failure 409 {object} response.ErrorResponseDetail "Customer with this phone already exists"
// @Failure 422 {object} response.ErrorResponseDetail "Invalid data format"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Router /customers/register [post]
func (h *Handler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.RegisterRequest
	if err := commonRequest.ParseInput(c, &req, binding.JSON); err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	data := request.ToRegisterDto(&req)
	customerID, err := h.uc.Register(ctx, data)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.RegisterResponse{
		CustomerID: customerID,
	})
}

// Login godoc
// @Summary Customer login
// @Description Authenticate a customer and get a JWT token
// @Tags customers
// @Accept json
// @Produce json
// @Param request body customer_request.LoginRequest true "Customer login credentials"
// @Success 200 {object} customer_response.LoginResponse "Login successful"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request format"
// @Failure 401 {object} response.ErrorResponseDetail "Invalid credentials"
// @Failure 404 {object} response.ErrorResponseDetail "Customer not found"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Router /customers/login [post]
func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.LoginRequest
	if err := commonRequest.ParseInput(c, &req, binding.JSON); err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	data := request.ToLoginDto(&req)
	challengeID, err := h.uc.Login(ctx, data)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		ChallengeID: challengeID,
	})
}

// VerifyOtp godoc
// @Summary Verify OTP code
// @Description Verify OTP code for the authentication challenge
// @Tags customers
// @Accept json
// @Produce json
// @Param challenge_id path string true "Challenge ID"
// @Param request body customer_request.VerifyOtpRequest true "OTP verification data"
// @Success 200 {object} customer_response.VerifyOtpResponse "OTP verified"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request format"
// @Failure 401 {object} response.ErrorResponseDetail "Invalid or expired code"
// @Failure 404 {object} response.ErrorResponseDetail "Challenge not found"
// @Failure 422 {object} response.ErrorResponseDetail "Invalid data format"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Router /customers/auth-challenges/{challenge_id} [patch]
func (h *Handler) VerifyOtp(c *gin.Context) {
	ctx := c.Request.Context()
	challengeID := c.Param("challenge_id")

	var req request.VerifyOtpRequest
	if err := commonRequest.ParseInput(c, &req, binding.JSON); err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	data := request.ToVerifyOtp(challengeID, &req)
	token, err := h.uc.VerifyOtp(ctx, data)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.VerifyOtpResponse{
		Token: token,
	})
}

// RequestPasswordReset godoc
// @Summary Request password reset
// @Description Send password reset email to customer
// @Tags customers
// @Accept json
// @Produce json
// @Param request body customer_request.RequestPasswordResetRequest true "Password reset request data"
// @Success 204 "Password reset email sent"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request format"
// @Failure 404 {object} response.ErrorResponseDetail "Customer not found"
// @Failure 422 {object} response.ErrorResponseDetail "Invalid data format"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Router /customers/password-resets [post]
func (h *Handler) RequestPasswordReset(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.RequestPasswordResetRequest
	if err := commonRequest.ParseInput(c, &req, binding.JSON); err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	err := h.uc.RequestPasswordReset(ctx, req.Email)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// CompletePasswordReset godoc
// @Summary Complete password reset
// @Description Reset customer password using token
// @Tags customers
// @Accept json
// @Produce json
// @Param token path string true "Password reset token"
// @Param request body customer_request.CompletePasswordReset true "New password data"
// @Success 204 "Password reset completed"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request format"
// @Failure 401 {object} response.ErrorResponseDetail "Invalid or expired token"
// @Failure 404 {object} response.ErrorResponseDetail "Token not found"
// @Failure 422 {object} response.ErrorResponseDetail "Invalid data format"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Router /customers/password-resets/{token} [patch]
func (h *Handler) CompletePasswordReset(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.Param("token")

	var req request.CompletePasswordReset
	if err := commonRequest.ParseInput(c, &req, binding.JSON); err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	err := h.uc.CompletePasswordReset(ctx, token, req.NewPassword)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
