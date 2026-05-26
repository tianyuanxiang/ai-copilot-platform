package ai_kb_public

import (
	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	"ai-copilot-platform/gateway/internal/types"
)

func publicKbItemFromRPC(item *aiknowledgeclient.KnowledgeBaseItem) types.AiPublicKbItem {
	if item == nil {
		return types.AiPublicKbItem{}
	}
	return types.AiPublicKbItem{
		KbId:          item.KbId,
		DomainId:      item.DomainId,
		DomainName:    item.DomainName,
		Name:          item.Name,
		Description:   item.Description,
		Visibility:    item.Visibility,
		Status:        int(item.Status),
		DocumentCount: item.DocumentCount,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}
}
