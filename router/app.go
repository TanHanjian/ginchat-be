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
	return r
}
