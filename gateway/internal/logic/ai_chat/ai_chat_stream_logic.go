// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_chat

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	aichatclient "ai-copilot-platform/ai-rpc/client/aichatservice"
	"ai-copilot-platform/ai-rpc/pb"
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

		content := event.Content
		citations := streamCitationsFromRPC(event.Citations)
		references := []streamReference(nil)
		if event.Type == "token" {
			answer.WriteString(content)
		}
		if event.Type == "done" {
			content = ""
			finalAnswer := answer.String()
			references = streamReferencesFromAnswer(finalAnswer, citations)
			if err := writeChatSSE(w, streamEvent{
				Type:       event.Type,
				Content:    content,
				Answer:     finalAnswer,
				TraceID:    event.TraceId,
				Citations:  citations,
				References: references,
			}); err != nil {
				return err
			}
			if flusher != nil {
				flusher.Flush()
			}
			continue
		}

		if err := writeChatSSE(w, streamEvent{
			Type:    event.Type,
			Content: content,
			TraceID: event.TraceId,
		}); err != nil {
			return err
		}
		if flusher != nil {
			flusher.Flush()
		}
	}
}

type streamEvent struct {
	Type       string            `json:"type"`
	Content    string            `json:"content,omitempty"`
	Answer     string            `json:"answer,omitempty"`
	TraceID    string            `json:"traceId,omitempty"`
	Citations  []streamCitation  `json:"citations,omitempty"`
	References []streamReference `json:"references,omitempty"`
}

type streamCitation struct {
	RefIndex   int     `json:"refIndex"`
	RefText    string  `json:"refText"`
	DocumentId int64   `json:"documentId"`
	ChunkId    int64   `json:"chunkId"`
	Title      string  `json:"title"`
	Snippet    string  `json:"snippet"`
	Score      float64 `json:"score"`
}

type streamReference struct {
	RefIndex int            `json:"refIndex"`
	RefText  string         `json:"refText"`
	Start    int            `json:"start"`
	End      int            `json:"end"`
	Citation streamCitation `json:"citation"`
}

func streamCitationsFromRPC(citations []*pb.Citation) []streamCitation {
	items := make([]streamCitation, 0, len(citations))
	for _, item := range citations {
		if item == nil {
			continue
		}
		refIndex := len(items) + 1
		items = append(items, streamCitation{
			RefIndex:   refIndex,
			RefText:    "[" + strconv.Itoa(refIndex) + "]",
			DocumentId: item.DocumentId,
			ChunkId:    item.ChunkId,
			Title:      item.Title,
			Snippet:    item.Snippet,
			Score:      item.Score,
		})
	}
	return items
}

func streamReferencesFromAnswer(answer string, citations []streamCitation) []streamReference {
	matches := citationMarkerPattern.FindAllStringSubmatchIndex(answer, -1)
	references := make([]streamReference, 0, len(matches))
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}
		refIndex, err := strconv.Atoi(answer[match[2]:match[3]])
		if err != nil || refIndex <= 0 || refIndex > len(citations) {
			continue
		}
		references = append(references, streamReference{
			RefIndex: refIndex,
			RefText:  answer[match[0]:match[1]],
			Start:    utf8.RuneCountInString(answer[:match[0]]),
			End:      utf8.RuneCountInString(answer[:match[1]]),
			Citation: citations[refIndex-1],
		})
	}
	return references
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
