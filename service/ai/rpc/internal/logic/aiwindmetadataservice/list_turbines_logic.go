package aiwindmetadataservicelogic

import (
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"context"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTurbinesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTurbinesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTurbinesLogic {
	return &ListTurbinesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListTurbinesLogic) ListTurbines(in *pb.ListWindTurbineReq) (*pb.ListWindTurbineResp, error) {
	towers, err := l.svcCtx.WindTowerModel.ListTurbines(l.ctx, in.FarmCode, in.Keyword)
	if err != nil {
		l.Logger.Errorf("ListTurbines err %v", err)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "查询风机列表错误")
	}
	items := make([]*pb.WindTurbineItem, 0, len(towers))
	for _, row := range towers {
		items = append(items, &pb.WindTurbineItem{
			TowerId:   row.TowerId,
			TowerCode: row.TowerCode + "号风机",
			FarmId:    row.FarmId,
			FarmCode:  row.FarmCode,
			FarmName:  row.FarmName,
			RiskLevel: row.RiskLevel,
			AiEnabled: row.AiEnabled,
			Remark:    row.Remark,
		})
	}
	return &pb.ListWindTurbineResp{Total: int64(len(items)), List: items, Message: "wind turbine metadata scaffold ready"}, nil
}
