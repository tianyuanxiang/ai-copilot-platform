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

type AiWindListFarmsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindListFarmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindListFarmsLogic {
	return &AiWindListFarmsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindListFarmsLogic) AiWindListFarms(req *types.AiWindListFarmReq) (resp *types.AiWindListFarmResp, err error) {
	result, err := l.svcCtx.AiWindMetadataClient.ListFarms(l.ctx, &pb.ListWindFarmReq{Keyword: req.Keyword})
	if err != nil {
		return nil, err
	}
	list := make([]types.AiWindFarmItem, 0, len(result.List))
	for _, item := range result.List {
		if item == nil {
			continue
		}
		list = append(list, types.AiWindFarmItem{
			FarmId:     item.FarmId,
			FarmCode:   item.FarmCode,
			FarmName:   item.FarmName,
			Province:   item.Province,
			Location:   item.Location,
			TdDatabase: item.TdDatabase,
			AiEnabled:  item.AiEnabled,
		})
	}
	return &types.AiWindListFarmResp{Total: result.Total, List: list, Message: result.Message}, nil
}
