package customer

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router gin.IRouter, handler *Handler) {
	customers := router.Group("/customers")
	{
		customers.POST("/register", handler.Register)
		customers.POST("/login", handler.Login)
	}
}
