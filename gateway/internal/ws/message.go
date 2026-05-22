package ws

import (
	"encoding/json"
	"time"
)

// 消息类型常量
const (
	// 通知类消息
	MsgTypeNotification = "notification"
	// 数据更新消息
	MsgTypeDataUpdate = "data_update"
	// 系统消息
	MsgTypeSystem = "system"
	// 心跳消息
	MsgTypeHeartbeat = "heartbeat"
	// 错误消息
	MsgTypeError = "error"
)

// Message WebSocket 消息结构体
type Message struct {
	// Type 消息类型
	Type string `json:"type"`
	// Data 消息数据
	Data interface{} `json:"data"`
	// Timestamp 时间戳
	Timestamp int64 `json:"timestamp"`
	// UserID 目标用户ID（0=广播）
	UserID int64 `json:"userId,omitempty"`
}

// NewMessage 创建新消息
func NewMessage(msgType string, data interface{}, userID int64) *Message {
	return &Message{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now().Unix(),
		UserID:    userID,
	}
}

// ToJSON 消息序列化为 JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON 从 JSON 反序列化消息
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// NotificationData 通知消息数据结构
type NotificationData struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Level   string `json:"level,omitempty"` // info, warning, error
}

// DataUpdateData 数据更新消息结构
type DataUpdateData struct {
	ResourceType string      `json:"resourceType"` // 资源类型（如：user, role, order）
	ResourceID   int64       `json:"resourceId"`   // 资源ID
	Action       string      `json:"action"`       // 操作类型（如：create, update, delete）
	Data         interface{} `json:"data,omitempty"`
}

// SystemData 系统消息数据结构
type SystemData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
