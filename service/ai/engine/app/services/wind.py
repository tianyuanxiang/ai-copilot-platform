import json
import logging
from collections import Counter, defaultdict
from statistics import mean
from typing import Any, Literal, TypedDict

from langgraph.graph import START, END, StateGraph

from app.core.config import get_settings
from app.schemas.chat import ChatStreamRequest
from app.schemas.wind import (
    WindEvidenceRequest,
    WindHealthReportDraftRequest,
    WindScaffoldResponse,
    WindTicketDraftRequest,
)
from app.services import llm

"""
从零写一个 LangGraph 流程，思维框架:
  第一步：定义 State
      → 想清楚每个步骤需要读什么、写什么
  第二步：定义节点（流程分几步？每步干什么？）
      → 每个节点是一个函数：读 state → 干活 → 写回 state
  第三步：定义边（步骤之间的顺序是什么？有没有分支？）
      → 线性就用 add_edge
      → 有分支就用 add_conditional_edges
  第四步：设置入口（从哪个节点开始？）
      → set_entry_point
  第五步：compile + invoke
      → 把初始数据塞进 state，启动图
"""

logger = logging.getLogger(__name__)


class WindDraftState(TypedDict, total=False):
    payload: WindEvidenceRequest # 原始请求，始终不变
    task: Literal["timeseries", "alarm", "health_report", "ticket"]
    evidence: list[dict[str, Any]]
    alarms: list[dict[str, Any]]
    points: list[dict[str, Any]]
    metrics: dict[str, Any]
    draft: WindScaffoldResponse


WindDraftTask = Literal["timeseries", "alarm", "health_report", "ticket"]


# Python Engine 的边界很重要：它不直接访问 PostgreSQL、TDengine 或 ES。
# Go RPC 是事实查询和权限控制层，Python 只消费 Go 传入的 evidence，
# 这样后续接入 LLM 时也不会让模型绕过白名单去“自由查库”。
def _collect_evidence(payload: WindEvidenceRequest) -> list[dict[str, Any]]:
    items: list[dict[str, Any]] = list(payload.evidence)
    if items:
        return items
    if not payload.evidence_json: # 如果都不是，返回空列表
        return items

    # 尝试解析 evidence_json 字符串
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
        summary="未收到可用于生成草稿的风机事实证据，不能生成异常判断、归因结论或处置结论。",
        evidence_count=0,
        message="请先由 Go RPC 调用 PG/TDengine 工具形成有效 evidence，再调用 Python 生成草稿。",
        sections=[
            {
                "name": "证据要求",
                "content": "至少需要包含时间范围、数据来源、风场/风机/设备标识，以及查询返回的测点或告警记录。",
            }
        ],
        recommendations=["补充有效 evidence 后重试", "不要让 Python 直接访问 TDengine"],
        todo=["接入更多 evidence 字段解释", "后续可在此处接 LLM 生成更自然的文本"],
    )


def _numeric(value: Any) -> float | None:
    if value is None or isinstance(value, bool):
        return None
    try:
        return float(value)
    except (TypeError, ValueError):
        return None

# 看有没有 points、rows、values、ts 这些 key
def _timeseries_points(evidence: list[dict[str, Any]]) -> list[dict[str, Any]]:
    points: list[dict[str, Any]] = []
    for item in evidence:
        # 告警聚合 evidence 包含 by_level，其 samples 是告警记录而非测点，跳过
        if (
            item.get("source") == "tdengine.alarm"
            or item.get("by_level") is not None
            or (item.get("alarm_count") is not None and "level_counts" in item)
        ):
            continue

        if isinstance(item.get("points"), list):
            points.extend(point for point in item["points"] if isinstance(point, dict))
        elif isinstance(item.get("rows"), list):
            points.extend(point for point in item["rows"] if isinstance(point, dict))
        elif isinstance(item.get("samples"), list):
            points.extend(point for point in item["samples"] if isinstance(point, dict))
        elif "values" in item or "ts" in item:
            points.append(item)

    return points

# 看有没有 alarm_code、list、alarms 这些 key
# 优先使用 Go 侧聚合指标(by_level 或 alarm_count+level_counts)，
# 只有没有聚合指标时才回退扫描 list/alarms 原始行
def _alarms(evidence: list[dict[str, Any]]) -> list[dict[str, Any]]:
    alarms: list[dict[str, Any]] = []
    for item in evidence:
        # 新版聚合 evidence: 包含 by_level 字段
        if item.get("by_level") is not None:
            alarms.append(item)
            continue
        # 旧版聚合 evidence: 包含 alarm_count + level_counts
        if item.get("alarm_count") is not None and "level_counts" in item:
            alarms.append(item)
            continue

        if isinstance(item.get("list"), list):
            alarms.extend(alarm for alarm in item["list"] if isinstance(alarm, dict))
        elif isinstance(item.get("alarms"), list):
            alarms.extend(alarm for alarm in item["alarms"] if isinstance(alarm, dict))
        elif "alarm_code" in item or "alarmCode" in item:
            alarms.append(item)
    return alarms


def _point_values(point: dict[str, Any]) -> dict[str, Any]:
    values = point.get("values")
    if isinstance(values, dict):
        return values
    return {k: v for k, v in point.items() if k not in {"ts", "time", "timestamp"}}


# 对每个数值字段算：count(条数)、min(最小)、max(最大)、avg(均值)、latest(最新值)
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


def _level_value(alarm: dict[str, Any]) -> str:
    value = alarm.get("alarm_level", alarm.get("alarmLevel", ""))
    return "" if value is None else str(value)


def _aggregate_alarm_metrics(alarms: list[dict[str, Any]]) -> dict[str, Any]:
    # 新版聚合 evidence: 包含 by_level 字段（确定性 SQL 聚合结果）
    for alarm in alarms:
        by_level = alarm.get("by_level")
        if by_level is not None and isinstance(by_level, dict):
            level_counter = Counter()
            for k, v in by_level.items():
                level_counter[str(k)] = int(v) if not isinstance(v, int) else v

            by_status = alarm.get("by_status", {})
            status_counts = {str(k): int(v) if not isinstance(v, int) else v for k, v in by_status.items()}

            tower_top = alarm.get("by_tower_top", [])
            tower_counts = {str(t.get("key", "")): t.get("count", 0) for t in tower_top if isinstance(t, dict)}

            return {
                "alarm_count": alarm.get("total", 0),
                "level_counts": dict(level_counter),
                "status_counts": status_counts,
                "tower_counts": tower_counts,
                "risk": _risk_from_alarm_level(level_counter),
                "truncated": alarm.get("truncated", False),
                "returned_records": alarm.get("sampled", 0),
                "granularity": alarm.get("granularity", ""),
                "peak_bucket": alarm.get("peak_bucket"),
                "by_alarm_code_top": alarm.get("by_alarm_code_top", [])[:5],
                "time_bucket_count": alarm.get("time_bucket_count", 0),
                "first_ts": alarm.get("first_ts", ""),
                "last_ts": alarm.get("last_ts", ""),
            }

    # 回退：没有聚合指标时，从原始告警行计算
    level_counts: Counter[str] = Counter()
    status_counts: Counter[str] = Counter()
    tower_counts: Counter[str] = Counter()
    for alarm in alarms:
        level_counts[_level_value(alarm)] += 1
        status_counts[str(alarm.get("status", ""))] += 1
        tower_counts[str(alarm.get("tower_code", alarm.get("towerCode", alarm.get("tower_id", ""))))] += 1

    return {
        "alarm_count": len(alarms),
        "level_counts": dict(level_counts),
        "status_counts": dict(status_counts),
        "tower_counts": dict(tower_counts),
        "risk": _risk_from_alarm_level(level_counts),
    }


def _priority_from_risk(risk: str, fallback: str) -> str:
    if risk == "critical":
        return "urgent"
    if risk == "high":
        return "high"
    if risk == "warning":
        return "normal"
    return fallback or "normal"


def _validate_evidence(state: WindDraftState) -> WindDraftState:
    payload = state["payload"]
    evidence = _collect_evidence(payload)
    logger.info("[validate_evidence] 提取到 %d 条 evidence, task=%s", len(evidence), state["task"])
    if not evidence:
        logger.warning("[validate_evidence] evidence 为空, farm_code=%s tower_code=%s",
                       payload.farm_code, payload.tower_code)
    return {**state, "evidence": evidence}


# 把证据分成两类——告警记录和测点数据
def _normalize_evidence(state: WindDraftState) -> WindDraftState:
    evidence = state.get("evidence", [])
    alarms = _alarms(evidence)
    points = _timeseries_points(evidence)
    logger.info("[normalize_evidence] 从 %d 条 evidence 中拆分出 %d 条告警, %d 条测点",
                len(evidence), len(alarms), len(points))
    return {**state, "alarms": alarms, "points": points}


def _aggregate_facts(state: WindDraftState) -> WindDraftState:
    points = state.get("points", [])
    alarms = state.get("alarms", [])
    alarm_metrics = _aggregate_alarm_metrics(alarms)
    ts_metrics = _summarize_numeric_fields(points)
    logger.info("[aggregate_facts] 告警统计: count=%d risk=%s level_counts=%s",
                alarm_metrics.get("alarm_count", 0),
                alarm_metrics.get("risk", "normal"),
                alarm_metrics.get("level_counts", {}))
    logger.info("[aggregate_facts] 测点统计: point_count=%d fields=%s",
                len(points), list(ts_metrics.keys()))
    return {
        **state,
        "metrics": {
            "point_count": len(points),
            "alarm": alarm_metrics, # 告警统计
            "timeseries": ts_metrics, # 测点统计
        },
    }

# 构建草稿
def _build_draft(state: WindDraftState) -> WindDraftState:
    task = state["task"]
    payload = state["payload"]
    evidence = state.get("evidence", [])
    alarms = state.get("alarms", [])
    points = state.get("points", [])
    metrics = state.get("metrics", {})

    if not evidence or (not alarms and not points and task in {"alarm", "timeseries", "ticket"}):
        title = {
            "alarm": "告警分析摘要",
            "timeseries": "测点时序摘要",
            "ticket": "维修工单草稿",
        }.get(task, "健康报告草稿")
        logger.warning("[build_draft] 证据不足, 返回 insufficient_evidence, task=%s evidence=%d alarms=%d points=%d",
                       task, len(evidence), len(alarms), len(points))
        return {**state, "draft": _insufficient(title, payload)}

    if task == "timeseries":
        draft = _build_timeseries_draft(payload, evidence, points, metrics)
    elif task == "alarm":
        draft = _build_alarm_draft(payload, evidence, alarms, metrics)
    elif task == "health_report":
        draft = _build_health_report_draft(payload, evidence, alarms, points, metrics)
    else:
        draft = _build_ticket_draft(payload, evidence, alarms, metrics)
    logger.info("[build_draft] 草稿构建完成 task=%s status=%s title=%s evidence_count=%d",
                task, draft.status, draft.title, draft.evidence_count)
    return {**state, "draft": draft}


def _build_timeseries_draft(
    payload: WindEvidenceRequest,
    evidence: list[dict[str, Any]],
    points: list[dict[str, Any]],
    metrics: dict[str, Any],
) -> WindScaffoldResponse:
    field_metrics = metrics.get("timeseries", {})
    if not field_metrics:
        return WindScaffoldResponse(
            title="测点时序摘要",
            status="insufficient_evidence",
            summary="已收到 evidence，但没有识别到可统计的数值型测点。",
            evidence_count=len(evidence),
            message="请确认 Go 侧 evidence 中包含 points/rows，以及 values 字段或数值列。",
            metrics={"point_count": len(points)},
            sections=[{"name": "原始证据", "content": f"收到 {len(evidence)} 条 evidence，{len(points)} 条疑似测点记录。"}],
            recommendations=["补充数值型字段", "检查 wind_device_meta.column_name 和 TDengine 返回字段是否一致"],
            todo=["补充单位和阈值解释", "接入 LLM 生成更自然的趋势描述"],
        )

    first_field = next(iter(field_metrics))
    first = field_metrics[first_field]
    summary = (
        f"共识别 {len(points)} 条测点记录，字段 {first_field} 最新值 {first['latest']}，"
        f"范围 {first['min']} 至 {first['max']}，均值 {first['avg']}。"
    )
    return WindScaffoldResponse(
        title="测点时序摘要",
        status="ok",
        summary=summary,
        evidence_count=len(evidence),
        message="已基于 Go 侧 evidence 完成 LangGraph 规则统计，未直接访问 TDengine。",
        metrics={"point_count": len(points), "fields": field_metrics},
        sections=[
            {"name": "统计结果", "content": summary},
            {"name": "AI 边界", "content": "当前结果只基于传入 evidence，不代表 Python 自行查库。"},
        ],
        recommendations=["结合 threshold_config 判断是否超限", "补充前后时间窗用于趋势对比"],
        todo=["接入阈值配置", "接入 LLM 生成异常解释"],
    )

# 把前面算出来的统计数字，填进一个模板里，生成结构化的分析草稿。
def _build_alarm_draft(
    payload: WindEvidenceRequest,
    evidence: list[dict[str, Any]],
    alarms: list[dict[str, Any]],
    metrics: dict[str, Any],
) -> WindScaffoldResponse:
    alarm_metrics = metrics.get("alarm", {})
    if not alarms:
        return _insufficient("告警分析摘要", payload)

    risk = str(alarm_metrics.get("risk", "normal"))
    alarm_count = alarm_metrics.get("alarm_count", len(alarms))
    granularity = alarm_metrics.get("granularity", "")

    summary = f"识别到 {alarm_count} 条告警，风险等级建议为 {risk}。"
    if granularity:
        summary += f" 聚合粒度: {granularity}。"

    sections = [
        {"name": "告警概览", "content": summary},
        {"name": "可能影响", "content": "当前为草稿结论，需结合告警前后测点窗口确认是否存在持续异常。"},
        {"name": "后续归因需要", "content": "告警前后测点窗口、设备元数据、SOP 检索结果和历史案例。"},
    ]

    # 新版聚合指标增强: 告警高峰和高频告警码
    peak = alarm_metrics.get("peak_bucket")
    if peak and isinstance(peak, dict):
        sections.append({
            "name": "告警高峰",
            "content": f"告警集中在 {peak.get('bucket_start', '')} 时段，共 {peak.get('count', 0)} 条。",
        })

    top_codes = alarm_metrics.get("by_alarm_code_top", [])
    if top_codes and isinstance(top_codes, list):
        code_parts = []
        for c in top_codes[:5]:
            if isinstance(c, dict):
                code_parts.append(f"代码{c.get('key', '?')}({c.get('count', 0)}次)")
        if code_parts:
            sections.append({"name": "高频告警码", "content": "、".join(code_parts)})

    return WindScaffoldResponse(
        title="告警分析摘要",
        status="ok",
        summary=summary,
        evidence_count=len(evidence),
        message="已按告警等级、状态和风机位置完成 LangGraph 规则聚合。",
        metrics=alarm_metrics,
        sections=sections,
        recommendations=["优先核对高等级未删除告警记录", "拉取告警前后 30 分钟关键测点趋势"],
        todo=["补充相似 SOP 检索", "补充多证据归因排序"],
    )


def _build_health_report_draft(
    payload: WindEvidenceRequest,
    evidence: list[dict[str, Any]],
    alarms: list[dict[str, Any]],
    points: list[dict[str, Any]],
    metrics: dict[str, Any],
) -> WindScaffoldResponse:
    if not evidence or (not alarms and not points):
        return _insufficient("健康报告草稿", payload)

    alarm_metrics = metrics.get("alarm", {})
    alarm_count = alarm_metrics.get("alarm_count", len(alarms))
    timeseries_metrics = metrics.get("timeseries", {})
    risk = str(alarm_metrics.get("risk", "normal"))
    alarm_summary = f"识别到 {alarm_count} 条告警，风险等级建议为 {risk}。" if alarms else "本次 evidence 未包含可识别告警记录。"
    report_title = f"{payload.farm_code or '风场'} {payload.tower_code or '全部风机'} 健康报告草稿"
    sections = [
        {"name": "一、概览", "content": f"报告类型：{getattr(payload, 'report_type', 'health')}；时间范围：{getattr(payload, 'start_time', '') or '-'} 至 {getattr(payload, 'end_time', '') or '-'}。"},
        {"name": "二、测点趋势", "content": f"识别到 {len(points)} 条测点记录，数值字段 {len(timeseries_metrics)} 个。"},
        {"name": "三、告警情况", "content": alarm_summary},
        {"name": "四、风险建议", "content": "当前为规则模板草稿，提交前需由运维人员结合现场情况复核。"},
        {"name": "五、待补证据", "content": "在线率、缺测率、阈值配置、SOP 引用和历史故障案例。"},
    ]
    return WindScaffoldResponse(
        title=report_title,
        status="draft",
        summary="健康报告草稿已生成，包含概览、测点趋势、告警情况、风险建议和待补证据。",
        evidence_count=len(evidence),
        message="报告由 LangGraph 规则模板生成，事实完全来自 Go 侧 evidence。",
        metrics={"timeseries": timeseries_metrics, "alarm": alarm_metrics},
        sections=sections,
        recommendations=["保存到 ai_health_report 后允许人工编辑", "后续补齐日报、周报、单机报告和故障复盘模板"],
        todo=["接入健康评分", "接入 SOP 引用", "接入 LLM 润色"],
    )


def _build_ticket_draft(
    payload: WindTicketDraftRequest,
    evidence: list[dict[str, Any]],
    alarms: list[dict[str, Any]],
    metrics: dict[str, Any],
) -> WindScaffoldResponse:
    if not alarms:
        return _insufficient("维修工单草稿", payload)

    alarm_metrics = metrics.get("alarm", {})
    alarm_count = alarm_metrics.get("alarm_count", len(alarms))
    risk = str(alarm_metrics.get("risk", "normal"))
    priority = _priority_from_risk(risk, payload.priority)
    summary = f"识别到 {alarm_count} 条告警，建议工单优先级为 {priority}。"
    title = f"{payload.farm_code or '风场'} {payload.tower_code or '风机'} 告警处置工单草稿"
    if payload.alarm_code:
        title += f" - {payload.alarm_code}"
    sections = [
        {"name": "问题描述", "content": summary},
        {"name": "优先级", "content": priority},
        {"name": "建议步骤", "content": "核对告警状态；检查对应设备；拉取前后测点趋势；按 SOP 执行现场确认。"},
        {"name": "证据引用", "content": f"本草稿引用 {len(evidence)} 条 evidence，提交前需人工复核。"},
    ]
    return WindScaffoldResponse(
        title=title,
        status="draft",
        summary="维修工单草稿已生成，包含问题描述、优先级、建议步骤和证据引用。",
        evidence_count=len(evidence),
        message="一期只生成可编辑草稿，不自动提交到工单系统。",
        metrics=alarm_metrics,
        sections=sections,
        recommendations=["人工确认后再提交工单", "补充现场照片、备件和负责人信息"],
        todo=["接入工单系统", "接入 SOP 检索", "补充自动优先级映射"],
    )

# 可选 LLM 润色
async def _optional_llm_polish(state: WindDraftState) -> WindDraftState:
    draft = state["draft"]
    if draft.status == "insufficient_evidence":
        logger.info("[optional_llm_polish] 跳过 LLM 润色, 原因: 证据不足")
        return state

    settings = get_settings()
    if settings.mock_llm or settings.llm_provider == "mock" or not settings.llm_api_key:
        logger.info("[optional_llm_polish] 跳过 LLM 润色, 原因: LLM 未配置 (mock_llm=%s, provider=%s, api_key=%s)",
                    settings.mock_llm, settings.llm_provider, bool(settings.llm_api_key))
        return state

    logger.info("[optional_llm_polish] 开始 LLM 润色, trace_id=%s", state["payload"].trace_id)
    payload = state["payload"]

    # 从 state 构建紧凑事实包（兼容扁平/嵌套 metrics）
    polish_context = _build_polish_context(state)
    if not polish_context:
        logger.warning("[optional_llm_polish] 跳过 LLM 润色, 原因: context 构建失败或超长, trace_id=%s", payload.trace_id)
        return state

    prompt = (
        "你是风机混塔智能运维 Copilot。只能基于下面已生成的草稿和 evidence 指标润色摘要，"
        "不得新增事实、不得编造原因。只输出一段中文摘要,按问题描述、建议步骤、证据引用来分层。\n\n"
        f"标题：{draft.title}\n"
        f"当前摘要：{draft.summary}\n"
        f"事实：{json.dumps(polish_context, ensure_ascii=False)}"
    )

    # 硬性长度保护：截断或跳过 LLM polish
    if len(prompt) > _MAX_POLISH_PROMPT_CHARS:
        logger.warning(
            "[optional_llm_polish] prompt 超长 (%d chars > %d), 跳过 LLM 润色, trace_id=%s",
            len(prompt), _MAX_POLISH_PROMPT_CHARS, payload.trace_id,
        )
        return state

    try:
        request = ChatStreamRequest(
            user_id=payload.user_id,
            kb_id="",
            conversation_id=payload.trace_id,
            question=prompt,
            history=[],
        )
    except Exception as exc:
        logger.warning(
            "[optional_llm_polish] ChatStreamRequest 构造失败, 跳过润色, trace_id=%s error=%s",
            payload.trace_id, exc,
        )
        return state

    try:
        parts: list[str] = []
        async for event in llm.stream_chat(request, trace_id=payload.trace_id or None):
            if event.type == "token":
                parts.append(event.content)
            elif event.type == "error":
                logger.warning("wind llm polish failed trace_id=%s error=%s", payload.trace_id, event.content)
                return state
        polished = "".join(parts).strip()
        if polished:
            original_summary = draft.summary
            logger.info("[optional_llm_polish] 润色前: %s", original_summary)
            logger.info("[optional_llm_polish] 润色后: %s", polished)
            draft.summary = polished
            draft.message += "；已执行可选 LLM 润色"
            logger.info("[optional_llm_polish] LLM 润色完成, 原始长度=%d, 润色后长度=%d",
                        len(original_summary), len(polished))
        else:
            logger.warning("[optional_llm_polish] LLM 返回为空, 保留原始摘要")
    except Exception:
        logger.exception("wind llm polish exception trace_id=%s", payload.trace_id)
    return {**state, "draft": draft}


# 精简 metrics 用于 LLM 润色 prompt，避免超长
_MAX_POLISH_PROMPT_CHARS = 50000
_MAX_SAMPLES_DEFAULT = 10
_MAX_SAMPLES_REDUCED = 3
_SECTION_CONTENT_LIMIT = 200


def _build_polish_context(state: WindDraftState) -> dict[str, Any]:
    """从 state 构建 LLM 润色所需的紧凑事实包。

    直接从 state["metrics"]（嵌套结构）+ state["evidence"] + draft + payload 取事实，
    兼容扁平 alarm metrics 和嵌套 {"alarm": ...} metrics。

    包含：风场/风机/告警码、标题、规则草稿摘要、已有 sections、
    告警总数、风险等级、等级分布、状态分布、top 告警码、top 风机、
    峰值时间桶、首末时间、代表性 samples。
    不放完整 by_time_bucket，只放 time_bucket_count。
    """
    payload = state["payload"]
    metrics = state.get("metrics", {})
    draft = state.get("draft")
    evidence = state.get("evidence", [])

    # 提取告警指标：兼容嵌套 {"alarm": {...}} 和扁平 {"alarm_count": ...} 结构
    alarm_metrics: dict[str, Any] = {}
    if "alarm" in metrics and isinstance(metrics["alarm"], dict):
        # 嵌套结构（来自 _aggregate_facts）
        alarm_metrics = metrics["alarm"]
    elif "alarm_count" in metrics:
        # 扁平结构（兼容旧格式）
        alarm_metrics = metrics

    # 构建紧凑事实包
    context: dict[str, Any] = {}

    # 基本信息
    context["farm_code"] = getattr(payload, "farm_code", "") or ""
    context["tower_code"] = getattr(payload, "tower_code", "") or ""
    context["alarm_code"] = getattr(payload, "alarm_code", "") or ""

    # 草稿信息
    if draft:
        context["title"] = draft.title or ""
        context["summary"] = draft.summary or ""
        # sections 只保留 name 和截断后的 content
        if draft.sections:
            context["sections"] = [
                {
                    "name": s.get("name", ""),
                    "content": str(s.get("content", ""))[:_SECTION_CONTENT_LIMIT],
                }
                for s in draft.sections
                if isinstance(s, dict)
            ]

    # 告警指标摘要
    if alarm_metrics:
        context["alarm_count"] = alarm_metrics.get("alarm_count", 0)
        context["risk"] = alarm_metrics.get("risk", "normal")
        context["level_counts"] = alarm_metrics.get("level_counts", {})
        context["status_counts"] = alarm_metrics.get("status_counts", {})
        context["tower_counts"] = alarm_metrics.get("tower_counts", {})
        context["truncated"] = alarm_metrics.get("truncated", False)
        context["returned_records"] = alarm_metrics.get("returned_records", 0)
        context["granularity"] = alarm_metrics.get("granularity", "")
        context["peak_bucket"] = alarm_metrics.get("peak_bucket")
        context["time_bucket_count"] = alarm_metrics.get("time_bucket_count", 0)
        context["first_ts"] = alarm_metrics.get("first_ts", "")
        context["last_ts"] = alarm_metrics.get("last_ts", "")

        top_codes = alarm_metrics.get("by_alarm_code_top", [])
        if isinstance(top_codes, list):
            context["by_alarm_code_top"] = top_codes[:5]

    # 测点指标摘要
    ts_metrics = metrics.get("timeseries", {})
    if ts_metrics and isinstance(ts_metrics, dict):
        ts_summary: dict[str, Any] = {}
        for field, stats in list(ts_metrics.items())[:20]:
            if isinstance(stats, dict):
                ts_summary[field] = {
                    "count": stats.get("count"),
                    "latest": stats.get("latest"),
                    "min": stats.get("min"),
                    "max": stats.get("max"),
                    "avg": stats.get("avg"),
                }
        if ts_summary:
            context["timeseries"] = ts_summary

    if "point_count" in metrics:
        context["point_count"] = metrics["point_count"]

    # 提取代表性 samples（最多 N 条）
    samples = _extract_samples(evidence, _MAX_SAMPLES_DEFAULT)
    if samples:
        context["samples"] = samples

    # 长度保护：序列化后检查，超长则减少 samples
    serialized = json.dumps(context, ensure_ascii=False)
    if len(serialized) > _MAX_POLISH_PROMPT_CHARS:
        context["samples"] = _extract_samples(evidence, _MAX_SAMPLES_REDUCED)
        serialized = json.dumps(context, ensure_ascii=False)
        if len(serialized) > _MAX_POLISH_PROMPT_CHARS:
            logger.warning(
                "[build_polish_context] context 仍超长 (%d chars), 返回空",
                len(serialized),
            )
            return {}

    return context


def _extract_samples(evidence: list[dict[str, Any]], max_count: int) -> list[dict[str, Any]]:
    """从 evidence 中提取代表性告警 samples，最多 max_count 条。"""
    samples: list[dict[str, Any]] = []
    for item in evidence:
        if len(samples) >= max_count:
            break
        item_samples = item.get("samples", [])
        if isinstance(item_samples, list):
            for s in item_samples:
                if len(samples) >= max_count:
                    break
                if isinstance(s, dict):
                    samples.append(s)
    return samples


def _finalize_response(state: WindDraftState) -> WindDraftState:
    draft = state.get("draft")
    if draft:
        logger.info("[finalize_response] 图执行完成, task=%s status=%s title=%s",
                    state["task"], draft.status, draft.title)
    return state


def _build_graph():
    graph = StateGraph(WindDraftState)
    graph.add_node("validate_evidence", _validate_evidence)
    graph.add_node("normalize_evidence", _normalize_evidence)
    graph.add_node("aggregate_facts", _aggregate_facts)
    graph.add_node("build_draft", _build_draft)
    graph.add_node("optional_llm_polish", _optional_llm_polish)
    graph.add_node("finalize_response", _finalize_response)

    graph.add_edge(START, "validate_evidence")
    graph.add_edge("validate_evidence", "normalize_evidence")
    graph.add_edge("normalize_evidence", "aggregate_facts")
    graph.add_edge("aggregate_facts", "build_draft")
    graph.add_edge("build_draft", "optional_llm_polish")
    graph.add_edge("optional_llm_polish", "finalize_response")
    graph.add_edge("finalize_response", END)
    return graph.compile()


_WIND_DRAFT_GRAPH = _build_graph()


async def _run_wind_graph(payload: WindEvidenceRequest, task: WindDraftTask) -> WindScaffoldResponse:
    logger.info("[wind_graph] 开始执行 task=%s farm_code=%s tower_code=%s trace_id=%s",
                task, payload.farm_code, payload.tower_code, payload.trace_id)
    state = await _WIND_DRAFT_GRAPH.ainvoke({"payload": payload, "task": task})
    draft = state["draft"]
    logger.info("[wind_graph] 执行完成 task=%s status=%s title=%s", task, draft.status, draft.title)
    return draft


async def summarize_timeseries(payload: WindEvidenceRequest) -> WindScaffoldResponse:
    return await _run_wind_graph(payload, "timeseries")


async def summarize_alarm(payload: WindEvidenceRequest) -> WindScaffoldResponse:
    return await _run_wind_graph(payload, "alarm")


async def draft_health_report(payload: WindHealthReportDraftRequest) -> WindScaffoldResponse:
    return await _run_wind_graph(payload, "health_report")


async def draft_ticket(payload: WindTicketDraftRequest) -> WindScaffoldResponse:
    return await _run_wind_graph(payload, "ticket")
