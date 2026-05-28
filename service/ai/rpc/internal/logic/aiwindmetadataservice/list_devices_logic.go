package aiwindmetadataservicelogic

import (
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"context"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDevicesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDevicesLogic {
	return &ListDevicesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListDevicesLogic) ListDevices(in *pb.ListWindDeviceReq) (*pb.ListWindDeviceResp, error) {

	windDevices, err := l.svcCtx.WindDeviceModel.ListDevices(l.ctx, in.FarmCode, in.TowerCode, in.DeviceTypeCode, in.Keyword)
	if err != nil {
		l.Logger.Errorf("ListDevices err %v", err)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "查询风机设备列表错误")
	}

	items := make([]*pb.WindDeviceItem, 0, len(windDevices))
	for _, row := range windDevices {
		items = append(items, &pb.WindDeviceItem{
			DeviceId:       row.DeviceId,
			DeviceCode:     row.DeviceCode,
			DeviceTypeCode: row.DeviceTypeCode,
			DeviceTypeName: row.DeviceTypeName,
			TowerId:        row.TowerId,
			TowerCode:      row.TowerCode,
			StructureCode:  row.StructureCode,
			StructureName:  row.StructureName,
			Status:         row.Status,
		})
	}
	return &pb.ListWindDeviceResp{Total: int64(len(items)), List: items, Message: "wind device metadata scaffold ready"}, nil
}
