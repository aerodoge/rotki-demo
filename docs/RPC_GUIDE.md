# RPC 节点完整指南

## 目录

1. [概述](#概述)
2. [快速启动](#快速启动)
3. [功能特性](#功能特性)
4. [使用指南](#使用指南)
5. [配置示例](#配置示例)
6. [负载均衡策略](#负载均衡策略)
7. [API 参考](#api-参考)
8. [监控和维护](#监控和维护)
9. [故障排查](#故障排查)
10. [最佳实践](#最佳实践)

## 概述

RPC 节点功能允许你为每个区块链网络配置多个 RPC 端点。这提供了：

- **负载均衡**：根据权重在多个节点之间分配请求
- **故障转移**：如果主节点失败，自动切换到备用节点
- **灵活性**：轻松从 DeBank API 切换到你自己的 RPC 节点
- **连接监控**：自动健康检查所有配置的节点
- **成本优化**：混合使用免费和付费节点

## 快速启动

### 1. 运行数据库迁移

应用 RPC 节点表迁移：

```bash
# 选项 1：使用 MySQL 命令行
mysql -u root -p rotki_demo < migrations/002_add_rpc_nodes_table.sql

# 选项 2：使用 GORM 自动迁移（在 main.go 中取消注释）
# 编辑 cmd/server/main.go 并取消注释 AutoMigrate 部分
```

### 2. 启动后端

```bash
# 构建并运行
go run cmd/server/main.go

# 或使用 make
make run
```

API 将在 `http://localhost:8080` 可用

### 3. 启动前端

```bash
cd frontend
npm install  # 如果尚未完成
npm run dev
```

UI 将在 `http://localhost:3000` 可用

### 4. 访问 RPC 节点设置

1. 打开 http://localhost:3000
2. 点击侧边栏菜单
3. 展开**设置**
4. 点击 **RPC 节点**

## 功能特性

### 后端实现

1. **数据库模型**（`internal/models/models.go`）
   - `RPCNode` 模型，包含字段：
     - `chain_id`：哪个区块链（eth、bsc、polygon 等）
     - `name`：显示名称（例如，"0xRPC"、"PublicNode"）
     - `url`：RPC 端点 URL
     - `weight`：负载均衡权重（0-100）
     - `is_enabled`：启用/禁用切换
     - `is_connected`：连接状态（自动更新）
     - `priority`：优先级更高的节点优先
     - `timeout`：请求超时（秒）

2. **仓储层**（`internal/repository/rpc_node_repository.go`）
   - CRUD 操作
   - 按链查询
   - 获取启用的节点
   - 按链分组
   - 更新连接状态

3. **服务层**（`internal/service/rpc_node_service.go`）
   - 连接测试（通过 JSON-RPC 调用 `eth_blockNumber`）
   - 创建时自动检查连接
   - 批量连接检查
   - 节点管理的业务逻辑

4. **API 端点**
   ```
   POST   /api/v1/rpc-nodes              创建新的 RPC 节点
   GET    /api/v1/rpc-nodes              列出所有节点（按 chain_id 过滤）
   GET    /api/v1/rpc-nodes/grouped      获取按链分组的节点
   GET    /api/v1/rpc-nodes/:id          获取单个节点
   PUT    /api/v1/rpc-nodes/:id          更新节点
   DELETE /api/v1/rpc-nodes/:id          删除节点
   POST   /api/v1/rpc-nodes/:id/check    检查特定节点的连接
   POST   /api/v1/rpc-nodes/check-all    检查所有节点
   ```

### 前端实现

1. **设置页面**（`frontend/src/views/RPCNodesSettings.vue`）
   - 类似 Rotki 的 RPC 节点设置的 UI
   - 链标签页，便于在网络之间导航
   - 表格视图，显示节点信息、权重和连接状态
   - 添加/编辑/删除节点模态框
   - 启用/禁用切换开关
   - 连接状态徽章

2. **功能**
   - **链标签页**：在不同区块链之间切换
   - **节点权重**：配置负载均衡的百分比（0-100%）
   - **连接性**：实时连接状态显示
   - **启用/禁用**：每个节点的切换开关
   - **编辑**：使用模态框进行内联编辑
   - **删除**：带确认的删除节点

## 使用指南

### 访问 RPC 节点设置

导航到：侧边栏中的**设置 → RPC 节点**

### 添加新节点

点击**"+ Add Node"**按钮：
- **Chain**：选择区块链网络
- **Node Name**：给它一个描述性的名称
- **RPC URL**：输入端点 URL
- **Weight**：设置负载均衡权重（0-100）
- **Priority**：可选的优先级级别
- **Timeout**：请求超时（默认：30 秒）
- **Enable**：选中以立即激活

### 配置多个节点

以太坊的示例配置：

```
Node 1: 0xRPC
URL: https://0xrpc.io/eth
Weight: 50%
Priority: 1
Status: CONNECTED

Node 2: PublicNode
URL: https://ethereum-rpc.publicnode.com
Weight: 30%
Priority: 0
Status: CONNECTED

Node 3: Alchemy Backup
URL: https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY
Weight: 20%
Priority: 2
Status: CONNECTED
```

### 基于权重的负载均衡

系统将根据以下内容路由请求：
1. **优先级**：首先尝试优先级更高的节点
2. **权重**：在相同优先级的节点中，权重决定分配
3. **状态**：仅使用已启用和已连接的节点

示例：使用上述配置：
- 50% 的请求发送到 0xRPC
- 30% 发送到 PublicNode
- 20% 发送到 Alchemy Backup

### 连接监控

- 创建时测试节点
- 通过"Check Connection"按钮进行手动检查
- 可以安排自动后台检查
- 连接状态实时更新

## 配置示例

### 以太坊主网

#### 免费公共节点

```json
{
  "chain_id": "eth",
  "name": "LlamaNodes",
  "url": "https://eth.llamarpc.com",
  "weight": 40,
  "priority": 1,
  "is_enabled": true
}
```

```json
{
  "chain_id": "eth",
  "name": "PublicNode",
  "url": "https://ethereum-rpc.publicnode.com",
  "weight": 30,
  "priority": 0,
  "is_enabled": true
}
```

```json
{
  "chain_id": "eth",
  "name": "1RPC",
  "url": "https://1rpc.io/eth",
  "weight": 30,
  "priority": 0,
  "is_enabled": true
}
```

#### 付费节点（更高优先级）

```json
{
  "chain_id": "eth",
  "name": "Alchemy",
  "url": "https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY",
  "weight": 100,
  "priority": 10,
  "is_enabled": true
}
```

```json
{
  "chain_id": "eth",
  "name": "Infura",
  "url": "https://mainnet.infura.io/v3/YOUR_PROJECT_ID",
  "weight": 100,
  "priority": 10,
  "is_enabled": true
}
```

### 其他链配置

#### BSC（币安智能链）

```json
{
  "chain_id": "bsc",
  "name": "BSC RPC",
  "url": "https://bsc-dataseed1.binance.org",
  "weight": 50,
  "priority": 1,
  "is_enabled": true
}
```

#### Polygon

```json
{
  "chain_id": "matic",
  "name": "Polygon RPC",
  "url": "https://polygon-rpc.com",
  "weight": 100,
  "priority": 1,
  "is_enabled": true
}
```

#### Arbitrum

```json
{
  "chain_id": "arb",
  "name": "Arbitrum One",
  "url": "https://arb1.arbitrum.io/rpc",
  "weight": 100,
  "priority": 1,
  "is_enabled": true
}
```

#### Optimism

```json
{
  "chain_id": "op",
  "name": "Optimism",
  "url": "https://mainnet.optimism.io",
  "weight": 100,
  "priority": 1,
  "is_enabled": true
}
```

## 负载均衡策略

### 优先级优先

节点选择基于：
1. **优先级**（较高的优先）
2. **权重**（相同优先级内的百分比）
3. **连接状态**（仅连接的节点）

### 示例配置

```
优先级 10（高级节点）：
- Alchemy：100% 权重
- Infura：100% 权重
→ 50% 的请求发送到 Alchemy，50% 发送到 Infura

优先级 1（免费第一层）：
- LlamaNodes：50% 权重
- 1RPC：50% 权重
→ 仅在优先级 10 节点失败时使用

优先级 0（回退）：
- PublicNode：100% 权重
→ 仅在所有更高优先级节点失败时使用
```

## API 参考

### 创建节点

```bash
curl -X POST http://localhost:8080/api/v1/rpc-nodes \
  -H "Content-Type: application/json" \
  -d '{
    "chain_id": "eth",
    "name": "0xRPC",
    "url": "https://0xrpc.io/eth",
    "weight": 100,
    "priority": 1,
    "timeout": 30,
    "is_enabled": true
  }'
```

### 列出链的节点

```bash
curl http://localhost:8080/api/v1/rpc-nodes?chain_id=eth
```

### 获取分组的节点

```bash
curl http://localhost:8080/api/v1/rpc-nodes/grouped
```

### 检查连接

```bash
curl -X POST http://localhost:8080/api/v1/rpc-nodes/1/check
```

### 更新节点

```bash
curl -X PUT http://localhost:8080/api/v1/rpc-nodes/1 \
  -H "Content-Type: application/json" \
  -d '{
    "chain_id": "eth",
    "name": "0xRPC Updated",
    "url": "https://0xrpc.io/eth",
    "weight": 80,
    "priority": 2,
    "timeout": 30,
    "is_enabled": true
  }'
```

### 删除节点

```bash
curl -X DELETE http://localhost:8080/api/v1/rpc-nodes/1
```

## 监控和维护

### 检查节点状态

通过 API：
```bash
curl http://localhost:8080/api/v1/rpc-nodes/grouped | jq '.'
```

通过 UI：
1. 转到设置 → RPC 节点
2. 查看"连接性"列
3. 绿色徽章 = 已连接
4. 红色徽章 = 已断开连接

### 手动连接检查

在 UI 中点击"Check Connection"按钮，或：

```bash
curl -X POST http://localhost:8080/api/v1/rpc-nodes/1/check
```

### 禁用有问题的节点

在 UI 中：
1. 切换开关以禁用
2. 节点将在路由中被跳过

通过 API：
```bash
curl -X PUT http://localhost:8080/api/v1/rpc-nodes/1 \
  -H "Content-Type: application/json" \
  -d '{"is_enabled": false, ...其他字段...}'
```

## 故障排查

### 节点显示为断开连接

1. 检查 URL 是否正确
2. 验证网络连接
3. 手动测试：
   ```bash
   curl -X POST https://eth-rpc.example.com \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
   ```
4. 检查是否需要 API 密钥
5. 增加超时值

### 无法创建节点

1. 确保链表中有 chain_id
2. 检查外键约束
3. 验证链是否存在：
   ```bash
   curl http://localhost:8080/api/v1/chains
   ```

### 前端未显示节点

1. 检查后端是否正在运行
2. 打开浏览器控制台查看错误
3. 验证 API 端点：
   ```bash
   curl http://localhost:8080/api/v1/rpc-nodes/grouped
   ```
4. 检查 CORS 设置

### 权重不按预期工作

- 确保每条链的总权重不超过 100%
- 检查节点是否已启用
- 验证连接状态
- 查看优先级设置

### 连接检查失败

- 增加超时值
- 检查 RPC 端点是否支持 `eth_blockNumber`
- 对于非 EVM 链，修改连接测试方法
- 如果需要，检查 API 密钥

## 最佳实践

### 1. 每条链使用多个节点

- 最少 2 个节点以实现冗余
- 混合免费和付费以优化成本
- 不同提供商以实现多样性

### 2. 设置适当的超时

- 公共节点：30-60 秒
- 付费节点：10-30 秒
- 高流量：5-10 秒

### 3. 定期健康检查

- 每天运行 check-all
- 监控连接状态
- 替换失败的节点

### 4. 权重分配

- 高级节点：100% 权重，高优先级
- 免费可靠：50-100% 权重，中等优先级
- 实验性：20-50% 权重，低优先级

### 5. 从 DeBank 逐步迁移

- 第 1 周：配置节点，使用 10% 流量测试
- 第 2 周：增加到 50% 流量
- 第 3 周：监控稳定性，增加到 80%
- 第 4 周：完全迁移，保留 DeBank 作为回退

## 架构优势

### 1. 提供商抽象

系统维护现有的提供商接口：

```go
type DataProvider interface {
    GetTokenList(ctx context.Context, address string, chainIDs []string) ([]TokenInfo, error)
    GetTotalBalance(ctx context.Context, address string) (float64, error)
    GetUsedChains(ctx context.Context, address string) ([]string, error)
}
```

这意味着：
- ✅ 易于从 DeBank 切换到 RPC 节点
- ✅ 可以混合 DeBank 和 RPC 节点
- ✅ 可使用模拟进行测试
- ✅ 不需要更改现有代码

### 2. 渐进式迁移

你可以：
1. 配置 RPC 节点与 DeBank 一起使用
2. 使用特定地址测试 RPC 节点
3. 逐渐增加 RPC 节点权重
4. 最终移除 DeBank 依赖

### 3. 成本优化

- 用于测试的免费公共节点
- 用于生产的付费节点
- 根据负载混合免费和付费
- 实时监控成本

## 未来增强

### 立即后续步骤

1. **与数据提供商集成**
   - 修改 `internal/provider/provider.go` 以使用 RPC 节点
   - 实现加权轮询选择
   - 添加故障转移逻辑

2. **自动健康检查**
   - 后台 goroutine 定期检查所有节点
   - 自动更新连接状态
   - 节点故障时通知

3. **每个节点的速率限制**
   - 配置每个 RPC 节点的速率限制
   - 跟踪请求计数
   - 防止配额耗尽

### 高级功能

1. **响应时间跟踪**
   - 测量每个节点的延迟
   - 根据性能自动调整权重
   - 在 UI 中显示平均响应时间

2. **成本跟踪**
   - 跟踪每个节点的 API 调用
   - 计算成本（对于付费端点）
   - 预算警报

3. **负载测试**
   - 测试节点容量
   - 识别最快的节点
   - 优化权重分配

4. **地理分布**
   - 按区域标记节点
   - 根据用户位置路由
   - 延迟优化

5. **智能合约调用**
   - 使用 RPC 节点进行合约交互
   - 批量调用以提高效率
   - 对于复杂查询，回退到 DeBank

## 总结

RPC 节点功能为以下内容提供了完整的基础设施：
- ✅ 管理每条链的多个 RPC 端点
- ✅ 带可配置权重的负载均衡
- ✅ 连接监控和健康检查
- ✅ 轻松从 DeBank 迁移到自托管节点
- ✅ 通过混合节点策略进行成本优化
- ✅ 匹配 Rotki 设计的生产就绪 UI

此功能为完全独立的区块链数据检索奠定了基础，消除了第三方 API 依赖，同时保持灵活性和可靠性。

## 支持

如有问题或疑问：
1. 检查日志：`logs/app.log`
2. 查看架构文档：`docs/ARCHITECTURE.md`
3. 使用脚本测试：`scripts/test_rpc_nodes.sh`
4. 验证数据库：检查 `rpc_nodes` 表
