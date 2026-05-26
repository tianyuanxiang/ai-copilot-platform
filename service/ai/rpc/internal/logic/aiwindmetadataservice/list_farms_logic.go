package aiwindmetadataservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFarmsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFarmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFarmsLogic {
	return &ListFarmsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListFarmsLogic) ListFarms(in *pb.ListWindFarmReq) (*pb.ListWindFarmResp, error) {
	rows, err := l.svcCtx.WindFarmModel.ListFarms(l.ctx, in.Keyword)
	if err != nil {
		return nil, err
	}
	items := make([]*pb.WindFarmItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &pb.WindFarmItem{
			FarmId:     row.FarmId,
			FarmCode:   row.FarmCode,
			FarmName:   row.FarmName,
			Province:   row.Province,
			Location:   row.Location,
			TdDatabase: row.TdDatabase,
			AiEnabled:  row.AiEnabled,
		})
	}
	return &pb.ListWindFarmResp{Total: int64(len(items)), List: items, Message: "wind farm metadata scaffold ready"}, nil
}
