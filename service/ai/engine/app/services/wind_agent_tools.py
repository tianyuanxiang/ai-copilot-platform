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
        "查询风场、风机和设备元数据。缺少编码时优先使用它补齐事实。",
        {
            "farmCode": {"type": "string"},
            "towerCode": {"type": "string"},
            "deviceTypeCode": {"type": "string"},
        },
    ),
    _function(
        "search_maintenance_sop",
        "检索维护 SOP。该工具通过 Go SearchKnowledge 执行权限受控 RAG。",
        {
            "query": {"type": "string"},
            "kbId": {"type": "integer"},
            "searchScope": {"type": "string", "enum": ["public", "personal", "all"]},
            "domainId": {"type": "integer"},
            "documentIds": {"type": "array", "items": {"type": "integer"}},
            "topK": {"type": "integer"},
            "answerMode": {"type": "string", "enum": ["quick", "deep"]},
        },
        ["query"],
    ),
    _function(
        "query_alarm_events",
        "查询 TDengine 告警聚合、时间桶和抽样 evidence。如果是扶余风场farmCode参数为FY，榆树风场为YS，不存在其他选项",
        {
            "farmCode": {"type": "string"},
            "towerCode": {"type": "string"},
            "alarmCode": {"type": "string"},
            "startTime": {"type": "string"},
            "endTime": {"type": "string"},
            "status": {"type": "integer"},
            "hasStatus": {"type": "boolean"},
        },
        ["farmCode"],
    ),
    _function(
        "query_sensor_timeseries",
        "查询指定设备类型和测点字段的 TDengine 时序数据。",
        {
            "farmCode": {"type": "string"},
            "towerCode": {"type": "string"},
            "deviceCode": {"type": "string"},
            "deviceTypeCode": {"type": "string"},
            "field": {"type": "array", "items": {"type": "string"}, "maxItems": 8},
            "startTime": {"type": "string"},
            "endTime": {"type": "string"},
            "page": {"type": "integer"},
            "pageSize": {"type": "integer"},
            "indexId": {"type": "integer"},
        },
        ["farmCode", "deviceTypeCode"],
    ),
    _function(
        "compare_sensor_trend",
        "对传感器测点做趋势对比、阈值解析和风险判断。",
        {
            "farmCode": {"type": "string"},
            "towerCode": {"type": "string"},
            "deviceCode": {"type": "string"},
            "deviceTypeCode": {"type": "string"},
            "field": {"type": "array", "items": {"type": "string"}, "maxItems": 8},
            "startTime": {"type": "string"},
            "endTime": {"type": "string"},
            "indexId": {"type": "integer"},
        },
        ["farmCode", "deviceTypeCode"],
    ),
    _function(
        "generate_alarm_analysis_draft",
        "根据 evidence 生成并保存告警分析草稿。调用前必须人工确认。",
        {
            "farmCode": {"type": "string"},
            "towerCode": {"type": "string"},
            "alarmCode": {"type": "string"},
            "startTime": {"type": "string"},
            "endTime": {"type": "string"},
            "evidenceJson": {"type": "string"},
        },
        ["farmCode"],
    ),
    _function(
        "generate_health_report",
        "根据 evidence 生成并保存健康报告草稿。调用前必须人工确认。",
        {
            "reportType": {"type": "string"},
            "farmCode": {"type": "string"},
            "towerCode": {"type": "string"},
            "startTime": {"type": "string"},
            "endTime": {"type": "string"},
            "evidenceJson": {"type": "string"},
        },
        ["farmCode"],
    ),
    _function(
        "create_maintenance_ticket_draft",
        "根据 evidence 生成并保存维修工单草稿。调用前必须人工确认。",
        {
            "farmCode": {"type": "string"},
            "towerCode": {"type": "string"},
            "alarmCode": {"type": "string"},
            "priority": {"type": "string"},
            "evidenceJson": {"type": "string"},
        },
        ["farmCode"],
    ),
    _function(
        "request_clarification",
        "信息不足时向用户提出一个简短、具体的问题。",
        {
            "question": {"type": "string"},
            "reason": {"type": "string"},
        },
        ["question"],
    ),
    _function(
        "finish_answer",
        "已有信息足够时结束工具选择并生成最终回答。",
        {
            "evidenceRequirement": {"type": "string", "enum": ["none", "sop", "operational"]},
            "answerFocus": {"type": "string"},
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
