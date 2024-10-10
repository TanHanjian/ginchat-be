package chatroom

import (
	user_models "ginchat/models/user_basic"
	"ginchat/mydb"
	"time"
)

// ChatroomType 定义聊天室类型
type ChatroomType string

// 定义聊天室类型的常量
const (
	Single ChatroomType = "single"
	Group  ChatroomType = "group"
)

type Chatroom struct {
	ID        uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string       `gorm:"not null" json:"name"`
	Type      ChatroomType `json:"type"` // 聊天室类型
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

type ChatroomUsers struct {
	ChatroomID uint `gorm:"primaryKey"`
	UserID     uint `gorm:"primaryKey"`
}

type CreateChatroomData struct {
	Name string
	Type ChatroomType
}

type MultiAddUserRes struct {
	UserID uint
	Error  error
}

func AddUserToChatroom(userID, chatroomID uint) error {
	chatroomUser := ChatroomUsers{
		ChatroomID: chatroomID,
		UserID:     userID,
	}
	err := mydb.DB.Table("chatroom_users").Create(&chatroomUser).Error
	return err
}

func RemoveUserFromChatroom(userID, chatroomID uint) (ChatroomUsers, error) {
	var rel ChatroomUsers
	err := mydb.DB.Table("chatroom_users").Where("chatroom_users.chatroom_id = ? AND chatroom_users.user_id = ?", chatroomID, userID).Delete(&rel).Error
	return rel, err
}

func MultiAddUserToChatRoom(userIDs []uint, chatroomID uint) []MultiAddUserRes {
	resList := make([]MultiAddUserRes, 0, len(userIDs))
	ch := make(chan MultiAddUserRes)
	for _, friendID := range userIDs {
		id := friendID
		go func() {
			err := AddUserToChatroom(id, chatroomID)
			res := MultiAddUserRes{
				Error:  err,
				UserID: id,
			}
			ch <- res
		}()
	}
	for res := range ch {
		resList = append(resList, res)
	}
	return resList
}

func CreateChatroom(data CreateChatroomData) (Chatroom, error) {
	chatroom := Chatroom{
		Name: data.Name,
		Type: data.Type,
	}
	err := mydb.DB.Table("chatroom").Create(&chatroom).Error
	return chatroom, err
}

func GetAllChatsByUserId(user_id uint) ([]Chatroom, error) {
	var chatrooms []Chatroom

	// 使用 JOIN 查询获取与用户关联的聊天室
	err := mydb.DB.Table("chatroom").
		Select("chatroom.*").
		Joins("JOIN chatroom_users ON chatroom_users.chatroom_id = chatroom.id").
		Where("chatroom_users.user_id = ?", user_id).
		Scan(&chatrooms).Error

	return chatrooms, err
}

func GetAllUsersByChatroomId(chatroom_id uint) ([]user_models.UserBasic, error) {
	var users []user_models.UserBasic
	err := mydb.DB.Table("user_basic").
		Select("user_basic.*").
		Joins("JOIN chatroom_users ON chatroom_users.user_id = user_baisc.id").
		Where("chatroom_users.chatroom_id = ?", chatroom_id).
		Scan(&users).Error
	return users, err
}
