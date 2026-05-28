package aiwindtimeseriesservicelogic

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	maxFields    = 8
	queryTimeout = 50 * time.Second
)

// CompareTrendLogic 趋势对比 RPC 编排层。
// 负责参数校验、依赖调度和响应组装，业务逻辑在本包内各文件拆分。
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

// CompareTrend 执行风机趋势对比分析。
func (l *CompareTrendLogic) CompareTrend(in *pb.WindTrendCompareReq) (*pb.WindTrendCompareResp, error) {
	// 1. 参数校验
	if in.DeviceTypeCode == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "DeviceTypeCode不能为空")
	}
	if len(in.Field) > maxFields {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, fmt.Sprintf("查询字段数量超过上限（%d）", maxFields))
	}
	if err := validateTowerCode(in.TowerCode); err != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, err.Error())
	}

	req := trendRequest{
		FarmCode:       in.FarmCode,
		TowerCode:      in.TowerCode,
		DeviceTypeCode: in.DeviceTypeCode,
		Fields:         in.Field,
		StartTime:      in.StartTime,
		EndTime:        in.EndTime,
		IndexID:        in.IndexId,
	}

	// 2. 解析数据库/表/字段
	database := l.svcCtx.WindFarmModel.FarmDatabase(l.ctx, in.FarmCode)

	stableName := model.DeviceTypeStableFallback[in.DeviceTypeCode]
	if stableName == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, fmt.Sprintf("不支持的设备类型: %s", in.DeviceTypeCode))
	}

	fields, err := l.svcCtx.WindDeviceMetaModel.FieldsForDeviceType(l.ctx, in.DeviceTypeCode, in.Field)
	if err != nil {
		l.Logger.Errorf("FieldsForDeviceType failed: %v", err)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "通过设备类型获取测点名称失败")
	}
	if len(fields) == 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, fmt.Sprintf("设备类型 %s 无可用测点字段", in.DeviceTypeCode))
	}

	// 3. WHERE 条件
	whereParts := append(model.TimeWhere(in.StartTime, in.EndTime), model.DeviceWhere(in.TowerCode, "")...)
	where := model.JoinWhere(whereParts)

	// 4. 查询策略
	durationSeconds := calcDurationSeconds(in.StartTime, in.EndTime)

	plan := planQuery(database, stableName, fields, where, durationSeconds)

	plan.TimeRange = timeRangeInfo{
		StartTime:       in.StartTime,
		EndTime:         in.EndTime,
		DurationSeconds: durationSeconds,
	}

	// 5. 范围过大直接返回 scaffold evidence
	if plan.isRangeTooLarge() {
		return &pb.WindTrendCompareResp{
			Summary:      "时间范围超过 31 天，请缩小查询范围后重试。",
			EvidenceJson: buildScaffoldEvidence(req, queryModeRangeTooLarge, "range too large"),
			Message:      "time range too large; please limit to 31 days or less",
		}, nil
	}

	// 6. TDengine 未配置，返回 scaffold evidence
	if !l.svcCtx.TdengineModel.IsConfigured() {
		l.Logger.Info("TDengine not configured, returning scaffold evidence")

		return &pb.WindTrendCompareResp{
			Summary:      fmt.Sprintf("%s %s 趋势分析完成（TDengine 未配置）。", in.TowerCode, in.DeviceTypeCode),
			EvidenceJson: buildScaffoldEvidence(req, plan.Mode, "TDengine not configured"),
			Message:      "TDEngine is not configured; returning scaffold evidence only",
		}, nil
	}

	// 7. 执行查询
	queryCtx, cancel := context.WithTimeout(l.ctx, queryTimeout)
	defer cancel()

	if plan.isExactMode() {
		return l.compareExact(queryCtx, plan, req)
	}
	return l.compareMinuteBucket(queryCtx, plan, req)
}

// compareExact 精确模式：basic stats + boundary rows。
func (l *CompareTrendLogic) compareExact(ctx context.Context, plan queryPlan, req trendRequest) (*pb.WindTrendCompareResp, error) {
	stats, err := l.svcCtx.TdengineModel.QueryBasicStats(ctx, plan.Database, plan.Stable, plan.Fields, plan.Where)
	if err != nil {
		l.Logger.Errorf("QueryBasicStats failed: %v", err)
		return l.handleQueryError(plan, req, err, stats != nil)
	}

	first, last, err := l.svcCtx.TdengineModel.QueryBoundaryRows(ctx, plan.Database, plan.Stable, plan.Fields, plan.Where)
	if err != nil {
		l.Logger.Errorf("QueryBoundaryRows failed: %v", err)
		return l.handleQueryError(plan, req, err, stats != nil)
	}

	result := l.buildResult(ctx, plan, req, &trendCalcParams{
		Stats:           stats,
		First:           first,
		Last:            last,
		DurationSeconds: plan.TimeRange.DurationSeconds,
		DeviceTypeCode:  req.DeviceTypeCode,
		IsMinuteMode:    false,
	}, nil)

	return &pb.WindTrendCompareResp{
		Summary:      result.Summary,
		EvidenceJson: result.EvidenceJSON,
		Message:      result.Message,
	}, nil
}

// compareMinuteBucket 分钟级降采样模式。
func (l *CompareTrendLogic) compareMinuteBucket(ctx context.Context, plan queryPlan, req trendRequest) (*pb.WindTrendCompareResp, error) {
	buckets, err := l.svcCtx.TdengineModel.QueryMinuteBuckets(ctx, plan.Database, plan.Stable, plan.Fields, plan.Where)
	if err != nil {
		l.Logger.Errorf("QueryMinuteBuckets failed: %v", err)
		return l.handleQueryError(plan, req, err, len(buckets) > 0)
	}

	if len(buckets) == 0 {
		fieldStats := make(map[string]fieldStat)
		for _, f := range plan.Fields {
			fieldStats[f] = fieldStat{}
		}
		evidence := buildEvidence(plan, req, fieldStats, trendSignal{}, riskSignal{Level: riskNormal},
			dataQuality{}, "not_configured", nil, false, false)
		result := trendResult{
			Summary:      fmt.Sprintf("%s %s 暂无有效统计数据。", req.TowerCode, req.DeviceTypeCode),
			EvidenceJSON: evidence,
			Message:      "trend analysis completed; no data found",
		}
		return &pb.WindTrendCompareResp{
			Summary:      result.Summary,
			EvidenceJson: result.EvidenceJSON,
			Message:      result.Message,
		}, nil
	}

	result := l.buildResult(ctx, plan, req, nil, &trendBucketParams{
		Buckets:         buckets,
		DurationSeconds: plan.TimeRange.DurationSeconds,
		DeviceTypeCode:  req.DeviceTypeCode,
	})

	return &pb.WindTrendCompareResp{
		Summary:      result.Summary,
		EvidenceJson: result.EvidenceJSON,
		Message:      result.Message,
	}, nil
}

// buildResult 组装最终结果。exactParams 和 bucketParams 二选一传值。
func (l *CompareTrendLogic) buildResult(ctx context.Context, plan queryPlan, req trendRequest,
	exactParams *trendCalcParams, bucketParams *trendBucketParams) trendResult {
	var fieldStats map[string]fieldStat

	if exactParams != nil {
		fieldStats = l.calcFieldStatsExact(exactParams)
	} else {
		fieldStats = l.calcFieldStatsFromBuckets(bucketParams, plan.Fields)
	}

	// 阈值加载
	thresholdStatus := "not_configured"
	var thresholdFields map[string]model.FieldThreshold
	if l.svcCtx.ThresholdResolver != nil {
		thr, err := l.svcCtx.ThresholdResolver.ResolveThreshold(ctx, req.DeviceTypeCode, req.TowerCode)
		if err == nil && thr.Status == "configured" {
			thresholdStatus = thr.Status
			thresholdFields = thr.Fields
		}
	}

	// 风险信号
	riskSignals := make(map[string]riskSignal, len(fieldStats))
	var overallRisk riskSignal
	for field, fs := range fieldStats {
		var ft *model.FieldThreshold
		if tf, ok := thresholdFields[field]; ok {
			ft = &tf
		}
		rs := calcRiskSignal(fs, ft)
		riskSignals[field] = rs
		if riskLevelOrder(rs.Level) > riskLevelOrder(overallRisk.Level) {
			overallRisk = rs
		}
	}
	if overallRisk.Level == "" {
		overallRisk.Level = riskNormal
	}

	// 趋势信号（取第一个字段）
	var ts trendSignal
	for _, fs := range fieldStats {
		ts = trendSignal{
			Direction: fs.Trend,
			Change:    fs.Change,
			ChangePct: fs.ChangePct,
		}
		break
	}

	// 数据质量（取第一个字段的时间戳）
	var dq dataQuality
	for _, fs := range fieldStats {
		dq = calcDataQuality(fs.LastTs, fs.MissingRate)
		break
	}

	agentHints := buildAgentHints(req.TowerCode, req.DeviceTypeCode, fieldStats, riskSignals)
	summary := buildSummary(req.TowerCode, req.DeviceTypeCode, fieldStats)

	evidence := buildEvidence(plan, req, fieldStats, ts, overallRisk, dq, thresholdStatus, agentHints, false, false)

	return trendResult{
		Summary:      summary,
		EvidenceJSON: evidence,
		Message:      "trend analysis completed",
	}
}

// calcFieldStatsExact 精确模式：根据基础统计和边界行计算字段统计。
func (l *CompareTrendLogic) calcFieldStatsExact(params *trendCalcParams) map[string]fieldStat {
	fieldStats := make(map[string]fieldStat, len(params.Stats))
	expectedCount := calcExpectedCount(params.DeviceTypeCode, params.DurationSeconds, false)

	for field, stat := range params.Stats {
		fieldStats[field] = calcFieldStat(stat, expectedCount, params.First[field], params.First["ts"], params.Last[field], params.Last["ts"])
	}

	return fieldStats
}

// calcFieldStatsFromBuckets 分钟级模式：根据 bucket 数据计算字段统计。
func (l *CompareTrendLogic) calcFieldStatsFromBuckets(params *trendBucketParams, fields []string) map[string]fieldStat {
	durationMinutes := params.DurationSeconds / 60.0
	fieldStats := make(map[string]fieldStat, len(fields))

	for _, field := range fields {
		fieldStats[field] = calcFieldStatFromBuckets(params.Buckets, field, durationMinutes)
	}

	return fieldStats
}

// handleQueryError 处理查询错误，返回 partial 或完全失败的 evidence。
func (l *CompareTrendLogic) handleQueryError(plan queryPlan, req trendRequest, err error, hasPartialData bool) (*pb.WindTrendCompareResp, error) {
	if hasPartialData {
		partialEvidence := buildPartialEvidence(plan, req, nil, trendSignal{}, riskSignal{Level: riskNormal}, dataQuality{}, nil)
		return &pb.WindTrendCompareResp{
			Summary:      "趋势统计查询超时，未生成有效趋势结论。",
			EvidenceJson: partialEvidence,
			Message:      "trend query timeout; returned partial evidence",
		}, nil
	}
	return &pb.WindTrendCompareResp{
		Summary:      fmt.Sprintf("趋势统计查询失败: %v", err),
		EvidenceJson: buildScaffoldEvidence(req, plan.Mode, fmt.Sprintf("query error: %v", err)),
		Message:      fmt.Sprintf("trend query error: %v", err),
	}, nil
}

// validateTowerCode 校验 towerCode，非空时必须为数字字符串。
func validateTowerCode(towerCode string) error {
	towerCode = strings.TrimSpace(towerCode)
	if towerCode == "" {
		return nil
	}
	_, err := strconv.ParseInt(towerCode, 10, 64)
	if err != nil {
		return fmt.Errorf("towerCode 必须为数字: %s", towerCode)
	}
	return nil
}
