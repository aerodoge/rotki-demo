# DeFi 协议支持功能实现

## 问题诊断

通过查看 DeBank 网站显示的地址 `0x23a5e45f9556dc7ffb507db8a3cfb2589bc8adad`，发现：

- **DeBank 显示**: 总资产 $22,284,543
  - Wallet (钱包): $0
  - Aave V3 (协议): $22,284,543
    - Supplied: 72,140.5338 ETH ($260,262,847)
    - Borrowed: 66,034.3333 stETH ($237,978,304)
    - Net Value: $22,284,543

- **原系统问题**: 只查询钱包代币资产，未查询 DeFi 协议资产

## 解决方案

### 1. 数据库模型扩展

**新增 Protocol 模型** (`internal/models/models.go`):
```go
type Protocol struct {
    ID            uint
    AddressID     uint
    ProtocolID    string  // aave-v3, compound, etc.
    Name          string
    SiteURL       string
    LogoURL       string
    ChainID       string
    NetUSDValue   float64  // 净值
    AssetUSDValue float64  // 资产值（供应的）
    DebtUSDValue  float64  // 债务值（借出的）
    PositionType  string   // lending, staking, liquidity
    RawData       JSONMap  // 详细数据
    LastUpdated   time.Time
}
```

**更新 Address 模型**: 添加 `Protocols []Protocol` 关联

### 2. 数据库迁移

创建 `migrations/002_add_protocols.sql`:
- 新增 `protocols` 表
- 建立与 `addresses` 表的外键关系
- 添加唯一索引 (address_id, protocol_id)

### 3. Repository 层

新增 `internal/repository/protocol_repository.go`:
- `GetByAddressID()` - 获取地址的所有协议
- `UpsertBatch()` - 批量更新或插入协议
- `DeleteByAddressID()` - 删除地址的协议
- `GetTotalValueByAddress()` - 获取协议总价值

### 4. Service 层更新

更新 `internal/service/sync_service.go`:
- 在 `SyncAddress()` 方法中添加协议同步逻辑
- 调用 `dataProvider.GetProtocolList()` 获取协议数据
- 将协议数据转换并保存到数据库
- 错误处理：协议获取失败不影响代币同步

### 5. API 层更新

更新 `internal/api/handler/address_handler.go`:
- 在所有返回地址数据的接口中添加协议信息
- `GetAddress()` - 单个地址查询包含协议
- `ListAddresses()` - 地址列表包含协议
- `RefreshAddress()` - 刷新后返回协议数据

### 6. 主程序更新

更新 `cmd/server/main.go`:
- 初始化 ProtocolRepository
- 将 ProtocolRepository 注入到 SyncService
- 将 ProtocolRepository 注入到 AddressHandler

## 功能特性

### 自动同步
- 定时同步服务会自动获取所有地址的协议持仓
- 新添加地址时后台自动同步协议数据

### 数据完整性
- 协议数据与代币数据分开存储
- 支持多种协议类型：Lending, Staking, Liquidity Pool 等
- 存储完整的协议详情（供应、借出、净值）

### API 返回格式

```json
{
  "id": 1,
  "address": "0x23a5e45f9556dc7ffb507db8a3cfb2589bc8adad",
  "tokens": [...],
  "protocols": [
    {
      "id": 1,
      "protocol_id": "aave3",
      "name": "Aave V3",
      "chain_id": "eth",
      "net_usd_value": 22284543.0,
      "asset_usd_value": 260262847.0,
      "debt_usd_value": 237978304.0,
      "position_type": "lending",
      "raw_data": {
        "portfolio_items": [...]
      }
    }
  ]
}
```

## DeBank API 使用

使用的 API 端点：
- `/v1/user/all_token_list` - 获取钱包代币
- `/v1/user/all_complex_protocol_list` - 获取协议持仓

## 前端集成 (待实现)

前端需要：
1. 更新 TypeScript 类型定义添加 `Protocol` 接口
2. 在地址详情页显示协议列表
3. 计算总资产时包含协议净值
4. 按协议类型分组显示（Lending, Staking 等）

## 测试

### 测试地址
使用 `0x23a5e45f9556dc7ffb507db8a3cfb2589bc8adad` 测试：
1. 添加地址到系统
2. 等待自动同步或手动刷新
3. 查询地址详情，应该看到：
   - Tokens: 钱包代币列表
   - Protocols: Aave V3 等协议持仓
   - 总净值应与 DeBank 显示一致

### 验证步骤
```bash
# 1. 启动后端
go run cmd/server/main.go

# 2. 添加地址（通过API或前端）
curl -X POST http://localhost:8080/api/v1/addresses \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_id": 1,
    "address": "0x23a5e45f9556dc7ffb507db8a3cfb2589bc8adad",
    "label": "Test Address",
    "chain_type": "EVM"
  }'

# 3. 等待自动同步（或手动刷新）
curl -X POST http://localhost:8080/api/v1/addresses/1/refresh

# 4. 查询地址数据
curl http://localhost:8080/api/v1/addresses/1
```

## 性能考虑

- **缓存策略**: DeBank API 响应缓存 60 秒
- **速率限制**: 使用 token bucket 算法限制请求频率
- **批量处理**: 同步服务支持批量并发处理
- **错误容错**: 协议获取失败不影响代币同步

## 未来优化

1. **协议详情展开**: 显示 portfolio items 详细信息
2. **历史追踪**: 记录协议持仓的历史快照
3. **收益计算**: 计算 Lending/Staking 收益率
4. **风险指标**: 显示健康因子等风险指标
5. **多协议聚合**: 跨协议的资产聚合视图

## 相关文件

### 后端
- `internal/models/models.go` - Protocol 模型
- `internal/repository/protocol_repository.go` - 数据库操作
- `internal/service/sync_service.go` - 同步逻辑
- `internal/api/handler/address_handler.go` - API 处理
- `migrations/002_add_protocols.sql` - 数据库迁移

### 前端 (需要更新)
- `frontend/src/types/index.ts` - TypeScript 类型
- `frontend/src/views/EVMAccounts.vue` - 展示界面
- `frontend/src/stores/wallet.ts` - 状态管理
