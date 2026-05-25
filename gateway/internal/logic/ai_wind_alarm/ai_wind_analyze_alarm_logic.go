// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_wind_alarm

import (
	"context"

	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiWindAnalyzeAlarmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindAnalyzeAlarmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindAnalyzeAlarmLogic {
	return &AiWindAnalyzeAlarmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindAnalyzeAlarmLogic) AiWindAnalyzeAlarm(req *types.AiWindAlarmAnalyzeReq) (resp *types.AiWindScaffoldResp, err error) {
	result, err := l.svcCtx.AiWindAlarmClient.AnalyzeAlarm(l.ctx, &pb.WindAlarmAnalyzeReq{
		FarmCode:     req.FarmCode,
		TowerCode:    req.TowerCode,
		AlarmCode:    req.AlarmCode,
		EvidenceJson: req.EvidenceJson,
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
