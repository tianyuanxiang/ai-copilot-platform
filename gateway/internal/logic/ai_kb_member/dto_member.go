package ai_kb_member

import (
	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/types"
)

func kbMemberItemFromRPC(item *pb.KbMemberItem) types.AiKbMemberItem {
	if item == nil {
		return types.AiKbMemberItem{}
	}
	return types.AiKbMemberItem{
		KbId:      item.KbId,
		UserId:    item.UserId,
		Username:  item.Username,
		Nickname:  item.Nickname,
		Role:      item.Role,
		CreatedAt: item.CreatedAt,
	}
}
