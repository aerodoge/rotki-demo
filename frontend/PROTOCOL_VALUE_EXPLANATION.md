# 协议资产净值说明

## DeBank 数据示例

以地址 `0x23a5e45f9556dc7ffb507db8a3cfb2589bc8adad` 为例：

### Aave V3 协议
- **Supplied (提供的资产)**: 72,140.5338 ETH = $260,262,847
- **Borrowed (借出的资产)**: 66,034.3333 stETH = $237,978,304
- **Net Value (净值)**: $260,262,847 - $237,978,304 = **$22,284,543**

## 系统实现

### 数据流

1. **DeBank API 返回**:
```json
{
  "id": "aave3",
  "name": "Aave V3",
  "chain": "eth",
  "net_usd_value": 22284543.0,      // ← 这是净值！
  "asset_usd_value": 260262847.0,   // Supplied
  "debt_usd_value": 237978304.0     // Borrowed
}
```

2. **Provider 层** (`internal/provider/debank/debank.go`):
   - 直接映射 DeBank 返回的 `net_usd_value`
   - ✅ **不需要计算**，DeBank 已经计算好了

3. **数据库存储** (`protocols` 表):
```sql
CREATE TABLE protocols (
  id BIGINT PRIMARY KEY,
  net_usd_value DECIMAL(30,6),    -- 净值 (Asset - Debt)
  asset_usd_value DECIMAL(30,6),  -- 资产值
  debt_usd_value DECIMAL(30,6),   -- 债务值
  ...
);
```

4. **API 返回给前端**:
```json
{
  "protocols": [
    {
      "protocol_id": "aave3",
      "name": "Aave V3",
      "net_usd_value": 22284543.0,     // ← 前端应该使用这个！
      "asset_usd_value": 260262847.0,
      "debt_usd_value": 237978304.0
    }
  ]
}
```

## 前端展示建议

### 总资产计算
```typescript
// 计算地址总资产
const getTotalValue = (address: Address) => {
  // 钱包代币价值
  const tokenValue = address.tokens?.reduce((sum, token) =>
    sum + (token.usd_value || 0), 0
  ) || 0

  // 协议净值
  const protocolValue = address.protocols?.reduce((sum, protocol) =>
    sum + protocol.net_usd_value, 0  // ← 使用 net_usd_value
  ) || 0

  return tokenValue + protocolValue
}
```

### 协议详情展示
```vue
<div v-for="protocol in protocols" :key="protocol.id">
  <div class="protocol-header">
    <h3>{{ protocol.name }}</h3>
    <span class="net-value">
      ${{ protocol.net_usd_value.toLocaleString() }}
    </span>
  </div>

  <div class="protocol-details">
    <div class="supplied">
      Supplied: ${{ protocol.asset_usd_value.toLocaleString() }}
    </div>
    <div class="borrowed">
      Borrowed: ${{ protocol.debt_usd_value.toLocaleString() }}
    </div>
    <div class="net-value-detail">
      Net: ${{ protocol.net_usd_value.toLocaleString() }}
    </div>
  </div>
</div>
```

## 验证正确性

### 测试地址
`0x23a5e45f9556dc7ffb507db8a3cfb2589bc8adad`

### 期望结果
- 钱包代币: $0
- Aave V3 净值: $22,284,543
- **总资产**: $22,284,543

### API 测试
```bash
# 1. 添加地址
curl -X POST http://localhost:8080/api/v1/addresses \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_id": 1,
    "address": "0x23a5e45f9556dc7ffb507db8a3cfb2589bc8adad",
    "label": "Test Aave Position",
    "chain_type": "EVM"
  }'

# 2. 等待同步（约30秒）或手动刷新
curl -X POST http://localhost:8080/api/v1/addresses/1/refresh

# 3. 查询结果
curl http://localhost:8080/api/v1/addresses/1 | jq '.protocols[0]'

# 应该看到：
{
  "protocol_id": "aave3",
  "name": "Aave V3",
  "net_usd_value": 22284543.0,
  "asset_usd_value": 260262847.0,
  "debt_usd_value": 237978304.0
}
```

## 常见问题

### Q: 为什么不直接显示 asset_usd_value？
A: 因为那不是真实资产！如果我供应了 $260M 但借出了 $238M，我的净资产只有 $22M。

### Q: DeBank 如何计算 net_usd_value？
A: `net_usd_value = asset_usd_value - debt_usd_value`

### Q: 系统是否需要手动计算净值？
A: **不需要**！DeBank API 已经返回了计算好的 `net_usd_value`，我们直接使用即可。

### Q: 如果协议只有 Staking 没有 Borrow 呢？
A:
- `asset_usd_value`: Staking 的代币价值
- `debt_usd_value`: 0
- `net_usd_value`: 等于 asset_usd_value

### Q: 如何在前端显示详细信息？
A: 可以展开显示：
```
Aave V3: $22,284,543 ▼
  ├─ Supplied: $260,262,847 (72,140 ETH)
  ├─ Borrowed: $237,978,304 (66,034 stETH)
  └─ Net Value: $22,284,543
```

## 总结

✅ 系统**已经正确实现**了净值计算
✅ 数据库存储的是 `net_usd_value`（净值）
✅ API 返回的是 `net_usd_value`（净值）
✅ 前端只需要使用 `protocol.net_usd_value` 即可

**重要**: 计算总资产时，使用 `net_usd_value` 而不是 `asset_usd_value`！
