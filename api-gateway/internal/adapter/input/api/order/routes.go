package order

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router gin.IRouter, handler *Handler) {
	orders := router.Group("/orders")
	{
		orders.POST("", handler.Create)
		orders.GET("", handler.GetCustomerOrders)
		orders.PATCH("/:id/cancel", handler.CancelOrder)
		orders.PATCH("/:id/complete", handler.CompleteDelivery)
	}

	router.GET("/couriers/me/orders", handler.GetCourierOrders)
}
