package customer

import (
	"api-gateway/internal/adapter/input/api/customer/request"
	"api-gateway/internal/adapter/input/api/customer/response"
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
