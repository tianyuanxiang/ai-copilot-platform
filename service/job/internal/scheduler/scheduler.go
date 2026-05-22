package scheduler

import (
	"os"
	"os/signal"
	"syscall"

	"ai-copilot-platform/job/internal/logic"
	"ai-copilot-platform/job/internal/svc"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

type Scheduler struct {
	svcCtx *svc.ServiceContext
	cron   *cron.Cron
}

func NewScheduler(svcCtx *svc.ServiceContext) *Scheduler {
	return &Scheduler{
		svcCtx: svcCtx,
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithChain(
				cron.Recover(cron.DefaultLogger),
			),
		),
	}
}

func (s *Scheduler) Register() {
	if s.svcCtx.Config.Jobs.ClearLoginLog.Enable {
		_, err := s.cron.AddFunc(s.svcCtx.Config.Jobs.ClearLoginLog.Cron, func() {
			logic.NewClearLoginLogLogic(s.svcCtx).Run()
		})
		if err != nil {
			panic(err)
		}
	}

	if s.svcCtx.Config.Jobs.ClearOperLog.Enable {
		_, err := s.cron.AddFunc(s.svcCtx.Config.Jobs.ClearOperLog.Cron, func() {
			logic.NewClearOperLogLogic(s.svcCtx).Run()
		})
		if err != nil {
			panic(err)
		}
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
	defer s.cron.Stop()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	logx.Info("job service stopped")
}
