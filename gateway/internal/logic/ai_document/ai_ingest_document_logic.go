// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_document

import (
	"context"
	"strings"

	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiIngestDocumentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiIngestDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiIngestDocumentLogic {
	return &AiIngestDocumentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiIngestDocumentLogic) AiIngestDocument(req *types.AiIngestDocumentReq) (resp *types.AiTaskResp, err error) {
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if req.KbId <= 0 || strings.TrimSpace(req.FileName) == "" || strings.TrimSpace(req.FileType) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId、fileName、fileType 不能为空")
	}

	sourceType := strings.ToLower(strings.TrimSpace(req.SourceType))
	if sourceType == "" {
		if strings.TrimSpace(req.StoredPath) != "" {
			sourceType = "upload"
		} else {
			sourceType = "text"
		}
	}
	if sourceType != "text" && sourceType != "upload" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "sourceType 只支持 text 或 upload")
	}
	if sourceType == "text" && strings.TrimSpace(req.Content) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "sourceType=text 时 content 不能为空")
	}
	if sourceType == "upload" && strings.TrimSpace(req.StoredPath) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "sourceType=upload 时 storedPath 不能为空")
	}

	task, err := l.svcCtx.AiKnowledgeClient.IngestDocument(l.ctx, &aiknowledgeclient.IngestDocumentReq{
		UserId:      userId,
		KbId:        req.KbId,
		FileName:    strings.TrimSpace(req.FileName),
		FileType:    strings.TrimSpace(req.FileType),
		Content:     req.Content,
		SourceType:  sourceType,
		StoredPath:  strings.TrimSpace(req.StoredPath),
		ContentHash: strings.TrimSpace(req.ContentHash),
		OperatorId:  userId,
	})
	if err != nil {
		return nil, err
	}

	return &types.AiTaskResp{
		TaskId:  task.TaskId,
		Status:  task.Status,
		Message: task.Message,
	}, nil
}
