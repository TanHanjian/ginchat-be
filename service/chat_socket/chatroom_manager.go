package chatsocket

import (
	"errors"
	"sync"
)

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
