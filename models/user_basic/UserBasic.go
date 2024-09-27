package user_models

import (
	"fmt"
	"ginchat/mydb"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string `gorm:"unique;not null"`
	Password      string `gorm:"not null"`
	Phone         string `gorm:"unique;not null"`
	Email         string `gorm:"unique;not null"`
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
	mydb.DB.Find(&users)
	for _, user := range users {
		fmt.Println(user)
	}
	return users
}

func Create(user UserBasic) error {
	return mydb.DB.Create(&user).Error
}

func DeleteByUserID(user_id int) error {
	var user UserBasic
	user.ID = uint(user_id)
	return mydb.DB.Delete(&user).Error
}

func Update(user UserBasic) error {
	var existingUser UserBasic
	result := mydb.DB.First(&existingUser, user.ID)
	if result.Error != nil {
		// 如果找不到记录，返回错误
		return result.Error
	}

	// 更新记录
	result = mydb.DB.Model(&existingUser).Updates(user)
	return result.Error
}

func FindByPhone(user *UserBasic) (UserBasic, error) {
	var exist_user UserBasic
	return exist_user, mydb.DB.Model(user).Where("phone = ?", user.Phone).First(&exist_user).Error
}

func FindByEmail(user *UserBasic) (UserBasic, error) {
	var exist_user UserBasic
	return exist_user, mydb.DB.Model(user).Where("email = ?", user.Email).First(&exist_user).Error
}

func CheckRepeat(user *UserBasic) (UserBasic, error) {
	var exist_user UserBasic
	res := mydb.DB.Model(&UserBasic{}).Where("name = ? OR email = ? OR phone = ?", user.Name, user.Email, user.Phone).First(&exist_user)
	return exist_user, res.Error
}

func FindByID(id uint) (UserBasic, error) {
	var exist_user UserBasic
	return exist_user, mydb.DB.Model(&UserBasic{}).First(&exist_user, id).Error
}
