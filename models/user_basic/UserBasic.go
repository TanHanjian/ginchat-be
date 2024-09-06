package user_models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string
	Email         string
	ClientIp      string
	Identity      string
	ClientPort    string
	LoginTime     uint64
	HeartbeatTime uint64
	LogoutTime    uint64
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	users := make([]*UserBasic, 10)
	utils.DB.Find(&users)
	for _, user := range users {
		fmt.Println(user)
	}
	return users
}

func Create(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}
