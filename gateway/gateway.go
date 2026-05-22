// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"ai-copilot-platform/gateway/internal/config"
	"ai-copilot-platform/gateway/internal/handler"
	"ai-copilot-platform/gateway/internal/svc"
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"go-zero-rpc/common/response"
	"go-zero-rpc/common/rpcerr"
	"go-zero-rpc/gateway/pkg/sysgateway"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	logx.DisableStat()
	// 将慢查询阈值设为 1ms，使所有 SQL 语句都被记录到日志中（仅开发调试用）
	sqlx.SetSlowThreshold(time.Millisecond)

	// 全局错误处理：handler 调用 httpx.ErrorCtx(ctx, w, err) 时回调本函数。
	// 把 RPC 透传过来的 status 错误统一解析成业务码 + 消息，
	// 始终以 HTTP 200 + 统一 Response 结构返回，与 response.OkWithData 保持一致。
	httpx.SetErrorHandlerCtx(func(_ context.Context, err error) (int, any) {
		code, msg := rpcerr.FromStatus(err)
		return http.StatusOK, response.Response{
			Code: code,
			Msg:  msg,
		}
	})

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithFileServer("/uploads", http.Dir(c.Upload.Path)))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	// 再注册架构仓
	sysgateway.Register(server, sysgateway.Config{
		Auth: sysgateway.AuthConfig{
			AccessSecret: c.Auth.AccessSecret,
			AccessExpire: c.Auth.AccessExpire,
		},
		SysRpc: c.SysRpc,
		Upload: sysgateway.UploadConfig{
			Path:        c.Upload.Path,
			MaxSize:     c.Upload.MaxSize,
			AllowedExts: c.Upload.AllowedExts,
		},
		WebSocket: sysgateway.WebSocketConfig{
			HeartbeatInterval: c.WebSocket.HeartbeatInterval,
			MaxConnections:    c.WebSocket.MaxConnections,
		},
	})
	handler.RegisterHandlers(server, ctx)

	server.PrintRoutes()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
