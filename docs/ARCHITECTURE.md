# 架构文档

## 系统概述

Rotki Demo 是一个全栈 DeFi 资产管理应用程序，用于跟踪多个区块链地址的加密货币持仓。它使用 DeBank API 作为数据源，但设计了抽象层，可以轻松切换到其他数据提供商。

## 架构图

```
┌─────────────────────────────────────────────────────────┐
│                     Frontend (Vue.js)                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │   Sidebar    │  │ EVM Accounts │  │   Modals     │   │
│  └──────────────┘  └──────────────┘  └──────────────┘   │
│  ┌─────────────────────────────────────────────────┐    │
│  │           Pinia Store (State Management)        │    │
│  └─────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────┐    │
│  │              Axios (HTTP Client)                │    │
│  └─────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
                           │
                           │ HTTP/JSON
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    Backend (Go/Gin)                     │
│  ┌─────────────────────────────────────────────────┐    │
│  │              HTTP Router (Gin)                  │    │
│  │    /api/v1/wallets  /api/v1/addresses           │    │
│  └─────────────────────────────────────────────────┘    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │   Handlers   │  │   Services   │  │ Repositories │   │
│  └──────────────┘  └──────────────┘  └──────────────┘   │
│  ┌─────────────────────────────────────────────────┐    │
│  │         Data Provider Interface                 │    │
│  │  ┌────────────────┐   ┌─────────────────────┐   │    │
│  │  │ DeBank Client  │   │  Future: Custom RPC │   │    │
│  │  └────────────────┘   └─────────────────────┘   │    │
│  └─────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
           │                           │
           │                           │
           ▼                           ▼
    ┌────────────┐           ┌──────────────────┐
    │   MySQL    │           │   DeBank API     │
    │  Database  │           │ (Rate Limited)   │
    └────────────┘           └──────────────────┘
```

## 后端架构

### 层级分解

#### 1. API 层（`internal/api`）
- **Handler**：请求/响应处理、验证
- **Router**：路由定义和中间件设置
- **职责**：仅处理 HTTP 相关事务

#### 2. 服务层（`internal/service`）
- **SyncService**：后台同步逻辑
- **业务逻辑**：跨多个仓储协调操作
- **事务管理**：确保数据一致性

#### 3. 仓储层（`internal/repository`）
- **WalletRepository**：钱包 CRUD 操作
- **AddressRepository**：地址 CRUD 操作
- **TokenRepository**：代币数据管理
- **纯数据访问**：无业务逻辑

#### 4. 提供商层（`internal/provider`）
- **接口定义**：DataProvider 接口
- **DeBank 实现**：当前数据源
- **可扩展性**：易于添加新提供商

### 核心设计模式

#### 提供商模式
```go
type DataProvider interface {
    GetTotalBalance(ctx context.Context, address string) (*TotalBalanceResponse, error)
    GetTokenList(ctx context.Context, address string, chainIDs []string) ([]TokenInfo, error)
    // ... 其他方法
}
```

**优势：**
- 将数据源与业务逻辑解耦
- 易于在 DeBank、自定义 RPC 或其他 API 之间切换
- 可使用模拟提供商进行测试

#### 仓储模式
```go
type WalletRepository struct {
    db *gorm.DB
}

func (r *WalletRepository) GetByID(id uint) (*models.Wallet, error) {
    // 仅数据库访问
}
```

**优势：**
- 将数据访问与业务逻辑分离
- 更易于测试
- 可以切换数据库实现

#### 服务模式
```go
type SyncService struct {
    dataProvider provider.DataProvider
    addressRepo  *repository.AddressRepository
    tokenRepo    *repository.TokenRepository
}
```

**优势：**
- 封装复杂的业务逻辑
- 协调多个仓储和提供商
- 管理事务和错误处理

## 前端架构

### 组件结构

```
App.vue
├── Sidebar.vue（导航）
└── EVMAccounts.vue（主视图）
    ├── 钱包行（可展开）
    └── 地址行（嵌套）
        └── 代币显示
```

### 状态管理（Pinia）

```javascript
walletStore
├── state
│   ├── wallets[]
│   ├── addresses[]
│   ├── selectedWallet
│   └── loading
├── getters
│   ├── getAddressesByWallet()
│   ├── getTotalValueByWallet()
│   └── getTotalValue()
└── actions
    ├── fetchWallets()
    ├── createWallet()
    ├── refreshWallet()
    └── ...
```

## 数据流

### 添加新地址

1. **用户操作**：点击"Add Address"按钮
2. **前端**：打开模态框，用户填写表单
3. **API 调用**：POST /api/v1/addresses
4. **后端处理器**：验证请求
5. **仓储**：将地址插入数据库
6. **后台触发**：启动新地址的同步
7. **提供商**：调用 DeBank API
8. **数据存储**：将代币存储在数据库中
9. **响应**：将创建的地址返回给前端
10. **UI 更新**：显示带有代币的新地址

### 自动同步流程

```
SyncService.Start()
    │
    ├─> Ticker（每 5 分钟）
    │
    ├─> GetAllNeedingSync()
    │     │
    │     └─> 查询 last_synced_at > interval 的地址
    │
    ├─> ProcessInBatches()
    │     │
    │     └─> 并发处理 10 个地址
    │
    └─> 对于每个地址：
          │
          ├─> GetTokenList() 从提供商获取
          │
          ├─> UpsertBatch() 到数据库
          │
          └─> UpdateLastSynced()
```

## 数据库模式

### 核心表

**wallets**
- 存储钱包元数据
- 与地址一对多关系

**addresses**
- 存储区块链地址
- 链接到钱包
- 跟踪最后同步时间

**tokens**
- 存储代币余额
- 链接到地址和链
- 每次同步时更新

**asset_snapshots**
- 历史快照
- 用于跟踪余额随时间的变化

### 关系

```
wallets (1) ──< (N) addresses (1) ──< (N) tokens
                                  │
                                  └──< (N) asset_snapshots
```

## 速率限制策略

### 令牌桶算法

```go
limiter := rate.NewLimiter(
    rate.Limit(5),  // 每秒 5 个请求
    10,             // 突发 10 个
)

// 在每个请求之前
limiter.Wait(ctx)
```

### 优势
- 防止 API 限流
- 平滑流量分布
- 可根据环境配置

### 成本优化

1. **缓存**：响应的 60 秒 TTL
2. **批量端点**：使用 `all_token_list` 而不是每链调用
3. **定期同步**：可配置间隔（默认 5 分钟）
4. **按需刷新**：仅在用户请求时

## 可扩展性考虑

### 当前限制
- 单服务器实例
- 内存中的速率限制
- 无分布式缓存

### 未来改进

1. **水平扩展**
   - 添加 Redis 用于分布式速率限制
   - 添加 Redis 用于缓存
   - 负载均衡器用于多实例

2. **数据库优化**
   - 读副本用于查询
   - 大表分区
   - 物化视图用于聚合

3. **异步处理**
   - 消息队列（RabbitMQ/Kafka）用于同步作业
   - 工作池用于并行处理
   - 作业状态跟踪

## 安全考虑

### 当前实现
- 所有端点的输入验证
- GORM 参数化查询（SQL 注入保护）
- 带可配置源的 CORS 中间件
- API 密钥存储在配置中（不在代码中）

### 生产要求
- [ ] 添加身份验证（JWT）
- [ ] 每用户速率限制
- [ ] API 调用的请求签名
- [ ] 审计日志
- [ ] HTTPS/TLS 强制
- [ ] 密钥管理（Vault）

## 测试策略

### 推荐测试

**单元测试**
```go
// 仓储测试
func TestWalletRepository_Create(t *testing.T) {
    // 测试数据库操作
}

// 带模拟的服务测试
func TestSyncService_SyncAddress(t *testing.T) {
    mockProvider := &MockProvider{}
    // 测试业务逻辑
}
```

**集成测试**
```go
// API 测试
func TestWalletAPI_CreateWallet(t *testing.T) {
    // 测试完整的 HTTP 流程
}
```

**前端测试**
```javascript
// 组件测试
describe('EVMAccounts', () => {
  it('正确显示钱包', () => {
    // 测试组件渲染
  })
})
```

## 部署

### 开发环境
```bash
# 后端
go run cmd/server/main.go

# 前端
cd frontend && npm run dev
```

### 生产环境
```bash
# 构建后端
make build

# 构建前端
cd frontend && npm run build

# 运行
./bin/rotki-demo
```

### Docker（推荐）
```bash
docker-compose up -d
```

## 监控

### 要跟踪的指标
- API 响应时间
- DeBank API 调用计数
- 同步成功/失败率
- 数据库查询性能
- 代币余额差异

### 日志
- 使用 Zap 的结构化日志
- 日志级别：DEBUG、INFO、WARN、ERROR
- 请求/响应日志
- 错误堆栈跟踪

## 未来增强

### 阶段 2
- [ ] 支持比特币/Solana 地址
- [ ] 历史余额图表
- [ ] 交易历史
- [ ] NFT 跟踪

### 阶段 3
- [ ] 多用户支持
- [ ] 身份验证/授权
- [ ] 自定义警报
- [ ] 导出功能

### 阶段 4
- [ ] 自托管区块链节点
- [ ] 自定义数据提供商
- [ ] 移除 DeBank 依赖
- [ ] 高级分析
