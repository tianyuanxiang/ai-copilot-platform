package ws

import (
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	// 写超时时间
	writeWait = 10 * time.Second

	// Pong 等待超时时间
	pongWait = 60 * time.Second

	// Ping 发送间隔（必须小于 pongWait）
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 4096
)

// Client WebSocket 客户端连接
type Client struct {
	// ID 客户端唯一标识
	ID string
	// UserID 用户ID（已登录用户）
	UserID int64
	// Conn WebSocket 连接
	Conn *websocket.Conn
	// Send 发送缓冲区
	Send chan []byte
	// Hub 所属 Hub
	Hub *Hub
}

// ReadPump 读取客户端消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	// 设置读取限制
	c.Conn.SetReadLimit(maxMessageSize)
	// 设置 Pong 处理器
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// 读取消息
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logx.Errorf("WebSocket 读取错误: %v", err)
			}
			break
		}

		// 处理接收到的消息
		c.handleMessage(message)
	}
}

// WritePump 写入消息到客户端
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			// 设置写超时
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 写入消息
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logx.Errorf("WebSocket 写入错误: %v", err)
				return
			}

		case <-ticker.C:
			// 发送 Ping 消息
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logx.Errorf("WebSocket Ping 错误: %v", err)
				return
			}
		}
	}
}

// handleMessage 处理接收到的消息
func (c *Client) handleMessage(message []byte) {
	// 解析消息
	msg, err := FromJSON(message)
	if err != nil {
		logx.Errorf("消息解析错误: %v", err)
		return
	}

	// 根据消息类型处理
	switch msg.Type {
	case MsgTypeHeartbeat:
		// 心跳消息，回复 pong
		c.Send <- []byte(`{"type":"heartbeat","timestamp":` + strconv.FormatInt(time.Now().Unix(), 10) + `}`)
	default:
		// 其他消息可以在这里处理
		logx.Infof("收到消息: userId=%d, type=%s", c.UserID, msg.Type)
	}
}
