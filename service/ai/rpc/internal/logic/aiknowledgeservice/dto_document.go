package aiknowledgeservicelogic

import (
	"context"
	"strings"
	"time"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
)

const (
	sourceTypeText   = "text"
	sourceTypeUpload = "upload"

	documentStatusParsing  = "parsing"
	documentStatusIndexing = "indexing"
	documentStatusReady    = "ready"
	documentStatusFailed   = "failed"
)

type flatChunk struct {
	ParentDBID int64
	ChunkDBID  int64
	Content    string
}

func canMaintainKnowledgeBase(ctx context.Context, svcCtx *svc.ServiceContext, kb *model.AiKnowledgeBase, userID int64) bool {
	switch strings.ToLower(kb.KbType) {
	case "personal":
		return kb.OwnerUserId.Valid && kb.OwnerUserId.Int64 == userID
	case "public":
		member, err := svcCtx.AiKbMemberModel.FindByKbIDUserID(ctx, kb.Id, userID)
		if err != nil {
			return false
		}
		return member.Role == "editor" || member.Role == "manager"
	default:
		return false
	}
}

func canAccessKnowledgeBase(ctx context.Context, svcCtx *svc.ServiceContext, kb *model.AiKnowledgeBase, userID int64) bool {
	switch strings.ToLower(kb.KbType) {
	case "personal":
		return kb.OwnerUserId.Valid && kb.OwnerUserId.Int64 == userID
	case "public":
		if strings.ToLower(strings.TrimSpace(kb.Visibility)) == "public" {
			return true
		}
		member, err := svcCtx.AiKbMemberModel.FindByKbIDUserID(ctx, kb.Id, userID)
		if err != nil {
			return false
		}
		return member.Role == "viewer" || member.Role == "editor" || member.Role == "manager"
	default:
		return false
	}
}

func documentToPB(doc *model.AiDocument) *pb.DocumentItem {
	if doc == nil {
		return &pb.DocumentItem{}
	}
	return &pb.DocumentItem{
		DocumentId:   doc.Id,
		KbId:         doc.KbId,
		FileName:     doc.FileName,
		FileType:     doc.FileType,
		Status:       doc.Status,
		ErrorMessage: doc.ErrorMsg,
		UploadedBy:   doc.UploadedBy,
		CreatedAt:    formatDocumentTime(doc.CreatedAt),
		UpdatedAt:    formatDocumentTime(doc.UpdatedAt),
	}
}

func formatDocumentTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}
