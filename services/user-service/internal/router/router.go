package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ppthinh/ChatApp/services/user-service/internal/service"
)

func SetupRouter(us service.UserService) *gin.Engine {
	r := gin.Default()
	r.POST("/user", createUser(us))
	r.GET("/user/:email", getUserByEmail(us))
	r.GET("/user/phone/:phone", getUserByPhoneNumber(us))
	r.PUT("/user", updateUser(us))
	r.DELETE("/user/:id", deleteUser(us))
	return r{

}