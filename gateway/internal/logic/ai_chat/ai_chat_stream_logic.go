// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_chat

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	aichatclient "ai-copilot-platform/ai-rpc/client/aichatservice"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiChatStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

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

	for {
		event, err := rpcStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			_ = writeChatSSE(w, streamEvent{
				Type:    "error",
				Content: err.Error(),
			})
			if flusher != nil {
				flusher.Flush()
			}
			l.Errorf("receive rag chat stream failed: %v", err)
			return nil
		}

		if err := writeChatSSE(w, streamEvent{
			Type:      event.Type,
			Content:   event.Content,
			TraceID:   event.TraceId,
			Citations: citationsFromRPC(event.Citations),
		}); err != nil {
			return err
		}
		if flusher != nil {
			flusher.Flush()
		}
	}
}

type streamEvent struct {
	Type      string             `json:"type"`
	Content   string             `json:"content,omitempty"`
	TraceID   string             `json:"traceId,omitempty"`
	Citations []types.AiCitation `json:"citations,omitempty"`
}

func writeChatSSE(w http.ResponseWriter, event streamEvent) error {
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
