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
		ConversationId: req.ConversationId,
		Input:          req.Input,
	})
	if err != nil {
		return err
	}

	// 流式处理
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	flusher, _ := w.(http.Flusher)
	var answer strings.Builder

	for {
		event, err := rpcStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			_ = writeChatSSE(w, types.WindAgentStreamEvent{
				MessageType: "error",
				Content:     err.Error(),
			})
			if flusher != nil {
				flusher.Flush()
			}
			l.Errorf("receive rag chat stream failed: %v", err)
			return nil
		}

		content := event.Content
		citations := ai_wind_tool.CitationsFromRPC(event.Citations)
		references := []ai_wind_tool.StreamReference(nil)
		if event.Type == "token" {
			answer.WriteString(content)
		}
		if event.Type == "done" {
			content = ""
			finalAnswer := answer.String()
			references = ai_wind_tool.StreamReferencesFromAnswer(finalAnswer, citations)
			if err := writeChatSSE(w, types.WindAgentStreamEvent{
				MessageType:    event.Type,
				Content:        content,
				ToolCalls:      event.ToolCalls,
				TraceId:        event.TraceId,
				ConversationId: event.ConversationId,
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
			MessageType: event.Type,
			Content:     content,
			TraceId:     event.TraceId,
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
