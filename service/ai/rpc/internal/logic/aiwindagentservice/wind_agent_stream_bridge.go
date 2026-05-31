package aiwindagentservicelogic

import (
	"ai-copilot-platform/ai-rpc/internal/engine"
	"ai-copilot-platform/ai-rpc/pb"
)

type windAgentStreamSender interface {
	Send(*pb.WindAgentStreamEvent) error
}

// sendWindAgentStreamEvent 将 Python SSE 事件完整映射为 RPC 事件，不在 Go 层改变 Agent 语义。
func sendWindAgentStreamEvent(stream windAgentStreamSender, event engine.WindAgentStreamEvent) error {
	return stream.Send(windAgentStreamEventToPB(event))
}

func windAgentStreamEventToPB(event engine.WindAgentStreamEvent) *pb.WindAgentStreamEvent {
	return &pb.WindAgentStreamEvent{
		Type:           event.Type,
		TraceId:        event.TraceID,
		ConversationId: event.ConversationID,
		Content:        event.Content,
		ToolCall:       windAgentToolCallToPB(event.ToolCall),
		ToolCalls:      windAgentToolCallsToPB(event.ToolCalls),
		Citations:      windAgentCitationsToPB(event.Citations),
		Draft:          windAgentDraftRefToPB(event.Draft),
		ErrorMsg:       event.ErrorMsg,
	}
}

func windAgentToolCallToPB(toolCall *engine.WindAgentToolCall) *pb.ToolCall {
	if toolCall == nil {
		return nil
	}
	return &pb.ToolCall{
		ToolCallId:    toolCall.ToolCallID,
		ToolName:      toolCall.ToolName,
		Status:        toolCall.Status,
		ArgumentsJson: toolCall.ArgumentsJSON,
		ResultJson:    toolCall.ResultJSON,
		Message:       toolCall.Message,
		LatencyMs:     toolCall.LatencyMS,
	}
}

func windAgentToolCallsToPB(toolCalls []engine.WindAgentToolCall) []*pb.ToolCall {
	result := make([]*pb.ToolCall, 0, len(toolCalls))
	for i := range toolCalls {
		result = append(result, windAgentToolCallToPB(&toolCalls[i]))
	}
	return result
}

func windAgentCitationsToPB(citations []engine.WindAgentCitation) []*pb.Citation {
	result := make([]*pb.Citation, 0, len(citations))
	for _, citation := range citations {
		result = append(result, &pb.Citation{
			DocumentId: citation.DocumentID,
			ChunkId:    citation.ChunkID,
			Title:      citation.Title,
			Snippet:    citation.Snippet,
			Score:      citation.Score,
		})
	}
	return result
}

func windAgentDraftRefToPB(draft *engine.WindAgentDraftRef) *pb.WindDraftRef {
	if draft == nil {
		return nil
	}
	return &pb.WindDraftRef{
		DraftType: draft.DraftType,
		DraftId:   draft.DraftID,
		Title:     draft.Title,
	}
}
