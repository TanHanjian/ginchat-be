package chatroom_service

import (
	chatroom_models "ginchat/models/chatroom"
	"ginchat/utils"

	"github.com/gin-gonic/gin"
)

func CreateSingleChatroom(c *gin.Context) {
	userID, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	single_dto, err := utils.BodyToModel[CreateSingleRoomDto](c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	chatroom, err := chatroom_models.CreateChatroom(chatroom_models.CreateChatroomData{
		Name: single_dto.Name,
		Type: chatroom_models.Single,
	})
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	err = chatroom_models.AddUserToChatroom(userID, chatroom.ID)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	err = chatroom_models.AddUserToChatroom(single_dto.FriendID, chatroom.ID)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"retcode": 1,
	})
}

func CreateMultiChatroom(c *gin.Context) {
	userID, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	multiDto, err := utils.BodyToModel[CreateMultiRoomDto](c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	chatroom, err := chatroom_models.CreateChatroom(chatroom_models.CreateChatroomData{
		Name: multiDto.Name,
		Type: chatroom_models.Group,
	})
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	err = chatroom_models.AddUserToChatroom(userID, chatroom.ID)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	resList := chatroom_models.MultiAddUserToChatRoom(multiDto.FriendIDs, chatroom.ID)
	c.JSON(200, gin.H{
		"message": "success",
		"data":    resList,
		"retcode": 1,
	})
}

func JoinChatroom(c *gin.Context) {
	userID, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	joinDto, err := utils.BodyToModel[JoinChatroomDto](c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	err = chatroom_models.AddUserToChatroom(userID, joinDto.ChatroomID)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"retcode": 1,
	})
}

func QuitChatroom(c *gin.Context) {
	userID, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	joinDto, err := utils.BodyToModel[QuitChatroomDto](c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	_, err = chatroom_models.RemoveUserFromChatroom(userID, joinDto.ChatroomID)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"retcode": 1,
	})
}
