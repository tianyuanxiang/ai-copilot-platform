package aiwindmetadataservicelogic

import (
	"ai-copilot-platform/ai-rpc/internal/model"
	"context"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

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
	var rows []model.WindTower
	query := l.svcCtx.Orm.Table("wind_turbine").Where("is_delete = 0")
	if farmCode := strings.TrimSpace(in.FarmCode); farmCode != "" {
		query = query.Where("farm_code = ?", strings.ToUpper(farmCode))
	}
	if keyword := strings.TrimSpace(in.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("tower_code like ? or farm_name like ?", like, like)
	}
	if err := query.Order("farm_code asc, tower_code asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]*pb.WindTurbineItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &pb.WindTurbineItem{
			TurbineId: row.TowerId,
			TowerId:   row.TowerId,
			TowerCode: row.TowerCode,
			TowerName: row.TowerCode + "号风机",
			FarmCode:  row.FarmCode,
			FarmName:  row.FarmName,
			RiskLevel: row.RiskLevel,
			AiEnabled: row.AiEnabled,
		})
	}
	return &pb.ListWindTurbineResp{Total: int64(len(items)), List: items, Message: "wind turbine metadata scaffold ready"}, nil
}
