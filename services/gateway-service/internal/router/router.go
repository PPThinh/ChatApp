package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ppthinh/ChatApp/services/gateway-service/internal/handler"
)

func StartGateway(port string, h *handler.GatewayHandler) error {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
		api.GET("/users/email", h.GetUserByEmail)
		//TODO: add get user by phone number and get all users with page limit
	}

	return router.Run(port)
}
