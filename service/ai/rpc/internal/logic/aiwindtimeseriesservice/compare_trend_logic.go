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
	queryTimeout = 5000 * time.Second
)

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

func (l *CompareTrendLogic) CompareTrend(in *pb.WindTrendCompareReq) (*pb.WindTrendCompareResp, error) {
	deviceTypeCode := strings.ToUpper(strings.TrimSpace(in.DeviceTypeCode))
	if deviceTypeCode == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "DeviceTypeCode不能为空")
	}
	if len(in.Field) > maxFields {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, fmt.Sprintf("查询字段数量超过上限: %d", maxFields))
	}
	if err := validateTowerCode(in.TowerCode); err != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, err.Error())
	}
	if in.StartTime == "" || in.EndTime == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "查询时间必填，格式为 yyyy-MM-dd HH:mm:ss")
	}

	database := l.svcCtx.WindFarmModel.FarmDatabase(l.ctx, in.FarmCode)
	stableName := model.DeviceTypeStableFallback[deviceTypeCode]
	if stableName == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, fmt.Sprintf("不支持的设备类型: %s", deviceTypeCode))
	}

	displayMeta, err := l.svcCtx.WindDeviceMetaModel.DisplayMetaForDeviceType(l.ctx, deviceTypeCode)
	if err != nil {
		l.Logger.Errorf("DisplayMetaForDeviceType failed: %v", err)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "通过设备类型获取测点元数据失败")
	}
	fields, err := displayMeta.ResolveFields(in.Field)
	if err != nil {
		l.Logger.Errorf("ResolveFields failed: %v", err)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "设备类型不存在该测点字段")
	}

	requestedFields := splitRequestedTrendFields(in.Field)
	analysisFields, excludedFields, usedDefaultFields := resolveTrendAnalysisFields(deviceTypeCode, displayMeta, requestedFields, fields)
	if len(analysisFields) == 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, fmt.Sprintf("设备类型 %s 无可分析测点字段", deviceTypeCode))
	}

	req := trendRequest{
		FarmCode:          in.FarmCode,
		TowerCode:         strings.TrimSpace(in.TowerCode),
		DeviceTypeCode:    deviceTypeCode,
		DeviceCode:        strings.TrimSpace(in.DeviceCode),
		Fields:            fields,
		StartTime:         in.StartTime,
		EndTime:           in.EndTime,
		IndexID:           in.IndexId,
		RadarDistanceM:    in.RadarDistanceM,
		DisplayMeta:       displayMeta,
		RequestedFields:   requestedFields,
		AnalysisFields:    analysisFields,
		ExcludedFields:    excludedFields,
		UsedDefaultFields: usedDefaultFields,
	}

	whereParts := append(model.TimeWhere(in.StartTime, in.EndTime), model.DeviceWhere(req.TowerCode, req.DeviceCode)...)
	where := model.JoinWhere(whereParts)
	if deviceTypeCode == wprDeviceTypeCode {
		resolvedIndexID := in.IndexId
		var matchedDistanceM float64
		if resolvedIndexID == 0 && in.RadarDistanceM > 0 && l.svcCtx.TdengineModel.IsConfigured() {
			resolvedIndexID, matchedDistanceM, err = l.svcCtx.TdengineModel.ResolveRadarIndexByDistance(
				l.ctx, database, stableName, where, float64(in.RadarDistanceM),
			)
			if err != nil {
				l.Logger.Errorf("TdengineModel.ResolveRadarIndexByDistance failed: %v", err)
				return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "通过距离获取雷达 index_id 失败")
			}
		}
		req.ResolvedIndexID = resolvedIndexID
		req.MatchedDistanceM = matchedDistanceM
		if resolvedIndexID > 0 {
			where = fmt.Sprintf("(%s) AND index_id=%d", where, resolvedIndexID)
		}
	}

	durationSeconds := calcDurationSeconds(in.StartTime, in.EndTime)
	plan := planQuery(database, stableName, analysisFields, where, durationSeconds)
	plan.TimeRange = timeRangeInfo{
		StartTime:       in.StartTime,
		EndTime:         in.EndTime,
		DurationSeconds: durationSeconds,
	}

	if plan.isRangeTooLarge() {
		return &pb.WindTrendCompareResp{
			Summary:      "时间范围超过 31 天，请缩小查询范围后重试。",
			EvidenceJson: buildScaffoldEvidence(req, queryModeRangeTooLarge, "range too large"),
			Message:      "time range too large; please limit to 31 days or less",
		}, nil
	}

	if !l.svcCtx.TdengineModel.IsConfigured() {
		l.Logger.Info("TDengine not configured, returning scaffold evidence")
		return &pb.WindTrendCompareResp{
			Summary:      fmt.Sprintf("%s号风机 %s 趋势分析完成，TDengine 未配置。", req.TowerCode, displayMeta.DeviceTypeName),
			EvidenceJson: buildScaffoldEvidence(req, plan.Mode, "TDengine not configured"),
			Message:      "TDEngine is not configured; returning scaffold evidence only",
		}, nil
	}

	queryCtx, cancel := context.WithTimeout(l.ctx, queryTimeout)
	defer cancel()

	if plan.isExactMode() {
		return l.compareExact(queryCtx, plan, req)
	}
	return l.compareMinuteBucket(queryCtx, plan, req)
}

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
		return &pb.WindTrendCompareResp{
			Summary:      fmt.Sprintf("%s号风机 %s 暂无有效统计数据。", req.TowerCode, req.DisplayMeta.DeviceTypeName),
			EvidenceJson: evidence,
			Message:      "trend analysis completed; no data found",
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

func (l *CompareTrendLogic) buildResult(ctx context.Context, plan queryPlan, req trendRequest,
	exactParams *trendCalcParams, bucketParams *trendBucketParams) trendResult {
	var fieldStats map[string]fieldStat
	if exactParams != nil {
		fieldStats = l.calcFieldStatsExact(exactParams)
	} else {
		fieldStats = l.calcFieldStatsFromBuckets(bucketParams, plan.Fields)
	}

	thresholdStatus := "not_configured"
	var thresholdFields map[string]model.FieldThreshold
	if l.svcCtx.ThresholdResolver != nil {
		thr, err := l.svcCtx.ThresholdResolver.ResolveThreshold(ctx, req.DeviceTypeCode, req.TowerCode)
		if err == nil && thr.Status == "configured" {
			thresholdStatus = thr.Status
			thresholdFields = thr.Fields
		}
	}

	riskSignals := make(map[string]riskSignal, len(fieldStats))
	overallRisk := riskSignal{Level: riskNormal}
	for _, field := range plan.Fields {
		fs, ok := fieldStats[field]
		if !ok {
			continue
		}
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

	ts := trendSignal{}
	dq := dataQuality{}
	if fs, ok := firstFieldStat(plan.Fields, fieldStats); ok {
		ts = trendSignal{Direction: fs.Trend, Change: fs.Change, ChangePct: fs.ChangePct}
		dq = calcDataQuality(fs.LastTs, fs.MissingRate)
	}

	agentHints := buildAgentHints(req.TowerCode, req.DisplayMeta, plan.Fields, fieldStats, riskSignals)
	summary := buildSummary(req.TowerCode, req.DisplayMeta, plan.Fields, fieldStats)
	evidence := buildEvidence(plan, req, fieldStats, ts, overallRisk, dq, thresholdStatus, agentHints, false, false)

	return trendResult{
		Summary:      summary,
		EvidenceJSON: evidence,
		Message:      "trend analysis completed",
	}
}

func (l *CompareTrendLogic) calcFieldStatsExact(params *trendCalcParams) map[string]fieldStat {
	fieldStats := make(map[string]fieldStat, len(params.Stats))
	expectedCount := calcExpectedCount(params.DeviceTypeCode, params.DurationSeconds, false)
	for field, stat := range params.Stats {
		fieldStats[field] = calcFieldStat(stat, expectedCount, params.First[field], params.First["ts"], params.Last[field], params.Last["ts"])
	}
	return fieldStats
}

func (l *CompareTrendLogic) calcFieldStatsFromBuckets(params *trendBucketParams, fields []string) map[string]fieldStat {
	durationMinutes := params.DurationSeconds / 60.0
	fieldStats := make(map[string]fieldStat, len(fields))
	for _, field := range fields {
		fieldStats[field] = calcFieldStatFromBuckets(params.Buckets, field, durationMinutes)
	}
	return fieldStats
}

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
