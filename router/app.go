package router

import (
	"ginchat/middlewares"
	"ginchat/service"
	"ginchat/service/friend_service"
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
	friend := r.Group("/friend")
	friend.Use(middlewares.AuthMiddleware())
	{
		friend.POST("/list", friend_service.GetFriendList)
	}
	friend_apply := friend.Group("/apply")
	{
		friend_apply.POST("/create", friend_service.CreateFriendApply)
		friend_apply.POST("/agree", friend_service.AgreeFriendApply)
		friend_apply.POST("/reject", friend_service.RejectFriendApply)
		friend_apply.POST("/to_list", friend_service.GetFriendApplyToList)
		friend_apply.POST("/from_list", friend_service.GetFriendApplyFromList)
	}
	return r
}
