package user_service

import (
	user_models "ginchat/models/user_basic"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetUserList(c *gin.Context) {
	users := user_models.GetUserList()
	c.JSON(200, gin.H{
		"message": "success",
		"data":    users,
	})
}

func CreateUser(c *gin.Context) {
	var user_dto UserCreateDto
	if bind_error := c.ShouldBindBodyWith(&user_dto, binding.JSON); bind_error != nil {
		c.JSON(-1, gin.H{
			"message": bind_error.Error(),
		})
		return
	}
	if user_dto.Re_password != user_dto.Password {
		c.JSON(-1, gin.H{
			"message": "password is not equal to re_password!",
		})
		return
	}
	new_user := user_models.UserBasic{
		Name:     user_dto.Name,
		Password: user_dto.Password,
	}
	result := user_models.Create(new_user)
	if result.Error != nil {
		c.JSON(-1, gin.H{
			"message": result.Error.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "succeeded",
		})
	}
}
