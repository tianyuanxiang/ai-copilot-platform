package aichatservicelogic

import (
	"context"
	"strconv"
	"strings"
	"time"

	"ai-copilot-platform/ai-rpc/internal/engine"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RagChatStreamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRagChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RagChatStreamLogic {
	return &RagChatStreamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 执行流式 RAG 问答，逐步返回 token、引用和完成事件。

func (l *RagChatStreamLogic) RagChatStream(in *pb.RagChatReq, stream pb.AiChatService_RagChatStreamServer) error {
	chatLogic := NewRagChatLogic(l.ctx, l.svcCtx)
	run, err := chatLogic.prepareRagChat(in)
	if err != nil {
		return err
	}

	answer := insufficientEvidence
	status := "success"
	errorMsg := ""
	startedAt := time.Now()
	if run.ShouldCallLLM {
		var builder strings.Builder
		_, err := l.svcCtx.EngineCallClient.EngineChatStream(l.ctx, engine.ChatStreamRequest{
			UserID:         strconv.FormatInt(in.UserId, 10),
			KbID:           optionalInt64String(run.EffectiveKbID, run.HasEffectiveKbID),
			ConversationID: strconv.FormatInt(run.ConversationID, 10),
			Question:       run.Prompt,
			History:        chatHistoryToEngine(run.RecentMessages),
		}, func(event engine.ChatStreamEvent) error {
			switch event.Type {
			case "token":
				builder.WriteString(event.Content)
				return stream.Send(&pb.RagChatStreamEvent{
					Type:    "token",
					Content: event.Content,
					TraceId: run.TraceID,
				})
			case "error":
				status = "failed"
				errorMsg = event.Content
				return stream.Send(&pb.RagChatStreamEvent{
					Type:    "error",
					Content: event.Content,
					TraceId: run.TraceID,
				})
			default:
				return nil
			}
		})
		if err != nil {
			status = "failed"
			errorMsg = err.Error()
			chatLogic.writeLlmCallLog(run.TraceID, in.UserId, "chat-stream", run.Prompt, builder.String(), startedAt, status, errorMsg)
			_ = stream.Send(&pb.RagChatStreamEvent{
				Type:    "error",
				Content: err.Error(),
				TraceId: run.TraceID,
			})
			_ = stream.Send(&pb.RagChatStreamEvent{
				Type:    "done",
				TraceId: run.TraceID,
			})
			return nil
		}
		if status == "failed" {
			chatLogic.writeLlmCallLog(run.TraceID, in.UserId, "chat-stream", run.Prompt, builder.String(), startedAt, status, errorMsg)
			_ = stream.Send(&pb.RagChatStreamEvent{
				Type:    "done",
				TraceId: run.TraceID,
			})
			return nil
		}
		answer = strings.TrimSpace(builder.String())
		if answer == "" {
			if run.ShouldUseRag {
				answer = insufficientEvidence
			} else {
				answer = "暂时无法生成回答，请稍后重试。"
			}
		}
	} else if err := stream.Send(&pb.RagChatStreamEvent{
		Type:    "token",
		Content: answer,
		TraceId: run.TraceID,
	}); err != nil {
		return err
	}

	if err := chatLogic.finishRagChat(in.UserId, run, answer, startedAt, status, errorMsg); err != nil {
		return err
	}
	return stream.Send(&pb.RagChatStreamEvent{
		Type:      "done",
		TraceId:   run.TraceID,
		Citations: run.Citations,
	})
}
