package router

import (
	"ginchat/middlewares"
	"ginchat/service"
	"ginchat/service/user_service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", service.GetIndex)
	user := r.Group("/user")
	{
		user.GET("/list", user_service.GetUserList)
		user.POST("/create", user_service.CreateUser)
		user.POST("/delete", user_service.DeleteUserById)
		user.POST("/update", middlewares.AuthMiddleware(), user_service.UpdateUser)
		user.POST("/login/phone", user_service.LoginByUserPhone)
		user.POST("/login/email", user_service.LoginByUserEmail)
	}
	return r
}
