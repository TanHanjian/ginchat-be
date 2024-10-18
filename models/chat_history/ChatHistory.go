package chathisotry_models

import (
	"ginchat/mydb"
	"time"
)

type ChatHistory struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Content    string    `gorm:"type:varchar(500)" json:"content"`
	Type       int       `json:"type"` // 聊天记录类型 text:0、image:1、file:2
	ChatroomID int       `json:"chatroomId"`
	SenderID   int       `json:"senderId"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

type GetChatHistoryParams struct {
	ChatroomID uint
	PageSize   uint
	PageNo     uint
}

func GetChatHistory(data GetChatHistoryParams) ([]ChatHistory, error) {
	var list []ChatHistory
	err := mydb.DB.Table("chat_history").Where("ChatroomId = ?", data.ChatroomID).Offset(int(data.PageNo) - 1).Limit(int(data.PageSize)).Find(&list).Error
	return list, err
}

func AddChatHistory(data ChatHistory) error {
	err := mydb.DB.Table("chat_history").Save(data).Error
	return err
}
