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

func canMaintainKnowledgeBase(ctx context.Context, svcCtx *svc.ServiceContext, kb *model.AiKnowledgeBase, userID int64) (bool, error) {
	switch strings.ToLower(kb.KbType) {
	case "personal":
		return isKnowledgeBaseOwner(kb, userID), nil
	case "public":
		if isKnowledgeBaseOwner(kb, userID) {
			return true, nil
		}
		member, err := svcCtx.AiKbMemberModel.FindByKbIDUserID(ctx, kb.Id, userID)
		if err != nil {
			if err == model.ErrNotFound {
				return false, nil
			}
			return false, err
		}
		return member.Role == "editor" || member.Role == "manager", nil
	default:
		return false, nil
	}
}

func canManageKnowledgeBaseMembers(ctx context.Context, svcCtx *svc.ServiceContext, kb *model.AiKnowledgeBase, userID int64) (bool, error) {
	if isKnowledgeBaseOwner(kb, userID) {
		return true, nil
	}
	if strings.ToLower(kb.KbType) != "public" {
		return false, nil
	}
	member, err := svcCtx.AiKbMemberModel.FindByKbIDUserID(ctx, kb.Id, userID)
	if err != nil {
		if err == model.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return member.Role == "manager", nil
}

func isKnowledgeBaseOwner(kb *model.AiKnowledgeBase, userID int64) bool {
	if kb == nil || userID <= 0 {
		return false
	}
	if kb.OwnerUserId.Valid && kb.OwnerUserId.Int64 == userID {
		return true
	}
	return kb.CreatedBy > 0 && kb.CreatedBy == userID
}

func canAccessKnowledgeBase(ctx context.Context, svcCtx *svc.ServiceContext, kb *model.AiKnowledgeBase, userID int64) bool {
	switch strings.ToLower(kb.KbType) {
	case "personal":
		return isKnowledgeBaseOwner(kb, userID)
	case "public":
		if isKnowledgeBaseOwner(kb, userID) {
			return true
		}
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

func domainToPB(d *model.AiKbDomain) *pb.DomainItem {
	if d == nil {
		return &pb.DomainItem{}
	}
	return &pb.DomainItem{
		DomainId:    d.Id,
		Name:        d.Name,
		Code:        d.Code,
		Description: d.Description,
		Sort:        d.Sort,
		Status:      d.Status,
		CreatedAt:   formatDocumentTime(d.CreatedAt),
		UpdatedAt:   formatDocumentTime(d.UpdatedAt),
	}
}

func knowledgeBaseToPB(kb *model.AiKnowledgeBase, docCount int64, domainName string) *pb.KnowledgeBaseItem {
	if kb == nil {
		return &pb.KnowledgeBaseItem{}
	}
	return &pb.KnowledgeBaseItem{
		KbId:          kb.Id,
		KbType:        kb.KbType,
		OwnerUserId:   model.NullInt64Value(kb.OwnerUserId),
		DomainId:      model.NullInt64Value(kb.DomainId),
		DomainName:    domainName,
		Name:          kb.Name,
		Description:   kb.Description,
		Visibility:    kb.Visibility,
		Status:        kb.Status,
		DocumentCount: docCount,
		CreatedAt:     formatDocumentTime(kb.CreatedAt),
		UpdatedAt:     formatDocumentTime(kb.UpdatedAt),
	}
}

func formatDocumentTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}
