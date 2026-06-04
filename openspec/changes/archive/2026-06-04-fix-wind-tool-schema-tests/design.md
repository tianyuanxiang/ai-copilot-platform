## 上下文

`service/ai/engine/tests/test_wind_tool_rpc.py` 目前测试的是 `WindToolRPCClient`（gRPC 客户端适配器），而非 `wind_agent_tools.py` 中定义的 8 个工具 Schema。BOSS 明确指出需要测试的是工具 Schema 本身，而不是 gRPC 通信细节。

`wind_agent_tools.py` 导出了以下关键数据结构：
- `READ_ONLY_TOOLS`: 5 个只读工具集合
- `DRAFT_TOOLS`: 3 个草稿工具集合
- `GO_TOOLS`: 上述 8 个工具的并集
- `AGENT_TOOL_SCHEMAS`: 8 个工具 + 2 个控制动作的完整 OpenAI function schema 定义

## 目标 / 非目标

**目标：**
- 重写测试文件，聚焦于 8 个工具 Schema 的正确性验证
- 所有测试函数和注释使用中文
- 验证每个工具的 Schema 结构（名称、描述、参数定义、必填字段）
- 验证工具分类归属（READ_ONLY_TOOLS vs DRAFT_TOOLS vs CONTROL_ACTIONS）
- 验证 GO_TOOLS 并集的正确性

**非目标：**
- 不测试 gRPC 通信层（WindToolRPCClient 的测试保留在原有测试逻辑中，但不在此次重写范围）
- 不进行端到端的工具调用测试（那是集成测试的范畴）
- 不修改 wind_agent_tools.py 中的工具定义

## 决策

1. **测试组织方式**：按工具分组，每个工具一个测试类或测试函数组，包含 Schema 结构验证、参数验证、必填字段验证
2. **Schema 验证策略**：直接导入 `AGENT_TOOL_SCHEMAS`、`READ_ONLY_TOOLS`、`DRAFT_TOOLS`、`GO_TOOLS` 进行白盒验证，无需 Mock
3. **中文注释**：所有测试函数文档字符串、代码块注释均使用中文，与团队规范保持一致
4. **测试数据**：使用真实合理的风场数据（扶余风场 FY、榆树风场 YS 等），提高测试可读性

## 风险 / 权衡

- **风险**：删除原有 gRPC 客户端测试可能导致通信层回归问题 → **缓解**：原有测试覆盖的 RPC 通信逻辑属于基础设施层，若需要可在单独的集成测试中覆盖
- **权衡**：Schema 测试属于静态验证，无法捕获运行时行为问题 → 这是预期行为，运行时问题由集成测试和 E2E 测试覆盖
