package courier

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	couriers := router.Group("/couriers")
	{
		couriers.POST("/register", handler.Register)
		couriers.POST("/login", handler.Login)
	}
}
