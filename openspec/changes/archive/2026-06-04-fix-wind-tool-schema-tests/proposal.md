## 为什么

现有的 `test_wind_tool_rpc.py` 测试文件只覆盖了 `WindToolRPCClient`（gRPC 客户端适配器）的参数映射、错误包装、通道生命周期等底层通信行为，完全没有测试 `wind_agent_tools.py` 中定义的 8 个实际业务工具。BOSS 需要的是对这 8 个工具 Schema 的验证测试，确保工具定义的正确性（名称、参数、必填字段、分类归属等），而非 gRPC 通信细节。

## 变更内容

- 重写 `test_wind_tool_rpc.py`，将测试焦点从 gRPC 客户端适配器转移到 8 个工具 Schema 的验证
- 所有测试函数和注释使用中文，便于团队理解
- 测试覆盖：工具名称唯一性、Schema 结构完整性、参数定义正确性、必填字段声明、工具分类归属（READ_ONLY_TOOLS vs DRAFT_TOOLS）

## 功能 (Capabilities)

### 新增功能
- `wind-tool-schema-tests`: 对 `wind_agent_tools.py` 中 8 个工具 Schema 的完整验证测试

### 修改功能
<!-- 无 -->

## 影响

- 受影响文件: `service/ai/engine/tests/test_wind_tool_rpc.py`（重写）
- 依赖模块: `service/ai/engine/app/services/wind_agent_tools.py`
- 测试框架: pytest（与项目现有测试框架一致）
