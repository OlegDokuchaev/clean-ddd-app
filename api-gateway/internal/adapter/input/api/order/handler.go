package order

import (
	request "api-gateway/internal/adapter/input/api/order/request"
	response "api-gateway/internal/adapter/input/api/order/response"
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

// Create godoc
// @Summary Create a new order
// @Description Create a new order with items
// @Tags orders
// @Accept json
// @Produce json
// @Param request body order_request.CreateRequest true "Order details"
// @Success 201 "" "Created with location header"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request format"
// @Failure 401 {object} response.ErrorResponseDetail "Missing or invalid bearer token"
// @Failure 422 {object} response.ErrorResponseDetail "Invalid item data"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Security CustomerBearerAuth
// @Router /orders [post]
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

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel an order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 204 "" "No Content"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request"
// @Failure 401 {object} response.ErrorResponseDetail "Missing or invalid bearer token"
// @Failure 404 {object} response.ErrorResponseDetail "Order not found"
// @Failure 422 {object} response.ErrorResponseDetail "Invalid order ID format"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Security CustomerBearerAuth
// @Router /orders/{id}/cancel [patch]
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

// CompleteDelivery godoc
// @Summary Complete order
// @Description Mark an order as delivered (completed)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 "" "OK"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request"
// @Failure 401 {object} response.ErrorResponseDetail "Missing or invalid bearer token"
// @Failure 404 {object} response.ErrorResponseDetail "Order not found"
// @Failure 422 {object} response.ErrorResponseDetail "Invalid order ID format"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Security CourierBearerAuth
// @Router /orders/{id}/complete [patch]
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

	c.Status(http.StatusOK)
}

// GetCustomerOrders godoc
// @Summary Get customer orders
// @Description Get all orders for the authenticated customer
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {object} order_response.OrdersResponse "List of orders"
// @Failure 401 {object} response.ErrorResponseDetail "Missing or invalid bearer token"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Security CustomerBearerAuth
// @Router /orders [get]
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

// GetCourierOrders godoc
// @Summary Get courier orders
// @Description Get all current orders for the authenticated courier
// @Tags couriers
// @Accept json
// @Produce json
// @Success 200 {object} order_response.OrdersResponse "List of orders"
// @Failure 401 {object} response.ErrorResponseDetail "Missing or invalid bearer token"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Security CourierBearerAuth
// @Router /couriers/me/orders [get]
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
