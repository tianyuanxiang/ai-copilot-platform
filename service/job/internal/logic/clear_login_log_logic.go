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

type ClearLoginLogLogic struct {
	svcCtx *svc.ServiceContext
}

func NewClearLoginLogLogic(svcCtx *svc.ServiceContext) *ClearLoginLogLogic {
	return &ClearLoginLogLogic{
		svcCtx: svcCtx,
	}
}

func (l *ClearLoginLogLogic) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	lockKey := "job:lock:clear_login_log"
	owner := fmt.Sprintf("clear_login_log:%d", time.Now().UnixNano())

	ok, err := lockx.TryAcquire(
		ctx,
		l.svcCtx.RDB,
		lockKey,
		owner,
		time.Duration(l.svcCtx.Config.Jobs.ClearLoginLog.LockExpireSeconds)*time.Second,
	)
	if err != nil {
		logx.Errorf("clear login log acquire lock failed: %v", err)
		return
	}
	if !ok {
		logx.Infof("clear login log skipped, another instance is running")
		return
	}
	defer func() {
		releaseCtx, releaseCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer releaseCancel()

		if err := lockx.Release(releaseCtx, l.svcCtx.RDB, lockKey, owner); err != nil {
			logx.Errorf("clear login log release lock failed: %v", err)
		}
	}()

	_, err = l.svcCtx.SysRpc.ClearLoginLog(ctx, &sys.ClearLoginLogReq{})
	if err != nil {
		logx.Errorf("clear login log failed: %v", err)
		return
	}

	logx.Infof("clear login log finished")
}
