// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_wind_agent

import (
	"context"

	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiWindCreateTicketDraftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindCreateTicketDraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindCreateTicketDraftLogic {
	return &AiWindCreateTicketDraftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindCreateTicketDraftLogic) AiWindCreateTicketDraft(req *types.AiWindTicketDraftReq) (resp *types.AiWindScaffoldResp, err error) {
	result, err := l.svcCtx.AiWindAgentClient.CreateTicketDraft(l.ctx, &pb.WindTicketDraftReq{
		FarmCode:     req.FarmCode,
		TowerCode:    req.TowerCode,
		AlarmCode:    req.AlarmCode,
		EvidenceJson: req.EvidenceJson,
		UserId:       0,
	})
	if err != nil {
		return nil, err
	}
	return &types.AiWindScaffoldResp{
		Id:           result.Id,
		TraceId:      result.TraceId,
		Title:        result.Title,
		Content:      result.Content,
		EvidenceJson: result.EvidenceJson,
		Message:      result.Message,
	}, nil
}
