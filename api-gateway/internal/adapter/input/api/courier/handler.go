package courier

import (
	"api-gateway/internal/adapter/input/api/courier/request"
	"api-gateway/internal/adapter/input/api/courier/response"
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
