package api

import (
	courierApi "api-gateway/internal/adapter/input/api/courier"
	customerApi "api-gateway/internal/adapter/input/api/customer"
	orderApi "api-gateway/internal/adapter/input/api/order"
	warehouseApi "api-gateway/internal/adapter/input/api/warehouse"

	"github.com/gin-gonic/gin"
)

func NewAPI(
	orderHandler *orderApi.Handler,
	customerHandler *customerApi.Handler,
	courierHandler *courierApi.Handler,
	warehouseHandler *warehouseApi.Handler,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("")
	{
		orderApi.RegisterRoutes(api, orderHandler)
		customerApi.RegisterRoutes(api, customerHandler)
		courierApi.RegisterRoutes(api, courierHandler)
		warehouseApi.RegisterRoutes(api, warehouseHandler)
	}

	return r
}
