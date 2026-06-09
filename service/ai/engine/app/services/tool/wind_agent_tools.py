"""Tool schemas exposed to the Wind ReAct planner.

The eight domain tools are executed by Go. The two control actions are handled
inside LangGraph and never leave the Python process.
"""

from __future__ import annotations

from typing import Any


READ_ONLY_TOOLS = {
    "get_turbine_metadata",
    "search_maintenance_sop",
    "query_alarm_events",
    "query_sensor_timeseries",
    "compare_sensor_trend",
}

DRAFT_TOOLS = {
    "generate_alarm_analysis_draft",
    "generate_health_report",
    "create_maintenance_ticket_draft",
}

GO_TOOLS = READ_ONLY_TOOLS | DRAFT_TOOLS
CONTROL_ACTIONS = {"request_clarification", "finish_answer"}
ALL_PLANNER_ACTIONS = GO_TOOLS | CONTROL_ACTIONS


def _function(name: str, description: str, properties: dict[str, Any], required: list[str] | None = None) -> dict[str, Any]:
    return {
        "type": "function",
        "function": {
            "name": name,
            "description": description,
            "parameters": {
                "type": "object",
                "properties": properties,
                "required": required or [],
                "additionalProperties": False,
            },
        },
    }


AGENT_TOOL_SCHEMAS = [
    _function(
        "get_turbine_metadata",
        "查询风场、风机和设备元数据。当用户想要获取或查询风场、风机、设备类型数据的时候，通过该工具获取对应的元数据，"
        "此外，还用于补齐用户没有明确给出的 farmCode、towerCode、deviceTypeCode。仅当需要确认可用风机/风场/设备类型编号或风机列表时调用；不要用它查询告警、时序或 SOP。",
        {
            "farmCode": {
                "type": "string",
                "description": "风场编码，例如 FY 表示扶余风场，YS 表示榆树风场。未知时可省略以列出可用风场。",
            },
            "towerCode": {
                "type": "string",
                "description": "两位风机编号，例如 04、05。用户说4号风机、四号风机时必须转换为 04。",
                "pattern": "^[0-9]{2}$"
            },
            "deviceTypeCode": {
                "type": "string",
                "description": "传感器类型的编码，例如 WPR、ACCX、ACCY、HLS、STM、INSX、INSY、JMT、ATS、GNSS。需要筛选设备类型时提供。"
            },
        },
    ),
    _function(
        "search_maintenance_sop",
        "检索维护SOP、检修规程、处置步骤和安全注意事项。"
        "用户询问“怎么处理、操作规程、操作手册、检修步骤、SOP、维护要求”时调用；不要用它查询实时告警或传感器数据。"
        "返回 chunks 和 citations，最终回答引用 SOP 时 evidenceRequirement 应为 sop。",
        {
            "query": {
                "type": "string",
                "description": "面向SOP检索的自然语言查询，包含设备、告警、故障现象和想查的规程。"
            },
            "kbId": {
                "type": "integer",
                "description": "可选的知识库 ID；未知时省略。"
            },
            "searchScope": {
                "type": "string",
                "description": "检索范围；默认 public。",
                "enum": ["public", "personal", "all"]
            },
            "domainId": {
                "type": "integer",
                "description": "可选知识库领域 ID；未知时省略。"
            },
            "documentIds": {
                "type": "array",
                "description": "可选文档 ID 白名单。未知时省略。",
                "items": {"type": "integer"}
            },
            "topK": {
                "type": "integer",
                "description": "召回条数；默认 5，最大 20。"
            },
            "answerMode": {
                "type": "string",
                "description": "quick 用于普通检索，deep 用于复杂规程或多步骤问题。",
                "enum": ["quick", "deep"]
            },
        },
        ["query"],
    ),
    _function(
        "query_alarm_events",
        "查询 TDengine 告警 evidence，包括聚合统计、告警级别/状态分布、时间桶和抽样告警。"
        "用户询问具体风场/风机的告警、报警次数、告警趋势或告警证据时调用。需要 farmCode；如用户只说风场中文名，先用 FY=扶余风场、YS=榆树风场映射，无其他风场，如果用户询问其他风场，则返回无法处理。",
        {
            "farmCode": {"type": "string", "enum": ["FY", "YS"],"description": "风场编码。FY=扶余风场，YS=榆树风场。"},
            "towerCode": {"type": "string", "description": "可选风机编码。例如4号风机或者四号风机，需要转为04"},
            "alarmCode": {"type": "string","description": "可选告警编码。"},
            "startTime": {"type": "string","description": "开始时间，格式 yyyy-MM-dd HH:mm:ss 或 RFC3339；和 endTime 必须同时提供；默认最近 24 小时。"},
            "endTime": {"type": "string","description": "结束时间，必须晚于 startTime；查询范围最大 31 天。"},
            "status": {"type": "integer", "description": "可选告警状态值,已删除为1，未删除为0。若要按状态筛选，必须同时设置 hasStatus=true。"},
            "hasStatus": {"type": "boolean", "description": "是否启用 status 过滤；用于区分 status=0 和未传 status。"},
        },
        ["farmCode"],
    ),
    _function(
        "query_sensor_timeseries",
        "查询传感器/设备测点的原始 TDengine 时序数据。用户需要具体时间点、曲线数据、原始测点值或分页数据时调用；如果用户要趋势摘要、阈值解析或风险判断，优先调用 compare_sensor_trend。",
        {
            "farmCode": {"type": "string", "enum": ["FY", "YS"], "description": "风场编码。"},
            "towerCode": {"type": "string", "description": "可选风机编码。"},
            "deviceCode": {"type": "string", "description": "可选设备编码。"},
            "deviceTypeCode": {"type": "string", "description": "设备类型编码，必填，例如 WPR、ACCX。未知时先查 get_turbine_metadata。"},
            "field": {"type": "array", "items": {"type": "string"}, "maxItems": 8, "description": "测点字段名列表，最多 8 个；不知道 TDengine的列名时必须省略，不要编造字段。"},
            "startTime": {"type": "string", "description": "开始时间；默认最近 24 小时；最大 31 天。"},
            "endTime": {"type": "string", "description": "结束时间；必须和 startTime 同时提供。"},
            "page": {"type": "integer", "description": "页码，默认 1。"},
            "pageSize": {"type": "integer", "description": "每页条数，默认 500，最大 1000。"},
            "indexId": {"type": "integer", "description": "WPR 雷达距离层编号 1-10。用户说真实距离时优先传 radarDistanceM。"},
            "radarDistanceM": {"type": "integer", "description": "WPR 雷达真实探测距离，如 100、300；Go 侧会解析对应 indexId。"},
        },
        ["farmCode", "deviceTypeCode"],
    ),
    _function(
        "compare_sensor_trend",
        "对传感器测点做趋势对比、阈值解析和风险判断。用户询问“是否异常、趋势如何、风险等级、相比之前是否升高/下降、是否超过阈值”时调用；不要用于大量原始点位导出。",
        {
            "farmCode": {"type": "string", "enum": ["FY", "YS"], "description": "风场编码。"},
            "towerCode": {"type": "string", "description": "可选风机编码。"},
            "deviceCode": {"type": "string", "description": "可选设备编码。"},
            "deviceTypeCode": {"type": "string", "description": "设备类型编码，必填。"},
            "field": {"type": "array", "items": {"type": "string"}, "maxItems": 8, "description": "需要分析的测点字段，最多 8 个；不知道 TDengine的列名时必须省略，不要编造字段。"},
            "startTime": {"type": "string", "description": "分析开始时间；默认最近 24 小时；最大 31 天。"},
            "endTime": {"type": "string", "description": "分析结束时间；必须和 startTime 同时提供。"},
            "indexId": {"type": "integer", "description": "WPR 距离层编号 1-10。"},
            "radarDistanceM": {"type": "integer", "description": "WPR 真实探测距离，例如 100、200、250、300等；Go 侧会解析对应 indexId。"},
        },
        ["farmCode", "deviceTypeCode"],
    ),
    _function(
        "generate_alarm_analysis_draft",
        "根据已有或自动补充的告警 evidence 生成并保存告警分析草稿。只有用户明确要求生成/保存分析草稿时才选择；执行前，LangGraph 会进入人工确认；不会直接发布结论或执行现场操作。",
        {
            "farmCode": {"type": "string","enum": ["FY", "YS"], "description": "风场编码。"},
            "towerCode": {"type": "string", "description": "可选风机编码。"},
            "alarmCode": {"type": "string", "description": "可选告警编码。"},
            "startTime": {"type": "string", "description": "证据时间范围开始。"},
            "endTime": {"type": "string", "description": "证据时间范围结束。"},
            "evidenceJson": {"type": "string", "description": "已有工具 evidence 的 JSON 字符串；没有时可传 {}，Go 侧会按参数补充告警证据。"},
        },
        ["farmCode"],
    ),
    _function(
        "generate_health_report",
        "生成并保存健康报告草稿。只有用户明确要求生成健康报告/巡检报告/状态报告时调用；执行前需要人工确认；返回草稿引用，不代表报告已审核发布。",
        {
            "reportType": {"type": "string", "description": "报告类型，默认 health。"},
            "farmCode": {"type": "string", "enum": ["FY", "YS"], "description": "风场编码。"},
            "towerCode": {"type": "string", "description": "可选风机编号。"},
            "startTime": {"type": "string", "description": "报告时间范围开始。"},
            "endTime": {"type": "string", "description": "报告时间范围结束。"},
            "evidenceJson": {"type": "string", "description": "已有 evidence JSON；没有时可传 {}。"},
        },
        ["farmCode"],
    ),
    _function(
        "create_maintenance_ticket_draft",
        "创建维修工单草稿。只有用户明确要求创建工单/维修单/派工草稿时调用；执行前需要人工确认；该工具只保存草稿，不会真实派单。",
        {
            "farmCode": {"type": "string", "enum": ["FY", "YS"], "description": "风场编码。"},
            "towerCode": {"type": "string", "description": "可选风机编号。"},
            "alarmCode": {"type": "string", "description": "可选告警编码。"},
            "priority": {"type": "string", "description": "优先级，默认 normal。"},
            "evidenceJson": {"type": "string", "description": "支撑工单的 evidence JSON。"},
        },
        ["farmCode"],
    ),
    _function(
        "request_clarification",
        "缺少继续分析所必需的参数时，向用户提出一个简短、具体、一次只问一个重点的问题。",
        {
            "question": {"type": "string", "description": "要问用户的澄清问题。"},
            "reason": {"type": "string", "description": "内部原因，例如 missing_farmCode、missing_deviceTypeCode。"},
        },
        ["question"],
    ),
    _function(
        "finish_answer",
        "已有足够信息可以生成最终回答时调用。不要在仍缺少必需 operational evidence 或 SOP citation 时调用。",
        {
            "evidenceRequirement": {"type": "string", "enum": ["none", "sop", "operational"], "description": "none=通用解释；sop=回答依赖SOPcitation；operational=回答依赖风机/告警/趋势/草稿 evidence。"},
            "answerFocus": {"type": "string", "description": "最终回答关注点，例如 告警原因、处置建议、趋势结论。"},
        },
        ["evidenceRequirement"],
    ),
]


PLANNER_SYSTEM_PROMPT = """你是风机混塔智能运维 AI Copilot 的工具规划器。
你只能选择提供的一个工具或控制动作，不得自行编造数据库、告警、测点或 SOP 事实。
事实查询、RAG 检索和草稿保存必须调用工具。优先使用已有工具结果，避免重复调用。
缺少继续分析所需的关键参数时调用 request_clarification。已有足够证据时及时选择 finish_answer，
避免重复调用同一工具，草稿参数齐全后选择草稿工具，不要继续无关检索。
证据足够形成回答时调用 finish_answer，并正确选择 evidenceRequirement：
- none：通用解释，不依赖运维事实；
- sop：回答引用了维护 SOP，必须已有 citation；
- operational：回答涉及具体风机、告警、趋势或草稿，必须已有 Go evidence。
三个草稿工具只生成草稿，不会真实派单；它们会在执行前进入人工确认。"""
