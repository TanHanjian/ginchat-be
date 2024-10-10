package chatsocket

import (
	"encoding/json"
	"errors"
	"ginchat/service/chatroom_service"
	"ginchat/utils"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 定义枚举
type Status int

type Client struct {
	UserId uint
	Socket *websocket.Conn
	Send   chan []byte
}

func (c *Client) Read(manager *ClientManager) {
	defer func() {
		_ = c.Socket.Close()
		manager.UnRegister <- c
	}()
	for {
		_, data, err := c.Socket.ReadMessage()
		if err != nil {
			break
		}
		var msg WsMessage
		err = json.Unmarshal(data, &msg)
		if err != nil {
			break
		}
		bytes, _ := json.Marshal(&WsMessage{Type: 3, Data: msg.Data})
		c.Send <- bytes
	}
}

func (c *Client) Write(manager *ClientManager) {
	defer func() {
		_ = c.Socket.Close()
		manager.UnRegister <- c
	}()
	for msg := range c.Send {
		err := c.Socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

type ClientManager struct {
	Clients    map[uint]*Client // 记录在线用户
	Broadcast  chan []byte      //触发消息广播
	Register   chan *Client     // 触发新用户登陆
	UnRegister chan *Client     // 触发用户退出
	ChatroomID uint
}

type ChatroomManager struct {
	Chatrooms map[uint]*ClientManager
	mu        sync.RWMutex // 使用读写锁
}

var chatroomManger = ChatroomManager{}

func (m *ChatroomManager) GetChatroomClientManager(chatroomID uint) *ClientManager {
	m.mu.RLock()
	chatroomClient, exists := m.Chatrooms[chatroomID]
	if exists && chatroomClient != nil {
		m.mu.RUnlock()
		return chatroomClient
	}
	m.mu.RUnlock()
	newClient := ClientManager{
		ChatroomID: chatroomID,
	}
	m.mu.Lock()
	m.Chatrooms[chatroomID] = &newClient
	m.mu.Unlock()
	go newClient.Start()
	go newClient.Quit()
	go newClient.BroadcastSend()
	return &newClient
}

func (m *ChatroomManager) RemoveChatroomClientManager(chatroomID uint) error {
	m.mu.RLock()
	chatroomClient, exists := m.Chatrooms[chatroomID]
	if !exists || chatroomClient == nil {
		m.mu.RUnlock()
		return errors.New("chatroom not exist")
	}
	m.mu.RUnlock()
	for _, client := range chatroomClient.Clients {
		client.Socket.Close()
	}
	m.mu.Lock()
	delete(m.Chatrooms, chatroomID)
	m.mu.Unlock()
	close(chatroomClient.Broadcast)
	close(chatroomClient.Register)
	close(chatroomClient.UnRegister)
	return nil
}

func (manager *ClientManager) Start() {
	for client := range manager.Register {
		userID := client.UserId
		manager.Clients[userID] = client
		// 如果有新用户连接则发送最近聊天记录和在线人数给他
		count := len(manager.Clients)
		manager.InitSend(client, count, userID)
	}
}

func (manager *ClientManager) Quit() {
	for client := range manager.UnRegister {
		userID := client.UserId
		_, exists := manager.Clients[userID]
		if !exists {
			continue
		}
		delete(manager.Clients, userID)
		// 给客户端刷新在线人数
		if len(manager.Clients) > 0 {
			resp, _ := json.Marshal(&WsMessage{Type: 1, Data: len(manager.Clients), UserId: userID})
			manager.Broadcast <- resp
		} else {
			chatroomManger.RemoveChatroomClientManager(manager.ChatroomID)
		}
	}
}

func (manager *ClientManager) InitSend(cur *Client, count int, userId uint) {
	resp, _ := json.Marshal(&WsMessage{Type: 1, Data: count, UserId: userId})
	manager.Broadcast <- resp

	bytes, _ := json.Marshal(&WsMessage{Type: 2, Data: "123", UserId: userId})
	cur.Send <- bytes
}

func (manager *ClientManager) BroadcastSend() {
	for msg := range manager.Broadcast {
		for _, client := range manager.Clients {
			client.Send <- msg
		}
	}
}

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
	newClient.Write(chatroomManager)
}
