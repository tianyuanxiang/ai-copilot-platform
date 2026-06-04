package aiknowledgeservicelogic

import (
	"context"
	"path/filepath"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckDocumentImportCandidatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckDocumentImportCandidatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckDocumentImportCandidatesLogic {
	return &CheckDocumentImportCandidatesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckDocumentImportCandidatesLogic) CheckDocumentImportCandidates(in *pb.CheckDocumentImportCandidatesReq) (*pb.CheckDocumentImportCandidatesResp, error) {
	if in.UserId <= 0 || in.KbId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "user_id and kb_id cannot be empty")
	}
	if len(in.Candidates) > 100 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "at most 100 import candidates are allowed")
	}

	kb, err := l.svcCtx.AiKnowledgeBaseModel.FindByID(l.ctx, in.KbId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "knowledge base not found")
		}
		return nil, err
	}
	canMaintain, err := canMaintainKnowledgeBase(l.ctx, l.svcCtx, kb, in.UserId)
	if err != nil {
		return nil, err
	}
	if !canMaintain {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "no permission to maintain this knowledge base")
	}

	reader := NewIngestDocumentLogic(l.ctx, l.svcCtx)
	items := make([]*pb.DocumentImportCandidateItem, 0, len(in.Candidates))
	for _, candidate := range in.Candidates {
		item := &pb.DocumentImportCandidateItem{}
		if candidate == nil {
			item.Status = "unavailable"
			item.Message = "candidate is empty"
			items = append(items, item)
			continue
		}
		item.FileId = candidate.FileId
		item.FileName = strings.TrimSpace(candidate.FileName)
		item.StoredPath = strings.TrimSpace(candidate.StoredPath)
		item.FileType = normalizeFileType("", item.FileName)
		if item.FileType == "" {
			item.FileType = strings.ToLower(strings.TrimPrefix(filepath.Ext(item.StoredPath), "."))
		}
		if !isSupportedUploadedFileType(item.FileType) {
			item.Status = "unsupported"
			item.Message = "only md, markdown, txt and docx files can be imported"
			items = append(items, item)
			continue
		}
		data, err := reader.readUploadedFile(item.StoredPath)
		if err != nil {
			item.Status = "unavailable"
			item.Message = err.Error()
			items = append(items, item)
			continue
		}
		existing, err := l.svcCtx.AiDocumentModel.FindByKbIDContentHash(l.ctx, in.KbId, sha256Hex(data))
		if err != nil && err != model.ErrNotFound {
			return nil, err
		}
		if existing != nil {
			item.Status = "imported"
			item.DocumentId = existing.Id
			item.DocumentStatus = existing.Status
			item.Message = "document already exists in this knowledge base"
		} else {
			item.Status = "available"
			item.Message = "document can be imported"
		}
		items = append(items, item)
	}

	return &pb.CheckDocumentImportCandidatesResp{List: items}, nil
}
