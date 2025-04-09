package order

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	orders := router.Group("/orders")
	{
		orders.POST("", handler.Create)
		orders.GET("", handler.GetCustomerOrders)
		orders.PATCH("/:id/cancel", handler.CancelOrder)
		orders.PATCH("/:id/status", handler.CompleteDelivery)
	}

	router.GET("/couriers/me/orders", handler.GetCourierOrders)
}
