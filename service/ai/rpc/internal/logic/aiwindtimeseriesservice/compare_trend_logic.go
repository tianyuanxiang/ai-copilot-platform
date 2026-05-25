package aiwindtimeseriesservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompareTrendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompareTrendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompareTrendLogic {
	return &CompareTrendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CompareTrendLogic) CompareTrend(in *pb.WindTrendCompareReq) (*pb.WindTrendCompareResp, error) {
	database := l.svcCtx.WindMetadataModel.FarmDatabase(l.ctx, in.FarmCode)
	stable := l.svcCtx.WindMetadataModel.StableForDeviceType(l.ctx, in.DeviceTypeCode)
	fields := l.svcCtx.WindMetadataModel.FieldsForDeviceType(l.ctx, in.DeviceTypeCode, in.Field)
	evidence := map[string]any{
		"source":     "tdengine",
		"scaffold":   true,
		"farm_code":  in.FarmCode,
		"tower_code": in.TowerCode,
		"database":   database,
		"stable":     stable,
		"fields":     fields,
		"todo":       "后续补充最大值、最小值、均值、波动幅度、缺测率和持续超限判断。",
	}
	return &pb.WindTrendCompareResp{
		Summary:      "测点趋势对比脚手架已就绪，一期暂不执行完整趋势判断。",
		EvidenceJson: model.WindEvidenceJSON(evidence),
		Message:      "compare_sensor_trend scaffold ready",
	}, nil
}
