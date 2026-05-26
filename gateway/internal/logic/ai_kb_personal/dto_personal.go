package ai_kb_personal

import (
	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	"ai-copilot-platform/gateway/internal/types"
)

func personalKbItemFromRPC(item *aiknowledgeclient.KnowledgeBaseItem) types.AiPersonalKbItem {
	if item == nil {
		return types.AiPersonalKbItem{}
	}
	return types.AiPersonalKbItem{
		KbId:          item.KbId,
		Name:          item.Name,
		Description:   item.Description,
		Status:        int(item.Status),
		DocumentCount: item.DocumentCount,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}
}
