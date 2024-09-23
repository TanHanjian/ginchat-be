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

func checkNoRepeatUser(new_user *user_models.UserBasic) error {
	repeat_res, repeat_user := user_models.CheckRepeat(new_user)
	if repeat_res.Error == nil {
		if repeat_user.Name == new_user.Name {
			return errors.New("user name is existed")
		}
		if repeat_user.Email == new_user.Email {
			return errors.New("user email is existed")
		}
		if repeat_user.Phone == new_user.Phone {
			return errors.New("user phone is existed")
		}
	}
	return nil
}

func bodyToModel[T any](c *gin.Context) (error, T) {
	var user_dto T
	if bind_error := c.ShouldBindBodyWith(&user_dto, binding.JSON); bind_error != nil {
		c.JSON(-1, gin.H{
			"message": bind_error.Error(),
		})
		return bind_error, user_dto
	}
	err := utils.Go_validate.Struct(&user_dto)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return err, user_dto
	}
	return nil, user_dto
}

func CreateUser(c *gin.Context) {
	err, user_dto := bodyToModel[UserCreateDto](c)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
	}
	if user_dto.Re_password != user_dto.Password {
		c.JSON(-1, gin.H{
			"message": "password is not equal to re_password!",
		})
		return
	}
	new_user := user_models.UserBasic{
		Name:     user_dto.Name,
		Password: utils.Md5(user_dto.Password),
		Email:    user_dto.Email,
		Phone:    user_dto.Phone,
	}
	if err := checkNoRepeatUser(&new_user); err != nil {
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
	err, user_dto := bodyToModel[UserDeleteDto](c)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
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

func setUpdateUser[T any](user_dto *T, user *user_models.UserBasic) error {
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

func getUserIdFromToken(c *gin.Context) (uint, error) {
	userAny, exists := c.Get("user")
	if !exists {
		return 0, errors.New("user not found")
	}
	user, ok := userAny.(map[string]interface{})
	if !ok {
		return 0, errors.New("user data type error")
	}
	// 现在可以安全地使用 user 数据
	userID, ok := user["id"].(float64) // 假设 ID 是数字类型
	if !ok {
		return 0, errors.New("no user id")
	}
	return uint(userID), nil
}

func UpdateUser(c *gin.Context) {
	err, user_dto := bodyToModel[UserUpdateDto](c)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
	}
	user_id, err := getUserIdFromToken(c)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return
	}
	update_user := user_models.UserBasic{
		Model: gorm.Model{
			ID: user_id,
		},
	}
	if err := setUpdateUser[UserUpdateDto](&user_dto, &update_user); err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return
	}
	user_dto.Password = utils.Md5(user_dto.Password)
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

type FindUserFn func(user *user_models.UserBasic) (*gorm.DB, user_models.UserBasic)

func login[T any](c *gin.Context, fn FindUserFn) {
	err, user_dto := bodyToModel[T](c)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
	}
	var user user_models.UserBasic
	if err := setUpdateUser[T](&user_dto, &user); err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
	}
	res, exist_user := fn(&user)
	if res.Error != nil {
		c.JSON(-1, gin.H{
			"message": res.Error.Error(),
		})
	}
	if exist_user.Password != utils.Md5(user.Password) {
		c.JSON(-1, gin.H{
			"message": "password is wrong!",
		})
		return
	}
	exist_user.Password = ""
	token, err := utils.GenerateJWT(exist_user)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "succeeded",
		"token":   token,
	})
}

func LoginByUserPhone(c *gin.Context) {
	login[LoginByUserPhoneDto](c, user_models.FindByPhone)
}

func LoginByUserEmail(c *gin.Context) {
	login[LoginByUserEmailDto](c, user_models.FindByEmail)
}
