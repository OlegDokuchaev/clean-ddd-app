package warehouse

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	items := router.Group("/items")
	{
		items.GET("", handler.GetAllItems)
		items.PATCH("/increase", handler.IncreaseQuantity)
		items.PATCH("/decrease", handler.DecreaseQuantity)
	}

	products := router.Group("/products")
	{
		products.POST("", handler.CreateProduct)
		products.PUT("/:id/image", handler.UpdateProductImage)
		products.GET("/:id/image", handler.GetProductImage)
	}
}
