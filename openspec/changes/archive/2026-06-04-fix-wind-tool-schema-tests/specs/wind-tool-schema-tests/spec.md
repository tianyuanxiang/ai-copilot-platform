## 新增需求

### 需求:工具 Schema 结构验证

`wind_agent_tools.py` 中定义的每个工具 Schema 必须符合 OpenAI function calling 规范，包含完整的 name、description、parameters 字段。

#### 场景:每个工具 Schema 包含必要字段
- **当** 遍历 `AGENT_TOOL_SCHEMAS` 中的每个工具定义
- **那么** 每个工具必须包含 `type` 字段且值为 `"function"`
- **且** 每个工具必须包含 `function` 对象，其中包含 `name`、`description`、`parameters` 三个必填字段

#### 场景:工具名称在 8 个工具集合中唯一
- **当** 检查 `AGENT_TOOL_SCHEMAS` 中除控制动作外的工具名称
- **那么** 所有 8 个 Go 工具的名称必须唯一，无重复

#### 场景:工具参数定义为 JSON Schema 格式
- **当** 检查每个工具的 `parameters` 字段
- **那么** `parameters` 必须包含 `type`（值为 `"object"`）、`properties`、`required` 字段
- **且** `properties` 中的每个属性必须含有 `type` 字段

### 需求:工具分类归属正确

8 个 Go 工具的集合归属必须正确：5 个只读工具、3 个草稿工具。

#### 场景:READ_ONLY_TOOLS 包含正确的 5 个工具
- **当** 检查 `READ_ONLY_TOOLS` 集合
- **那么** 必须包含 `get_turbine_metadata`、`search_maintenance_sop`、`query_alarm_events`、`query_sensor_timeseries`、`compare_sensor_trend`
- **且** 集合大小必须为 5

#### 场景:DRAFT_TOOLS 包含正确的 3 个工具
- **当** 检查 `DRAFT_TOOLS` 集合
- **那么** 必须包含 `generate_alarm_analysis_draft`、`generate_health_report`、`create_maintenance_ticket_draft`
- **且** 集合大小必须为 3

#### 场景:GO_TOOLS 为 READ_ONLY_TOOLS 与 DRAFT_TOOLS 的并集
- **当** 计算 `READ_ONLY_TOOLS | DRAFT_TOOLS`
- **那么** 结果必须与 `GO_TOOLS` 完全一致
- **且** `GO_TOOLS` 大小必须为 8

#### 场景:CONTROL_ACTIONS 与 GO_TOOLS 无交集
- **当** 检查 `CONTROL_ACTIONS` 与 `GO_TOOLS` 的交集
- **那么** 交集必须为空

### 需求:工具参数定义正确

每个工具的 `required` 字段必须与 Schema 描述中的业务约束一致。

#### 场景:get_turbine_metadata 必填参数
- **当** 检查 `get_turbine_metadata` 的 `required` 字段
- **那么** 必填参数列表必须为空（所有参数可选，用于模糊查询）

#### 场景:search_maintenance_sop 必填参数
- **当** 检查 `search_maintenance_sop` 的 `required` 字段
- **那么** 必须包含 `query`
- **且** 其余参数（kbId、searchScope、domainId、documentIds、topK、answerMode）为可选

#### 场景:query_alarm_events 必填参数
- **当** 检查 `query_alarm_events` 的 `required` 字段
- **那么** 必须包含 `farmCode`

#### 场景:query_sensor_timeseries 必填参数
- **当** 检查 `query_sensor_timeseries` 的 `required` 字段
- **那么** 必须包含 `farmCode` 和 `deviceTypeCode`

#### 场景:compare_sensor_trend 必填参数
- **当** 检查 `compare_sensor_trend` 的 `required` 字段
- **那么** 必须包含 `farmCode` 和 `deviceTypeCode`

#### 场景:generate_alarm_analysis_draft 必填参数
- **当** 检查 `generate_alarm_analysis_draft` 的 `required` 字段
- **那么** 必须包含 `farmCode`

#### 场景:generate_health_report 必填参数
- **当** 检查 `generate_health_report` 的 `required` 字段
- **那么** 必须包含 `farmCode`

#### 场景:create_maintenance_ticket_draft 必填参数
- **当** 检查 `create_maintenance_ticket_draft` 的 `required` 字段
- **那么** 必须包含 `farmCode`

### 需求:工具描述不为空

每个工具必须包含中文描述，清晰地说明工具的用途和使用场景。

#### 场景:所有 Go 工具描述非空
- **当** 遍历 8 个 Go 工具的 Schema
- **那么** 每个工具的 `description` 字段必须非空字符串
- **且** 描述内容必须包含中文字符

### 需求:测试代码使用中文注释

测试文件中的所有注释、文档字符串必须使用中文，便于团队理解和维护。

#### 场景:测试函数文档字符串为中文
- **当** 读取测试文件中任意测试函数的 docstring
- **那么** 文档字符串必须包含中文说明，清晰表达测试目的

#### 场景:代码块注释为中文
- **当** 读取测试文件中的代码注释
- **那么** 注释必须使用中文，描述测试步骤和验证逻辑
