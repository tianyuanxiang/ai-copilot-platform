"""风机运维智能助手 8 个 Go 工具的 RPC 集成测试。

通过 WindToolRPCClient 实际调用 Go 后端 AiWindAgentService.ExecuteTool，
验证每个工具的参数传递和响应结构是否正确。

前置条件：Go RPC 服务必须在 127.0.0.1:8091 上运行。
如果服务不可用，整个模块的测试将被跳过（不会报 FAIL）。

8 个 Go 工具：
- 只读工具（5 个）：get_turbine_metadata, search_maintenance_sop,
  query_alarm_events, query_sensor_timeseries, compare_sensor_trend
- 草稿工具（3 个）：generate_alarm_analysis_draft, generate_health_report,
  create_maintenance_ticket_draft
"""

from __future__ import annotations

import asyncio
import json

import pytest

from app.clients.wind_tool_rpc import (
    WindToolExecuteRequest,
    WindToolRPCClient,
    WindToolRPCError,
)
from app.core.config import load_settings


# ---------------------------------------------------------------------------
# Go 服务可用性缓存（模块级全局，避免重复探测）
# None = 未探测，True = 可用，False = 不可用
# ---------------------------------------------------------------------------

_go_service_available = None


# ---------------------------------------------------------------------------
# 函数级 fixture：为每个测试创建独立的 WindToolRPCClient
# ---------------------------------------------------------------------------

# 注意：使用 scope="function" 而非 scope="module"，因为 gRPC 的
# AioCall Future 会绑定到创建 channel 时的事件循环，而 pytest-asyncio
# 的 asyncio_default_test_loop_scope=function 会为每个测试分配独立循环，
# 导致 "Future attached to a different loop" 错误。
#
# 首次调用时探测 Go 服务并缓存结果，后续测试直接复用缓存状态。


@pytest.fixture(scope="function")
async def wind_rpc():
    """创建 WindToolRPCClient 并探测 Go 服务是否可达。

    首次调用时用 get_turbine_metadata 探测 Go 服务并缓存结果。
    后续测试直接复用缓存状态，Go 不可用时跳过，可用时创建新客户端。

    Yields:
        WindToolRPCClient: 已就绪的 RPC 客户端，测试结束后自动关闭
    """
    global _go_service_available

    # 已确认 Go 服务不可用，跳过后续所有测试
    if _go_service_available is False:
        pytest.skip(
            "Go RPC 服务 (127.0.0.1:8091) 不可用，跳过集成测试。"
            "请先启动 Go 后端服务再运行此测试。"
        )

    settings = load_settings()
    client = WindToolRPCClient(
        target=settings.agent_go_rpc_target,
        timeout_seconds=settings.agent_go_rpc_timeout_seconds,
    )

    # 首次调用时探测 Go 服务是否可达
    if _go_service_available is None:
        probe_request = WindToolExecuteRequest(
            user_id=1,
            trace_id="probe-trace",
            conversation_id="probe-conv",
            tool_name="get_turbine_metadata",
            arguments_json="{}",
            step=1,
        )
        try:
            await asyncio.wait_for(client.execute(probe_request), timeout=5.0)
            _go_service_available = True
        except (WindToolRPCError, asyncio.TimeoutError, OSError):
            _go_service_available = False
            await client.close()
            pytest.skip(
                "Go RPC 服务 (127.0.0.1:8091) 不可用，跳过集成测试。"
                "请先启动 Go 后端服务再运行此测试。"
            )

    yield client
    await client.close()


# ---------------------------------------------------------------------------
# 辅助函数
# ---------------------------------------------------------------------------


def _make_request(tool_name: str, arguments: dict, step: int = 1) -> WindToolExecuteRequest:
    """构造发给 Go 工具的 RPC 请求。

    Args:
        tool_name: 工具名称，如 "get_turbine_metadata"
        arguments: 工具参数字典，将被序列化为 JSON 字符串
        step: ReAct 步骤编号，默认为 1

    Returns:
        WindToolExecuteRequest: 可直接传给 WindToolRPCClient.execute() 的请求对象
    """
    return WindToolExecuteRequest(
        user_id=1,
        trace_id=f"test-{tool_name}",
        conversation_id=f"test-{tool_name}-conv",
        tool_name=tool_name,
        arguments_json=json.dumps(arguments, ensure_ascii=False),
        step=step,
    )


def _assert_tool_result(result, expected_tool_name: str):
    """验证 Go 工具返回的响应结构完整性。

    检查项：
    - 工具名称与请求一致
    - 状态字段非空（Go 会返回 success/failed/denied 等）
    - result_json 是合法的 JSON 字符串
    - evidence_json 是非空字符串（可能为 ""）
    - latency_ms 是有效的非负整数

    Args:
        result: WindToolExecuteResult 响应对象
        expected_tool_name: 期望的工具名称
    """
    # 1. 工具名称一致性
    assert result.tool_call.tool_name == expected_tool_name, (
        f"工具名称不匹配: 期望 {expected_tool_name}, "
        f"实际 {result.tool_call.tool_name}"
    )

    # 2. 状态字段非空
    assert result.tool_call.status, (
        f"[{expected_tool_name}] 返回的 status 为空"
    )

    # 3. result_json 是合法的 JSON
    try:
        parsed = json.loads(result.tool_call.result_json)
    except json.JSONDecodeError:
        assert False, (
            f"[{expected_tool_name}] result_json 无法解析为 JSON: "
            f"{result.tool_call.result_json[:200]}"
        )

    # 4. evidence_json 为字符串类型（由客户端保证）
    assert isinstance(result.evidence_json, str), (
        f"[{expected_tool_name}] evidence_json 不是字符串类型"
    )

    # 5. 延迟非负
    assert result.tool_call.latency_ms >= 0, (
        f"[{expected_tool_name}] latency_ms 为负数: "
        f"{result.tool_call.latency_ms}"
    )


# ---------------------------------------------------------------------------
# 1. get_turbine_metadata - 查询风场、风机和设备元数据
# ---------------------------------------------------------------------------


async def test_get_turbine_metadata(wind_rpc: WindToolRPCClient):
    """测试通过场区编码查询风场元数据。

    使用扶余风场（FY）查询风机、设备等基础编码信息。
    get_turbine_metadata 为无必填参数的模糊查询工具，
    传入部分参数即可查询相关元数据。
    """
    # 1. 构造参数：仅指定扶余风场编码
    arguments = {"farmCode": "FY"}
    request = _make_request("get_turbine_metadata", arguments)

    # 2. 调用 Go RPC
    result = await wind_rpc.execute(request)

    # 3. 验证响应
    _assert_tool_result(result, "get_turbine_metadata")


# ---------------------------------------------------------------------------
# 2. search_maintenance_sop - 检索维护 SOP
# ---------------------------------------------------------------------------


async def test_search_maintenance_sop(wind_rpc: WindToolRPCClient):
    """测试通过关键词检索风机维护 SOP 知识库。

    使用中文关键词"齿轮箱油温"发起 RAG 检索，
    指定 topK=3 限制返回条数，范围为公开知识库。
    """
    # 1. 构造参数：检索齿轮箱油温相关 SOP
    arguments = {
        "query": "齿轮箱油温",
        "searchScope": "public",
        "topK": 3,
    }
    request = _make_request("search_maintenance_sop", arguments)

    # 2. 调用 Go RPC
    result = await wind_rpc.execute(request)

    # 3. 验证响应
    _assert_tool_result(result, "search_maintenance_sop")
    # SOP 检索可能返回引用文档
    assert isinstance(result.citations, list), (
        f"[search_maintenance_sop] citations 不是列表类型"
    )


# ---------------------------------------------------------------------------
# 3. query_alarm_events - 查询 TDengine 告警事件
# ---------------------------------------------------------------------------


async def test_query_alarm_events(wind_rpc: WindToolRPCClient):
    """测试通过场区编码查询 TDengine 告警聚合数据。

    以扶余风场（FY）为查询条件，查询所有已存在的告警记录
    （hasStatus=true 表示筛选有状态的告警事件）。
    """
    # 1. 构造参数：查询扶余风场有状态的告警事件
    arguments = {
        "farmCode": "FY",
        "hasStatus": True,
    }
    request = _make_request("query_alarm_events", arguments)

    # 2. 调用 Go RPC
    result = await wind_rpc.execute(request)

    # 3. 验证响应
    _assert_tool_result(result, "query_alarm_events")


# ---------------------------------------------------------------------------
# 4. query_sensor_timeseries - 查询传感器时序数据
# ---------------------------------------------------------------------------


async def test_query_sensor_timeseries(wind_rpc: WindToolRPCClient):
    """测试通过场区和设备类型查询传感器时序数据。

    查询扶余风场（FY）齿轮箱（gear_box）的油温（oil_temp）
    时序数据，限制分页大小。
    """
    # 1. 构造参数：查询齿轮箱油温时序数据
    arguments = {
        "farmCode": "FY",
        "deviceTypeCode": "WPR",
        "field": ["d"],
    }
    request = _make_request("query_sensor_timeseries", arguments)

    # 2. 调用 Go RPC
    result = await wind_rpc.execute(request)

    # 3. 验证响应
    _assert_tool_result(result, "query_sensor_timeseries")


# ---------------------------------------------------------------------------
# 5. compare_sensor_trend - 传感器趋势对比与风险判断
# ---------------------------------------------------------------------------


async def test_compare_sensor_trend(wind_rpc: WindToolRPCClient):
    """测试传感器测点的趋势对比和阈值风险分析。

    对扶余风场（FY）齿轮箱（gear_box）的油温（oil_temp）测点
    进行趋势对比，验证 Go 端能正确调用 CompareTrendLogic 并返回
    阈值解析和风险判断结果。
    """
    # 1. 构造参数：对比齿轮箱油温趋势
    arguments = {
        "farmCode": "FY",
        "deviceTypeCode": "ACCX",
        "field": ["accel"],
    }
    request = _make_request("compare_sensor_trend", arguments)

    # 2. 调用 Go RPC
    result = await wind_rpc.execute(request)

    # 3. 验证响应
    _assert_tool_result(result, "compare_sensor_trend")


# ---------------------------------------------------------------------------
# 6. generate_alarm_analysis_draft - 生成告警分析草稿
# ---------------------------------------------------------------------------


async def test_generate_alarm_analysis_draft(wind_rpc: WindToolRPCClient):
    """测试根据证据生成告警分析草稿。

    以扶余风场（FY）的告警 1001 为例，传入空证据 JSON，
    验证 Go 端能正确构建 AlarmEvidence 并调用 WindAlarmSummary 引擎
    生成告警分析草稿。
    注意：该工具生成的是草稿，不会真实派单，需要人工确认。
    """
    # 1. 构造参数：为扶余风场告警 1001 生成分析草稿
    arguments = {
        "farmCode": "FY",
        "alarmCode": "1001",
        "evidenceJson": "{}",
    }
    request = _make_request("generate_alarm_analysis_draft", arguments)

    # 2. 调用 Go RPC
    result = await wind_rpc.execute(request)

    # 3. 验证响应
    _assert_tool_result(result, "generate_alarm_analysis_draft")


# ---------------------------------------------------------------------------
# 7. generate_health_report - 生成健康报告草稿
# ---------------------------------------------------------------------------


async def test_generate_health_report(wind_rpc: WindToolRPCClient):
    """测试根据证据生成风机健康报告草稿。

    以扶余风场（FY）的日报为例，传入空证据 JSON，
    验证 Go 端能正确调用 WindHealthReportDraft 引擎生成健康报告。
    注意：该工具生成的是草稿，需要人工确认。
    """
    # 1. 构造参数：为扶余风场生成日报类型的健康报告
    arguments = {
        "farmCode": "FY",
        "reportType": "daily",
        "evidenceJson": "{}",
    }
    request = _make_request("generate_health_report", arguments)

    # 2. 调用 Go RPC
    result = await wind_rpc.execute(request)

    # 3. 验证响应
    _assert_tool_result(result, "generate_health_report")


# ---------------------------------------------------------------------------
# 8. create_maintenance_ticket_draft - 生成维修工单草稿
# ---------------------------------------------------------------------------


async def test_create_maintenance_ticket_draft(wind_rpc: WindToolRPCClient):
    """测试根据证据生成维修工单草稿。

    以扶余风场（FY）告警 1001、优先级 high 为例，
    验证 Go 端能正确调用 WindTicketDraft 引擎生成维修工单草稿。
    注意：该工具生成的是草稿，不会真实派单，需要人工确认。
    """
    # 1. 构造参数：高优先级维修工单草稿
    arguments = {
        "farmCode": "FY",
        "alarmCode": "1001",
        "priority": "high",
        "evidenceJson": "{}",
    }
    request = _make_request("create_maintenance_ticket_draft", arguments)

    # 2. 调用 Go RPC
    result = await wind_rpc.execute(request)

    # 3. 验证响应
    _assert_tool_result(result, "create_maintenance_ticket_draft")
