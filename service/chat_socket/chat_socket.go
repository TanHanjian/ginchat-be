package chatsocket

import (
	"ginchat/service/chatroom_service"
	"ginchat/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsMessage struct {
	Type   int         `json:"type"`
	Data   interface{} `json:"data"`
	UserId uint        `json:"userID"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
}

func StartChat(c *gin.Context) {
	userID, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	joinDto, err := utils.BodyToModel[chatroom_service.JoinChatroomDto](c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
			"retcode": -1,
		})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer conn.Close()
	sendCh := make(chan []byte)
	chatroomManager := chatroomManger.GetChatroomClientManager(joinDto.ChatroomID)
	newClient := Client{
		UserId: userID,
		Socket: conn,
		Send:   sendCh,
	}
	chatroomManager.Register <- &newClient
	go newClient.Read(chatroomManager)
	go newClient.HeatbeatCheck()
	newClient.Write(chatroomManager)
}
