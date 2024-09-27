package user_service

import (
	"errors"
	"fmt"
	user_models "ginchat/models/user_basic"
	"ginchat/utils"
	"reflect"

	"github.com/gin-gonic/gin"
)

func GetUserList(c *gin.Context) {
	users := user_models.GetUserList()
	c.JSON(200, gin.H{
		"message": "success",
		"data":    users,
	})
}

func checkNoRepeatUser(new_user *user_models.UserBasic) error {
	repeat_user, error := user_models.CheckRepeat(new_user)
	if error == nil {
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

func CreateUser(c *gin.Context) {
	user_dto, err := utils.BodyToModel[UserCreateDto](c)
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
	err = user_models.Create(new_user)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "succeeded",
		})
	}
}

func DeleteUserById(c *gin.Context) {
	user_dto, err := utils.BodyToModel[UserDeleteDto](c)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
	}
	err = user_models.DeleteByUserID(user_dto.User_id)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
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

func UpdateUser(c *gin.Context) {
	user_dto, err := utils.BodyToModel[UserUpdateDto](c)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
	}
	user_id, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return
	}
	update_user := user_models.UserBasic{
		Model: user_models.Model{
			ID: user_id,
		},
	}
	if user_dto.Name != "" {
		update_user.Name = user_dto.Name
	}
	if user_dto.Password != "" {
		update_user.Password = utils.Md5(user_dto.Password)
	}
	err = user_models.Update(update_user)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "succeeded",
		})
	}
}

type FindUserFn func(user *user_models.UserBasic) (user_models.UserBasic, error)

func login[T any](c *gin.Context, fn FindUserFn) {
	user_dto, err := utils.BodyToModel[T](c)
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
	exist_user, err := fn(&user)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
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
