package order

import (
	"api-gateway/internal/adapter/input/api/order/request"
	"api-gateway/internal/adapter/input/api/order/response"
	commonRequest "api-gateway/internal/adapter/input/api/request"
	commonResponse "api-gateway/internal/adapter/input/api/response"
	orderUseCase "api-gateway/internal/domain/usecases/order"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Handler struct {
	uc orderUseCase.UseCase
}

func NewHandler(orderUseCase orderUseCase.UseCase) *Handler {
	return &Handler{uc: orderUseCase}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateRequest
	if err := commonRequest.ParseInput(c, &req, binding.JSON); err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	token, err := commonRequest.ParseBearerToken(c)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	data := request.ToOrderCreateDto(&req)
	orderID, err := h.uc.Create(c, data, token)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	commonResponse.AddLocationHeaderWithID(c, orderID)
	c.Status(http.StatusCreated)
}

func (h *Handler) CancelOrder(c *gin.Context) {
	orderID, err := commonRequest.ParseParamUUID(c, "id")
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	token, err := commonRequest.ParseBearerToken(c)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	err = h.uc.CancelByCustomer(c, orderID, token)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) CompleteDelivery(c *gin.Context) {
	orderID, err := commonRequest.ParseParamUUID(c, "id")
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	token, err := commonRequest.ParseBearerToken(c)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	err = h.uc.Complete(c, orderID, token)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetCustomerOrders(c *gin.Context) {
	token, err := commonRequest.ParseBearerToken(c)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	orders, err := h.uc.GetByCustomer(c, token)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ToOrdersResponse(orders))
}

func (h *Handler) GetCourierOrders(c *gin.Context) {
	token, err := commonRequest.ParseBearerToken(c)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	orders, err := h.uc.GetCurrentByCourier(c, token)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ToOrdersResponse(orders))
}
