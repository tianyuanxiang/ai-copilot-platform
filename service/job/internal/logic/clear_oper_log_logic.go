package logic

import (
	"context"
	"fmt"
	"time"

	"ai-copilot-platform/job/internal/lockx"
	"ai-copilot-platform/job/internal/svc"
	sys "ai-copilot-platform/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearOperLogLogic struct {
	svcCtx *svc.ServiceContext
}

func NewClearOperLogLogic(svcCtx *svc.ServiceContext) *ClearOperLogLogic {
	return &ClearOperLogLogic{
		svcCtx: svcCtx,
	}
}

func (l *ClearOperLogLogic) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	lockKey := "job:lock:clear_oper_log"
	owner := fmt.Sprintf("clear_oper_log:%d", time.Now().UnixNano())

	ok, err := lockx.TryAcquire(
		ctx,
		l.svcCtx.RDB,
		lockKey,
		owner,
		time.Duration(l.svcCtx.Config.Jobs.ClearOperLog.LockExpireSeconds)*time.Second,
	)
	if err != nil {
		logx.Errorf("clear oper log acquire lock failed: %v", err)
		return
	}
	if !ok {
		logx.Infof("clear oper log skipped, another instance is running")
		return
	}
	defer func() {
		releaseCtx, releaseCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer releaseCancel()

		if err := lockx.Release(releaseCtx, l.svcCtx.RDB, lockKey, owner); err != nil {
			logx.Errorf("clear oper log release lock failed: %v", err)
		}
	}()

	_, err = l.svcCtx.SysRpc.ClearOperLog(ctx, &sys.ClearOperLogReq{})
	if err != nil {
		logx.Errorf("clear oper log failed: %v", err)
		return
	}

	logx.Infof("clear oper log finished")
}
