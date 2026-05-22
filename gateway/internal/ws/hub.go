package ws

import (
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

// Hub WebSocket 连接管理器
type Hub struct {
	// Clients 已注册的客户端（UserID -> Client 集合）
	Clients map[int64]map[*Client]bool
	// Broadcast 广播消息通道
	Broadcast chan []byte
	// Register 注册客户端通道
	Register chan *Client
	// Unregister 注销客户端通道
	Unregister chan *Client
	// mutex 读写锁（保护 Clients map）
	mutex sync.RWMutex
}

// NewHub 创建新的 Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int64]map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run 启动 Hub 主循环
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// 注册客户端
			h.registerClient(client)

		case client := <-h.Unregister:
			// 注销客户端
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			// 广播消息
			h.broadcastMessage(message)
		}
	}
}

// registerClient 注册客户端
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// 获取用户的客户端集合
	clients, exists := h.Clients[client.UserID]
	if !exists {
		clients = make(map[*Client]bool)
		h.Clients[client.UserID] = clients
	}

	// 添加客户端
	clients[client] = true

	logx.Infof("WebSocket 客户端已注册: userId=%d, clientId=%s, 总连接数=%d",
		client.UserID, client.ID, h.GetConnectionCount())
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// 获取用户的客户端集合
	clients, exists := h.Clients[client.UserID]
	if !exists {
		return
	}

	// 删除客户端
	if _, ok := clients[client]; ok {
		delete(clients, client)
		close(client.Send)

		// 如果用户没有其他连接，删除用户条目
		if len(clients) == 0 {
			delete(h.Clients, client.UserID)
		}

		logx.Infof("WebSocket 客户端已注销: userId=%d, clientId=%s, 总连接数=%d",
			client.UserID, client.ID, h.GetConnectionCount())
	}
}

// broadcastMessage 广播消息到所有客户端
func (h *Hub) broadcastMessage(message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for _, clients := range h.Clients {
		for client := range clients {
			select {
			case client.Send <- message:
				// 消息已发送
			default:
				// 发送缓冲区已满，关闭连接
				go func(c *Client) {
					h.Unregister <- c
				}(client)
			}
		}
	}
}

// SendToUser 发送消息给指定用户
func (h *Hub) SendToUser(userID int64, message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	clients, exists := h.Clients[userID]
	if !exists {
		return
	}

	for client := range clients {
		select {
		case client.Send <- message:
			// 消息已发送
		default:
			// 发送缓冲区已满，关闭连接
			go func(c *Client) {
				h.Unregister <- c
			}(client)
		}
	}
}

// BroadcastMessage 广播消息（公开方法）
func (h *Hub) BroadcastMessage(message []byte) {
	h.Broadcast <- message
}

// GetOnlineUsers 获取在线用户列表
func (h *Hub) GetOnlineUsers() []int64 {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	users := make([]int64, 0, len(h.Clients))
	for userID := range h.Clients {
		users = append(users, userID)
	}
	return users
}

// GetConnectionCount 获取总连接数
func (h *Hub) GetConnectionCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	count := 0
	for _, clients := range h.Clients {
		count += len(clients)
	}
	return count
}

// GetUserConnectionCount 获取指定用户的连接数
func (h *Hub) GetUserConnectionCount(userID int64) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	clients, exists := h.Clients[userID]
	if !exists {
		return 0
	}
	return len(clients)
}

// IsUserOnline 检查用户是否在线
func (h *Hub) IsUserOnline(userID int64) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	clients, exists := h.Clients[userID]
	return exists && len(clients) > 0
}
