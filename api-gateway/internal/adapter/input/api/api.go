package api

import (
	courierApi "api-gateway/internal/adapter/input/api/courier"
	customerApi "api-gateway/internal/adapter/input/api/customer"
	"api-gateway/internal/adapter/input/api/docs"
	"api-gateway/internal/adapter/input/api/middleware"
	orderApi "api-gateway/internal/adapter/input/api/order"
	warehouseApi "api-gateway/internal/adapter/input/api/warehouse"
	"api-gateway/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

func NewAPI(
	orderHandler *orderApi.Handler,
	customerHandler *customerApi.Handler,
	courierHandler *courierApi.Handler,
	warehouseHandler *warehouseApi.Handler,
	log logger.Logger,
) *gin.Engine {
	r := gin.New()

	r.Use(middleware.GinLoggingMiddleware(log))
	r.Use(gin.Recovery())

	api := r.Group("")
	{
		orderApi.RegisterRoutes(api, orderHandler)
		customerApi.RegisterRoutes(api, customerHandler)
		courierApi.RegisterRoutes(api, courierHandler)
		warehouseApi.RegisterRoutes(api, warehouseHandler)
		docs.RegisterRoutes(api)
	}

	return r
}
