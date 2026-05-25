package ws

import (
	"fmt"
	"net/http"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/ws"
	"go-zero-rpc/common/jwtx"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有来源（生产环境应该限制）
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsHandler WebSocket 升级处理
func WsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 从 query param 获取 token
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// 2. 验证 token
		claims, err := jwtx.ParseToken(token, svcCtx.Config.Auth.AccessSecret)
		if err != nil {
			logx.Errorf("WebSocket Token 验证失败: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 3. 检查 token 类型（只允许 access token）
		if claims.TokenType != jwtx.TokenTypeAccess {
			http.Error(w, "Invalid token type", http.StatusUnauthorized)
			return
		}

		// 4. 升级连接
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Errorf("WebSocket 升级失败: %v", err)
			return
		}

		// 5. 创建 Client 并注册到 Hub
		client := &ws.Client{
			ID:     fmt.Sprintf("%d", claims.UserId), // 使用 userId 作为 clientId
			UserID: claims.UserId,
			Conn:   conn,
			Send:   make(chan []byte, 256),
			Hub:    svcCtx.WsHub,
		}

		svcCtx.WsHub.Register <- client

		logx.Infof("WebSocket 连接已建立: userId=%d", claims.UserId)

		// 6. 启动读写 goroutine
		go client.WritePump()
		go client.ReadPump()
	}
}
