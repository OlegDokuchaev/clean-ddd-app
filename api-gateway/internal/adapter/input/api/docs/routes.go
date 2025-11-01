package docs

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(router gin.IRouter) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
