package utils

import (
	"errors"
	user_models "ginchat/models/user_basic"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetUserIdFromToken(c *gin.Context) (uint, error) {
	userAny, exists := c.Get("user")
	if !exists {
		return 0, errors.New("user not found")
	}
	res, ok := userAny.(user_models.UserBasic)
	if !ok {
		return 0, errors.New("user data type error")
	}
	// 现在可以安全地使用 user 数据
	return uint(res.ID), nil
}

func BodyToModel[T any](c *gin.Context) (T, error) {
	var user_dto T
	if bind_error := c.ShouldBindBodyWith(&user_dto, binding.JSON); bind_error != nil {
		c.JSON(-1, gin.H{
			"message": bind_error.Error(),
		})
		return user_dto, bind_error
	}
	err := Go_validate.Struct(&user_dto)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return user_dto, err
	}
	return user_dto, nil
}
