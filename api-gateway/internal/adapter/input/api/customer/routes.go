package customer

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router gin.IRouter, handler *Handler) {
	customers := router.Group("/customers")
	{
		customers.POST("/register", handler.Register)
		customers.POST("/login", handler.Login)

		challenges := customers.Group("/auth-challenges")
		{
			challenges.PATCH("/:challenge_id", handler.VerifyOtp)
		}

		passwordResets := customers.Group("/password-resets")
		{
			passwordResets.POST("", handler.RequestPasswordReset)
			passwordResets.PATCH("/:token", handler.CompletePasswordReset)
		}
	}
}
