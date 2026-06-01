// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_wind_agent

import (
	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/logic/ai_wind_tool"
	"context"
	"encoding/json"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"
	"io"
	"net/http"
	"strings"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiWindCopilotAgentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindCopilotAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindCopilotAgentLogic {
	return &AiWindCopilotAgentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindCopilotAgentLogic) AiWindCopilotAgent(req *types.AiWindAgentRunReq, w http.ResponseWriter) error {
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	if strings.TrimSpace(req.Input) == "" {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "输入不能为空")
	}

	rpcStream, err := l.svcCtx.AiWindAgentClient.RunAgentStream(l.ctx, &pb.WindAgentRunReq{
		UserId:         userID,
		AgentSessionId: req.AgentSessionId,
		Input:          req.Input,
	})
	if err != nil {
		return err
	}

	// 流式处理
	return forwardWindAgentStream(w, rpcStream, req.AgentSessionId)
}

type windAgentStreamReceiver interface {
	Recv() (*pb.WindAgentStreamEvent, error)
}

func forwardWindAgentStream(w http.ResponseWriter, rpcStream windAgentStreamReceiver, fallbackAgentSessionID string) error {
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	flusher, _ := w.(http.Flusher)

	for {
		event, err := rpcStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			_ = writeChatSSE(w, types.WindAgentStreamEvent{
				MessageType:    "error",
				AgentSessionId: fallbackAgentSessionID,
				Content:        err.Error(),
			})
			if flusher != nil {
				flusher.Flush()
			}
			return nil
		}

		agentSessionID := event.AgentSessionId
		if agentSessionID == "" {
			agentSessionID = fallbackAgentSessionID
		}
		content := event.Content
		citations := ai_wind_tool.CitationsFromRPC(event.Citations)
		if event.Type == "done" {
			if err := writeChatSSE(w, types.WindAgentStreamEvent{
				MessageType:    event.Type,
				Content:        content,
				ToolCalls:      ai_wind_tool.ToolCallsFromRPC(event.ToolCalls),
				TraceId:        event.TraceId,
				AgentSessionId: agentSessionID,
				Citations:      citations,
			}); err != nil {
				return err
			}
			if flusher != nil {
				flusher.Flush()
			}
			continue
		}

		if err := writeChatSSE(w, types.WindAgentStreamEvent{
			MessageType:    event.Type,
			Content:        content,
			TraceId:        event.TraceId,
			AgentSessionId: agentSessionID,
		}); err != nil {
			return err
		}
		if flusher != nil {
			flusher.Flush()
		}
	}
}

func writeChatSSE(w http.ResponseWriter, event types.WindAgentStreamEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte("data: ")); err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	_, err = w.Write([]byte("\n\n"))
	return err
}
