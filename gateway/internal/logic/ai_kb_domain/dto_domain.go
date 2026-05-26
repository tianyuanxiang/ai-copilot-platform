package ai_kb_domain

import (
	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/types"
)

// domainItemFromRPC 将 RPC DomainItem 映射为网关 AiKbDomainItem
func domainItemFromRPC(item *pb.DomainItem) types.AiKbDomainItem {
	if item == nil {
		return types.AiKbDomainItem{}
	}
	return types.AiKbDomainItem{
		DomainId:    item.DomainId,
		Name:        item.Name,
		Code:        item.Code,
		Description: item.Description,
		Sort:        int(item.Sort),
		Status:      int(item.Status),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
