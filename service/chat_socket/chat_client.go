package chatsocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserId   uint
	Socket   *websocket.Conn
	Send     chan []byte
	Heatbeat chan time.Time
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
		switch msg.Type {
		case msgHeatbeat:
			{
				c.Heatbeat <- time.Now()
			}
		default:
			{
				bytes, _ := json.Marshal(&WsMessage{Type: msg.Type, Data: msg.Data})
				c.Send <- bytes
			}
		}
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

func (c *Client) HeatbeatCheck() {
	ticker := time.NewTicker(7 * time.Second) // 每7秒发送一次心跳
	defer ticker.Stop()
	timer := time.NewTimer(2 * time.Second) // 创建一个定时器
	defer timer.Stop()
	for range ticker.C {
		pingReq := WsMessage{Type: msgHeatbeat, UserId: c.UserId}
		bytes, _ := json.Marshal(&pingReq)
		c.Send <- bytes
		// 重置定时器
		timer.Reset(5 * time.Second)
		select {
		case <-c.Heatbeat:
			{
				continue
			}
		case <-timer.C:
			{
				c.Socket.Close()
				break
			}
		}
	}
}
