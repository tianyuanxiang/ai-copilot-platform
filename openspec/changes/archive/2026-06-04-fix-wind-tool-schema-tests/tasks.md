## 1. 准备工作

- [x] 1.1 阅读 `wind_agent_tools.py` 中 8 个工具 Schema 的完整定义，确认每个工具的名称、参数、必填字段
- [x] 1.2 分析现有 `test_wind_tool_rpc.py`，确认哪些 import 和辅助函数可复用

## 2. 工具 Schema 结构验证

- [x] 2.1 编写测试：验证 `AGENT_TOOL_SCHEMAS` 中每个工具包含 `type`、`function` 字段
- [x] 2.2 编写测试：验证 8 个 Go 工具的名称在集合中唯一无重复
- [x] 2.3 编写测试：验证每个工具的 `parameters` 字段符合 JSON Schema 规范（含 `type`、`properties`、`required`）

## 3. 工具分类归属验证

- [x] 3.1 编写测试：验证 `READ_ONLY_TOOLS` 包含正确的 5 个工具名称且大小为 5
- [x] 3.2 编写测试：验证 `DRAFT_TOOLS` 包含正确的 3 个工具名称且大小为 3
- [x] 3.3 编写测试：验证 `GO_TOOLS` 等于 `READ_ONLY_TOOLS | DRAFT_TOOLS` 且大小为 8
- [x] 3.4 编写测试：验证 `CONTROL_ACTIONS` 与 `GO_TOOLS` 无交集

## 4. 工具参数定义验证

- [x] 4.1 编写测试：验证 `get_turbine_metadata` 的必填参数为空（模糊查询，参数全部可选）
- [x] 4.2 编写测试：验证 `search_maintenance_sop` 的必填参数仅含 `query`
- [x] 4.3 编写测试：验证 `query_alarm_events` 的必填参数含 `farmCode`
- [x] 4.4 编写测试：验证 `query_sensor_timeseries` 的必填参数含 `farmCode` 和 `deviceTypeCode`
- [x] 4.5 编写测试：验证 `compare_sensor_trend` 的必填参数含 `farmCode` 和 `deviceTypeCode`
- [x] 4.6 编写测试：验证 `generate_alarm_analysis_draft` 的必填参数含 `farmCode`
- [x] 4.7 编写测试：验证 `generate_health_report` 的必填参数含 `farmCode`
- [x] 4.8 编写测试：验证 `create_maintenance_ticket_draft` 的必填参数含 `farmCode`

## 5. 工具描述验证

- [x] 5.1 编写测试：验证 8 个 Go 工具的 `description` 字段均为非空中文字符串

## 6. 验证运行

- [x] 6.1 运行 `pytest test_wind_tool_rpc.py -v` 确认所有测试通过
- [x] 6.2 确认测试输出中每个测试名称清晰反映其验证内容
