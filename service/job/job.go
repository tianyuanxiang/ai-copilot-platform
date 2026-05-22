package main

import (
	"flag"
	"fmt"

	"ai-copilot-platform/job/internal/config"
	"ai-copilot-platform/job/internal/scheduler"
	"ai-copilot-platform/job/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/job.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)

	s := scheduler.NewScheduler(ctx)
	s.Register()

	logx.Infof("starting job service...")
	fmt.Println("Starting job service...")

	s.Start()
}
