// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_document

import (
	"context"
	"strings"

	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiCheckDocumentImportCandidatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiCheckDocumentImportCandidatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiCheckDocumentImportCandidatesLogic {
	return &AiCheckDocumentImportCandidatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiCheckDocumentImportCandidatesLogic) AiCheckDocumentImportCandidates(req *types.AiCheckDocumentImportCandidatesReq) (resp *types.AiCheckDocumentImportCandidatesResp, err error) {
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if req.KbId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId cannot be empty")
	}
	if len(req.Candidates) > 100 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "at most 100 import candidates are allowed")
	}

	candidates := make([]*pb.DocumentImportCandidateReq, 0, len(req.Candidates))
	for _, item := range req.Candidates {
		candidates = append(candidates, &pb.DocumentImportCandidateReq{
			FileId:     item.FileId,
			FileName:   strings.TrimSpace(item.FileName),
			StoredPath: strings.TrimSpace(item.StoredPath),
		})
	}
	result, err := l.svcCtx.AiKnowledgeClient.CheckDocumentImportCandidates(l.ctx, &aiknowledgeclient.CheckDocumentImportCandidatesReq{
		UserId:     userID,
		KbId:       req.KbId,
		Candidates: candidates,
	})
	if err != nil {
		return nil, err
	}
	items := make([]types.AiDocumentImportCandidateItem, 0, len(result.List))
	for _, item := range result.List {
		if item == nil {
			continue
		}
		items = append(items, types.AiDocumentImportCandidateItem{
			FileId:         item.FileId,
			FileName:       item.FileName,
			StoredPath:     item.StoredPath,
			FileType:       item.FileType,
			Status:         item.Status,
			DocumentId:     item.DocumentId,
			DocumentStatus: item.DocumentStatus,
			Message:        item.Message,
		})
	}
	return &types.AiCheckDocumentImportCandidatesResp{List: items}, nil
}
