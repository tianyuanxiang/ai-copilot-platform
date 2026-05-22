package ws

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	// Redis Pub/Sub 频道名
	RedisChannelWSBroadcast = "ws:broadcast"
)

// RedisSubscriber Redis Pub/Sub 订阅器
type RedisSubscriber struct {
	// Hub WebSocket 连接管理器
	Hub *Hub
	// RDB Redis 客户端
	RDB *redis.Client
	// ctx 上下文
	ctx context.Context
	// cancel 取消函数
	cancel context.CancelFunc
}

// NewRedisSubscriber 创建新的 Redis 订阅器
func NewRedisSubscriber(hub *Hub, rdb *redis.Client) *RedisSubscriber {
	ctx, cancel := context.WithCancel(context.Background())
	return &RedisSubscriber{
		Hub:    hub,
		RDB:    rdb,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start 启动 Redis 订阅
func (rs *RedisSubscriber) Start() {
	go rs.subscribe()
	logx.Info("Redis Pub/Sub 订阅已启动")
}

// Stop 停止 Redis 订阅
func (rs *RedisSubscriber) Stop() {
	rs.cancel()
	logx.Info("Redis Pub/Sub 订阅已停止")
}

// subscribe 订阅 Redis 频道
func (rs *RedisSubscriber) subscribe() {
	// 订阅频道
	pubsub := rs.RDB.Subscribe(rs.ctx, RedisChannelWSBroadcast)
	defer pubsub.Close()

	// 等待订阅确认
	_, err := pubsub.Receive(rs.ctx)
	if err != nil {
		logx.Errorf("Redis 订阅失败: %v", err)
		return
	}

	// 接收消息
	ch := pubsub.Channel()
	for {
		select {
		case <-rs.ctx.Done():
			// 上下文取消，退出
			return
		case msg, ok := <-ch:
			if !ok {
				// 通道关闭，退出
				return
			}
			// 收到其他实例的广播消息，转发给本地客户端
			rs.Hub.BroadcastMessage([]byte(msg.Payload))
		}
	}
}

// PublishToRedis 发布消息到 Redis
func (rs *RedisSubscriber) PublishToRedis(message []byte) error {
	return rs.RDB.Publish(rs.ctx, RedisChannelWSBroadcast, message).Err()
}
