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

type AiWindListDevicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindListDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindListDevicesLogic {
	return &AiWindListDevicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindListDevicesLogic) AiWindListDevices(req *types.AiWindListDeviceReq) (resp *types.AiWindListDeviceResp, err error) {
	result, err := l.svcCtx.AiWindMetadataClient.ListDevices(l.ctx, &pb.ListWindDeviceReq{
		FarmCode:       req.FarmCode,
		TowerCode:      req.TowerCode,
		DeviceTypeCode: req.DeviceTypeCode,
		Keyword:        req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.AiWindDeviceItem, 0, len(result.List))
	for _, item := range result.List {
		if item == nil {
			continue
		}
		list = append(list, types.AiWindDeviceItem{
			DeviceId:       item.DeviceId,
			DeviceCode:     item.DeviceCode,
			DeviceTypeCode: item.DeviceTypeCode,
			DeviceTypeName: item.DeviceTypeName,
			TowerId:        item.TowerId,
			TowerCode:      item.TowerCode,
			StructureCode:  item.StructureCode,
			StructureName:  item.StructureName,
			TdStable:       item.TdStable,
			Status:         int(item.Status),
			AiEnabled:      item.AiEnabled,
		})
	}
	return &types.AiWindListDeviceResp{Total: result.Total, List: list, Message: result.Message}, nil
}
