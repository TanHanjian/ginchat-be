package chatsocket

import "encoding/json"

type ClientManager struct {
	Clients    map[uint]*Client // 记录在线用户
	Broadcast  chan []byte      //触发消息广播
	Register   chan *Client     // 触发新用户登陆
	UnRegister chan *Client     // 触发用户退出
	ChatroomID uint
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
