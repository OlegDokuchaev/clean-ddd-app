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
// @Router /customers [post]
func (h *Handler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := commonRequest.ParseInput(c, &req, binding.JSON); err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	data := request.ToRegisterDto(&req)
	customerID, err := h.uc.Register(c, data)
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
	var req request.LoginRequest
	if err := commonRequest.ParseInput(c, &req, binding.JSON); err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	data := request.ToLoginDto(&req)
	token, err := h.uc.Login(c, data)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		Token: token,
	})
}
