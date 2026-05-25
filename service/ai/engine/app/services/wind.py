import json
from collections import Counter, defaultdict
from statistics import mean
from typing import Any

from app.schemas.wind import (
    WindEvidenceRequest,
    WindHealthReportDraftRequest,
    WindScaffoldResponse,
    WindTicketDraftRequest,
)


# Python Engine 的边界很重要：它不直接访问 PostgreSQL、TDengine 或 ES。
# Go RPC 是事实查询和权限控制层，Python 只消费 Go 传入的 evidence，
# 这样后续接入 LLM 时也不会让模型绕过白名单去“自由查库”。
def _collect_evidence(payload: WindEvidenceRequest) -> list[dict[str, Any]]:
    items: list[dict[str, Any]] = list(payload.evidence)
    if not payload.evidence_json:
        return items

    try:
        parsed = json.loads(payload.evidence_json)
    except json.JSONDecodeError:
        items.append({"type": "invalid_json", "raw": payload.evidence_json})
        return items

    if isinstance(parsed, list):
        items.extend(item for item in parsed if isinstance(item, dict))
    elif isinstance(parsed, dict):
        items.append(parsed)
    return items


def _insufficient(title: str, payload: WindEvidenceRequest) -> WindScaffoldResponse:
    return WindScaffoldResponse(
        title=title,
        status="insufficient_evidence",
        summary="未收到 Go 侧传入的事实证据，不能生成异常判断、归因结论或处置结论。",
        evidence_count=0,
        message="请先由 Go RPC 调用 PG/TDengine/RAG 工具形成 evidence，再调用 Python 生成摘要。",
        sections=[
            {
                "name": "证据要求",
                "content": "至少需要包含时间范围、数据来源、风场/风机/设备标识，以及查询返回的测点或告警记录。",
            }
        ],
        recommendations=["补充 evidence 后重试", "不要让 Python 直接访问 TDengine"],
        todo=["接入更多 evidence 字段解释", "后续可在此处接 LLM 生成更自然的文本"],
    )


def _numeric(value: Any) -> float | None:
    if value is None or isinstance(value, bool):
        return None
    try:
        return float(value)
    except (TypeError, ValueError):
        return None


def _timeseries_points(evidence: list[dict[str, Any]]) -> list[dict[str, Any]]:
    points: list[dict[str, Any]] = []
    for item in evidence:
        if isinstance(item.get("points"), list):
            points.extend(point for point in item["points"] if isinstance(point, dict))
        elif isinstance(item.get("rows"), list):
            points.extend(point for point in item["rows"] if isinstance(point, dict))
        elif "values" in item or "ts" in item:
            points.append(item)
    return points


def _point_values(point: dict[str, Any]) -> dict[str, Any]:
    values = point.get("values")
    if isinstance(values, dict):
        return values
    return {k: v for k, v in point.items() if k not in {"ts", "time", "timestamp"}}


def _summarize_numeric_fields(points: list[dict[str, Any]]) -> dict[str, Any]:
    buckets: dict[str, list[float]] = defaultdict(list)
    latest: dict[str, Any] = {}
    for point in points:
        for field, value in _point_values(point).items():
            number = _numeric(value)
            if number is None:
                continue
            buckets[field].append(number)
            latest[field] = number

    metrics: dict[str, Any] = {}
    for field, values in buckets.items():
        metrics[field] = {
            "count": len(values),
            "latest": latest.get(field),
            "min": min(values),
            "max": max(values),
            "avg": round(mean(values), 6),
            "range": round(max(values) - min(values), 6),
        }
    return metrics


def _risk_from_alarm_level(level_counts: Counter[str]) -> str:
    if any(int(level) >= 4 for level in level_counts if level.isdigit()):
        return "critical"
    if any(int(level) >= 3 for level in level_counts if level.isdigit()):
        return "high"
    if level_counts:
        return "warning"
    return "normal"


async def summarize_timeseries(payload: WindEvidenceRequest) -> WindScaffoldResponse:
    evidence = _collect_evidence(payload)
    if not evidence:
        return _insufficient("测点时序摘要", payload)

    points = _timeseries_points(evidence)
    metrics = _summarize_numeric_fields(points)
    if not metrics:
        return WindScaffoldResponse(
            title="测点时序摘要",
            status="no_numeric_points",
            summary="已收到 evidence，但没有识别到可统计的数值型测点。",
            evidence_count=len(evidence),
            message="请确认 Go 侧 evidence 中包含 points/rows，以及 values 字段或数值列。",
            metrics={"point_count": len(points)},
            sections=[{"name": "原始证据", "content": f"收到 {len(evidence)} 条 evidence，{len(points)} 条疑似测点记录。"}],
            recommendations=["补充数值型字段", "检查 wind_device_meta.column_name 和 TDengine 返回字段是否一致"],
            todo=["补充单位和阈值解释", "接入 LLM 生成更自然的趋势描述"],
        )

    first_field = next(iter(metrics))
    first = metrics[first_field]
    summary = (
        f"共识别 {len(points)} 条测点记录，字段 {first_field} 最新值 {first['latest']}，"
        f"范围 {first['min']} 至 {first['max']}，均值 {first['avg']}。"
    )
    return WindScaffoldResponse(
        title="测点时序摘要",
        status="ok",
        summary=summary,
        evidence_count=len(evidence),
        message="已基于 Go 侧 evidence 完成规则模板统计，未直接访问 TDengine。",
        metrics={"point_count": len(points), "fields": metrics},
        sections=[
            {"name": "统计结果", "content": summary},
            {"name": "AI 边界", "content": "当前结果只基于传入 evidence，不代表 Python 自行查库。"},
        ],
        recommendations=["结合 threshold_config 判断是否超限", "补充前后时间窗用于趋势对比"],
        todo=["接入阈值配置", "接入 LLM 生成异常解释"],
    )


async def summarize_alarm(payload: WindEvidenceRequest) -> WindScaffoldResponse:
    evidence = _collect_evidence(payload)
    if not evidence:
        return _insufficient("告警分析摘要", payload)

    alarms = []
    for item in evidence:
        if isinstance(item.get("list"), list):
            alarms.extend(alarm for alarm in item["list"] if isinstance(alarm, dict))
        elif isinstance(item.get("alarms"), list):
            alarms.extend(alarm for alarm in item["alarms"] if isinstance(alarm, dict))
        elif "alarm_code" in item or "alarmCode" in item:
            alarms.append(item)

    level_counts: Counter[str] = Counter()
    status_counts: Counter[str] = Counter()
    tower_counts: Counter[str] = Counter()
    for alarm in alarms:
        level_counts[str(alarm.get("alarm_level", alarm.get("alarmLevel", "")))] += 1
        status_counts[str(alarm.get("status", ""))] += 1
        tower_counts[str(alarm.get("tower_code", alarm.get("towerCode", alarm.get("tower_id", ""))))] += 1

    risk = _risk_from_alarm_level(level_counts)
    summary = f"识别到 {len(alarms)} 条告警，风险等级建议为 {risk}。"
    return WindScaffoldResponse(
        title="告警分析摘要",
        status="ok" if alarms else "no_alarm_rows",
        summary=summary,
        evidence_count=len(evidence),
        message="已按告警等级、状态和风机位置完成规则聚合。",
        metrics={
            "alarm_count": len(alarms),
            "level_counts": dict(level_counts),
            "status_counts": dict(status_counts),
            "tower_counts": dict(tower_counts),
            "risk": risk,
        },
        sections=[
            {"name": "告警概览", "content": summary},
            {"name": "后续归因需要", "content": "告警前后测点窗口、设备元数据、SOP 检索结果和历史案例。"},
        ],
        recommendations=["先确认高等级告警是否仍处于未恢复状态", "拉取告警前后 30 分钟关键测点趋势"],
        todo=["补充相似 SOP 检索", "补充多证据归因排序"],
    )


async def draft_health_report(payload: WindHealthReportDraftRequest) -> WindScaffoldResponse:
    evidence = _collect_evidence(payload)
    if not evidence:
        return _insufficient("健康报告草稿", payload)

    points = _timeseries_points(evidence)
    metrics = _summarize_numeric_fields(points)
    alarm_summary = await summarize_alarm(payload)
    report_title = f"{payload.farm_code or '风场'} {payload.tower_code or '全部风机'} 健康报告草稿"
    sections = [
        {"name": "一、概览", "content": f"报告类型：{payload.report_type}；时间范围：{payload.start_time or '-'} 至 {payload.end_time or '-'}。"},
        {"name": "二、测点趋势", "content": f"识别到 {len(points)} 条测点记录，数值字段 {len(metrics)} 个。"},
        {"name": "三、告警情况", "content": alarm_summary.summary},
        {"name": "四、风险建议", "content": "当前为规则模板草稿，后续可接入 LLM 生成正式报告语言。"},
        {"name": "五、待补证据", "content": "在线率、缺测率、阈值配置、SOP 引用和历史故障案例。"},
    ]
    return WindScaffoldResponse(
        title=report_title,
        status="draft",
        summary="健康报告草稿已生成，包含概览、测点趋势、告警情况、风险建议和待补证据。",
        evidence_count=len(evidence),
        message="报告由规则模板生成，事实完全来自 Go 侧 evidence。",
        metrics={"timeseries": metrics, "alarm": alarm_summary.metrics},
        sections=sections,
        recommendations=["保存到 ai_health_report 后允许人工编辑", "后续补齐日报、周报、单机报告和故障复盘模板"],
        todo=["接入健康评分", "接入 SOP 引用", "接入 LLM 润色"],
    )


async def draft_ticket(payload: WindTicketDraftRequest) -> WindScaffoldResponse:
    evidence = _collect_evidence(payload)
    if not evidence:
        return _insufficient("维修工单草稿", payload)

    alarm_summary = await summarize_alarm(payload)
    title = f"{payload.farm_code or '风场'} {payload.tower_code or '风机'} 告警处置工单草稿"
    if payload.alarm_code:
        title += f" - {payload.alarm_code}"
    sections = [
        {"name": "问题描述", "content": alarm_summary.summary},
        {"name": "优先级", "content": payload.priority},
        {"name": "建议步骤", "content": "核对告警状态；检查对应设备；拉取前后测点趋势；按 SOP 执行现场确认。"},
        {"name": "证据引用", "content": f"本草稿引用 {len(evidence)} 条 evidence，提交前需人工复核。"},
    ]
    return WindScaffoldResponse(
        title=title,
        status="draft",
        summary="维修工单草稿已生成，包含问题描述、优先级、建议步骤和证据引用。",
        evidence_count=len(evidence),
        message="一期只生成可编辑草稿，不自动提交到工单系统。",
        metrics=alarm_summary.metrics,
        sections=sections,
        recommendations=["人工确认后再提交工单", "补充现场照片、备件和负责人信息"],
        todo=["接入工单系统", "接入 SOP 检索", "补充自动优先级映射"],
    )
