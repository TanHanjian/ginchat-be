package router

import (
	"ginchat/middlewares"
	"ginchat/service"
	"ginchat/service/friend_service"
	"ginchat/service/user_service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setCorsConfig(r *gin.Engine) *gin.Engine {
	corsConfig := cors.Config{
		AllowAllOrigins: true, // 允许所有来源
		// 或者可以指定特定的来源
		// AllowOrigins: []string{"http://example.com"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * 3600, // 预检请求的有效期
	}
	r.Use(cors.New(corsConfig))
	return r
}

func Router() *gin.Engine {
	r := gin.Default()
	r = setCorsConfig(r)
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
		friend.POST("/delete", friend_service.DeleteFriend)
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
