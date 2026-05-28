// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_wind_metadata

import (
	"context"

	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiWindListTurbinesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindListTurbinesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindListTurbinesLogic {
	return &AiWindListTurbinesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindListTurbinesLogic) AiWindListTurbines(req *types.AiWindListTurbineReq) (resp *types.AiWindListTurbineResp, err error) {
	result, err := l.svcCtx.AiWindMetadataClient.ListTurbines(l.ctx, &pb.ListWindTurbineReq{
		FarmCode: req.FarmCode,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.AiWindTurbineItem, 0, len(result.List))
	for _, item := range result.List {
		if item == nil {
			continue
		}
		list = append(list, types.AiWindTurbineItem{
			TowerId:   item.TowerId,
			TowerCode: item.TowerCode,
			FarmId:    item.FarmId,
			FarmCode:  item.FarmCode,
			FarmName:  item.FarmName,
			RiskLevel: item.RiskLevel,
			AiEnabled: item.AiEnabled,
			Remark:    item.Remark,
		})
	}
	return &types.AiWindListTurbineResp{Total: result.Total, List: list, Message: result.Message}, nil
}
