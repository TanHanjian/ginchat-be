package router

import (
	"ginchat/service"
	"ginchat/service/user_service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", service.GetIndex)
	r.GET("/user/list", user_service.GetUserList)
	r.POST("/user/create", user_service.CreateUser)
	r.POST("/user/delete", user_service.DeleteUserById)
	r.POST("/user/update", user_service.UpdateUser)
	r.POST("/user/login/phone", user_service.LoginByUserPhone)
	r.POST("/user/login/email", user_service.LoginByUserEmail)
	return r
}
