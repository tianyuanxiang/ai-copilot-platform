package aistatusservicelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthLogic {
	return &HealthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Health checks the Python AI engine through its /v1/health endpoint.
func (l *HealthLogic) Health(in *pb.HealthReq) (*pb.HealthResp, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(l.svcCtx.Config.Engine.BaseURL), "/")
	if baseURL == "" {
		return &pb.HealthResp{
			Status:      "degraded",
			EngineReady: false,
			Message:     "AI engine base url is empty",
		}, nil
	}

	timeout := time.Duration(l.svcCtx.Config.Engine.TimeoutMs) * time.Millisecond
	if timeout <= 0 {
		timeout = 3 * time.Second
	}

	req, err := http.NewRequestWithContext(l.ctx, http.MethodGet, baseURL+"/v1/health", nil)
	if err != nil {
		l.Logger.Errorf("Request get ai health failed. %v", err)
		return nil, err
	}

	httpClient := &http.Client{Timeout: timeout}
	engineResp, err := httpClient.Do(req)
	if err != nil {
		l.Logger.Errorf("Send request ai health failed. %v", err)
		return &pb.HealthResp{
			Status:      "degraded",
			EngineReady: false,
			Message:     fmt.Sprintf("AI engine health check failed: %v", err),
		}, err
	}
	defer engineResp.Body.Close()

	if engineResp.StatusCode < http.StatusOK || engineResp.StatusCode >= http.StatusMultipleChoices {
		return &pb.HealthResp{
			Status:      "degraded",
			EngineReady: false,
			Message:     fmt.Sprintf("AI engine health check returned HTTP %d", engineResp.StatusCode),
		}, nil
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(engineResp.Body).Decode(&body); err != nil {
		return &pb.HealthResp{
			Status:      "degraded",
			EngineReady: false,
			Message:     fmt.Sprintf("AI engine health response is invalid: %v", err),
		}, nil
	}

	if !strings.EqualFold(body.Status, "ok") {
		return &pb.HealthResp{
			Status:      "degraded",
			EngineReady: false,
			Message:     fmt.Sprintf("AI engine status is %q", body.Status),
		}, nil
	}

	return &pb.HealthResp{
		Status:      "ok",
		EngineReady: true,
		Message:     "AI engine is ready",
	}, nil
}
