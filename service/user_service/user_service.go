package user_service

import (
	"errors"
	"fmt"
	user_models "ginchat/models/user_basic"
	"ginchat/utils"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

func GetUserList(c *gin.Context) {
	users := user_models.GetUserList()
	c.JSON(200, gin.H{
		"message": "success",
		"data":    users,
	})
}

func checkOnlyUserByName(new_user *user_models.UserBasic) error {
	if data := user_models.FindByName(new_user).Error; data == nil {
		return errors.New("user name is existed")
	}
	return nil
}

func checkOnlyUserByPhone(new_user *user_models.UserBasic) error {
	if data := user_models.FindByPhone(new_user).Error; data == nil {
		return errors.New("user phone is existed")
	}
	return nil
}

func checkOnlyUserByEmail(new_user *user_models.UserBasic) error {
	if data := user_models.FindByName(new_user).Error; data == nil {
		return errors.New("user email is existed")
	}
	return nil
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
	err := utils.Go_validate.Struct(&user_dto)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return
	}
	new_user := user_models.UserBasic{
		Name:     user_dto.Name,
		Password: user_dto.Password,
		Email:    user_dto.Email,
		Phone:    user_dto.Phone,
	}
	if err := checkOnlyUserByName(&new_user); err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := checkOnlyUserByPhone(&new_user); err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := checkOnlyUserByEmail(&new_user); err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return
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

func DeleteUserById(c *gin.Context) {
	var user_dto UserDeleteDto
	if bind_error := c.ShouldBindBodyWith(&user_dto, binding.JSON); bind_error != nil {
		c.JSON(-1, gin.H{
			"message": bind_error.Error(),
		})
		return
	}
	err := utils.Go_validate.Struct(&user_dto)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return
	}
	result := user_models.DeleteByUserID(user_dto.User_id)
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

func setUpdateUser(user_dto *UserUpdateDto, user *user_models.UserBasic) error {
	// 确保 user_dto 是指向结构体的指针
	if reflect.ValueOf(user_dto).Kind() != reflect.Ptr || reflect.ValueOf(user).Kind() != reflect.Ptr {
		return fmt.Errorf("both target and source must be pointers to structs")
	}

	targetValue := reflect.ValueOf(user).Elem() // 解引用
	sourceValue := reflect.ValueOf(user_dto).Elem()

	a := targetValue.Kind()
	b := sourceValue.Kind()
	fmt.Println(a, b)
	// 确保传入的是结构体
	if targetValue.Kind() != reflect.Struct || sourceValue.Kind() != reflect.Struct {
		return fmt.Errorf("both target and source must be structs")
	}

	for i := 0; i < sourceValue.NumField(); i++ {
		field := sourceValue.Type().Field(i)
		targetField := targetValue.FieldByName(field.Name)

		// 检查目标字段是否有效且可设置
		if targetField.IsValid() && targetField.CanSet() {
			// 检查类型是否匹配
			if targetField.Type() == field.Type {
				targetField.Set(sourceValue.Field(i)) // 设置目标字段值
			} else {
				return fmt.Errorf("type mismatch for field %s: expected %s, got %s", field.Name, targetField.Type(), field.Type)
			}
		}
	}
	return nil
}

func UpdateUser(c *gin.Context) {
	var user_dto UserUpdateDto
	if bind_error := c.ShouldBindBodyWith(&user_dto, binding.JSON); bind_error != nil {
		c.JSON(-1, gin.H{
			"message": bind_error.Error(),
		})
		return
	}
	err := utils.Go_validate.Struct(&user_dto)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return
	}
	update_user := user_models.UserBasic{
		Model: gorm.Model{
			ID: uint(user_dto.User_id),
		},
	}
	if err := setUpdateUser(&user_dto, &update_user); err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return
	}
	result := user_models.Update(update_user)
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
