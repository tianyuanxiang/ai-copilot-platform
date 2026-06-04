// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_chat

import (
	aichatclient "ai-copilot-platform/ai-rpc/client/aichatservice"
	"ai-copilot-platform/gateway/internal/logic/ai_wind_tool"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"context"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiChatStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

var citationMarkerPattern = regexp.MustCompile(`\[(\d+)\]`)

func NewAiChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiChatStreamLogic {
	return &AiChatStreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiChatStreamLogic) AiChatStream(req *types.AiChatReq, w http.ResponseWriter) error {
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if strings.TrimSpace(req.Question) == "" {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "question 不能为空")
	}

	rpcStream, err := l.svcCtx.AiChatClient.RagChatStream(l.ctx, &aichatclient.RagChatReq{
		UserId:         userID,
		HasKbId:        req.KbId > 0,
		KbId:           req.KbId,
		ConversationId: req.ConversationId,
		Question:       req.Question,
		AnswerMode:     req.AnswerMode,
		SearchScope:    req.SearchScope,
		HasDomainId:    req.DomainId > 0,
		DomainId:       req.DomainId,
		DocumentIds:    req.DocumentIds,
	})
	if err != nil {
		return err
	}

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
			_ = ai_wind_tool.WriteChatSSE(w, ai_wind_tool.StreamEvent{
				Type:    "error",
				Content: err.Error(),
			})
			if flusher != nil {
				flusher.Flush()
			}
			l.Errorf("receive rag chat stream failed: %v", err)
			return nil
		}

		content := event.Content
		citations := ai_wind_tool.StreamCitationsFromRPC(event.Citations)
		references := []ai_wind_tool.StreamReference(nil)
		if event.Type == "token" {
			answer.WriteString(content)
		}
		if event.Type == "done" {
			content = ""
			finalAnswer := answer.String()
			references = ai_wind_tool.StreamReferencesFromAnswer(finalAnswer, citations)
			if err := ai_wind_tool.WriteChatSSE(w, ai_wind_tool.StreamEvent{
				Type:           event.Type,
				Content:        content,
				Answer:         finalAnswer,
				TraceID:        event.TraceId,
				ConversationID: event.ConversationId,
				Citations:      citations,
				References:     references,
			}); err != nil {
				return err
			}
			if flusher != nil {
				flusher.Flush()
			}
			continue
		}

		if err := ai_wind_tool.WriteChatSSE(w, ai_wind_tool.StreamEvent{
			Type:           event.Type,
			Content:        content,
			TraceID:        event.TraceId,
			ConversationID: event.ConversationId,
		}); err != nil {
			return err
		}
		if flusher != nil {
			flusher.Flush()
		}
	}
}
