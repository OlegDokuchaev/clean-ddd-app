package courier

import (
	request "api-gateway/internal/adapter/input/api/courier/request"
	response "api-gateway/internal/adapter/input/api/courier/response"
	commonRequest "api-gateway/internal/adapter/input/api/request"
	commonResponse "api-gateway/internal/adapter/input/api/response"
	courierUseCase "api-gateway/internal/domain/usecases/courier"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Handler struct {
	uc courierUseCase.UseCase
}

func NewHandler(courierUseCase courierUseCase.UseCase) *Handler {
	return &Handler{uc: courierUseCase}
}

// Register godoc
// @Summary Register new courier
// @Description Register a new courier with name, password and phone
// @Tags couriers
// @Accept json
// @Produce json
// @Param request body courier_request.RegisterRequest true "Courier registration data"
// @Success 201 {object} courier_response.RegisterResponse "Courier created successfully"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request format"
// @Failure 409 {object} response.ErrorResponseDetail "Courier with this phone already exists"
// @Failure 422 {object} response.ErrorResponseDetail "Invalid data format"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Router /couriers/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := commonRequest.ParseInput(c, &req, binding.JSON); err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	data := request.ToRegisterDto(&req)
	courierID, err := h.uc.Register(c, data)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.RegisterResponse{
		CourierID: courierID,
	})
}

// Login godoc
// @Summary Courier login
// @Description Authenticate a courier and get a JWT token
// @Tags couriers
// @Accept json
// @Produce json
// @Param request body courier_request.LoginRequest true "Courier login credentials"
// @Success 200 {object} courier_response.LoginResponse "Login successful"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request format"
// @Failure 401 {object} response.ErrorResponseDetail "Invalid credentials"
// @Failure 404 {object} response.ErrorResponseDetail "Courier not found"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Router /couriers/login [post]
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
