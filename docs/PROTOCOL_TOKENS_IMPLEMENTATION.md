# 协议代币展平实现方案

## 背景

根据 Rotki 的展示方式（参考 img.png），DeFi 协议中的代币应该**直接展示在代币列表中**，而不是只显示协议的净值。

### Rotki 的展示方式

对于地址 `0x23a5e45f9556dc7ffb507db8a3cfb2589bc8adad`（Aave V3 持仓）：

```
Address: 0x23a5...aDAD
Total: $260,351,579.46

代币列表：
├─ aEthrsETH (Aave Ethereum rsETH)
│  Chain: Ethereum
│  Amount: 68,181.14
│  Price: $3,818.53
│  USD Value: $260,351,579.46  ← 正值（供应）
│
├─ stETH (Liquid staked Ether 2.0)
│  Chain: Ethereum
│  Amount: 0.000000
│  Price: $3,605.37
│  USD Value: $0.00
│
└─ variableDebtEthwstETH (Aave Ethereum Variable Debt wstETH)
   Chain: Ethereum
   Amount: 54,182.39
   Price: $0.00  ← Debt 代币价格需要特殊处理
   USD Value: $0.00  ← 应该显示负值！
```

### 问题

1. **只显示净值不符合 Rotki 风格**
2. **需要展开显示所有协议代币**（包括供应的和借出的）
3. **Debt 代币应该显示为负值**

## 解决方案

### 核心思路

**把协议中的代币"展平"到普通代币列表中**，使用 `protocol_id` 和 `is_debt` 字段标记来源。

### 数据流

```
DeBank API
  ↓
portfolio_item_list[].asset_token_list  ← 包含正负值的完整代币列表
  ├─ ETH: amount = 72140.53 (正值 = 供应)
  └─ stETH: amount = -66034.33 (负值 = 借出)
  ↓
Provider 层解析
  ↓
转换为 TokenDetail
  ├─ Amount: 保留符号（正/负）
  ├─ USDValue: amount * price
  └─ IsDebt: amount < 0
  ↓
存入 tokens 表
  ├─ balance: "72140.53" 或 "-66034.33"
  ├─ usd_value: 正数或负数
  ├─ protocol_id: "aave3"
  └─ is_debt: true/false
  ↓
前端展示
  └─ 所有代币统一展示在列表中
```

## 实现细节

### 1. 数据模型更新

#### Token 表结构
```sql
CREATE TABLE tokens (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    address_id BIGINT NOT NULL,
    chain_id VARCHAR(50) NOT NULL,
    token_id VARCHAR(255) NOT NULL,
    symbol VARCHAR(50),
    name VARCHAR(255),
    decimals INT,
    logo_url VARCHAR(512),
    balance VARCHAR(255),           -- 可以是负数字符串
    price DECIMAL(30, 6),
    usd_value DECIMAL(30, 6),       -- 可以是负数
    protocol_id VARCHAR(100),        -- NEW: 协议ID
    is_debt TINYINT(1) DEFAULT 0,   -- NEW: 债务标记
    ...
);
```

### 2. Provider 层增强

#### TokenDetail 结构
```go
type TokenDetail struct {
    TokenID   string
    ChainID   string
    Symbol    string
    Name      string
    Decimals  int
    LogoURL   string
    Amount    float64  // 可以是负数
    Price     float64
    USDValue  float64  // amount * price（保留符号）
    IsDebt    bool     // amount < 0
}
```

#### DeBank API 解析
```go
// asset_token_list 包含所有代币（正负值）
for _, token := range item.AssetTokenList {
    amount := token.Amount
    isDebt := amount < 0

    // 确保 debt 是负数
    if isDebt && amount > 0 {
        amount = -amount
    }

    tokens = append(tokens, TokenDetail{
        Amount: amount,
        USDValue: amount * token.Price,
        IsDebt: isDebt,
    })
}
```

### 3. Sync Service 展平逻辑

```go
// 收集所有协议代币
protocolTokens := make([]models.Token, 0)

for _, proto := range protocols {
    for _, item := range proto.PortfolioItems {
        for _, tokenDetail := range item.AssetTokenList {
            // 构造代币名称
            tokenName := tokenDetail.Name
            if tokenDetail.IsDebt {
                if !strings.Contains(tokenName, "debt") {
                    tokenName = "Debt " + tokenName
                }
            }

            // 添加到代币列表
            protocolTokens = append(protocolTokens, models.Token{
                AddressID:  addressID,
                ChainID:    tokenDetail.ChainID,
                TokenID:    tokenDetail.TokenID,
                Symbol:     tokenDetail.Symbol,
                Name:       tokenName,
                Balance:    fmt.Sprintf("%.18f", tokenDetail.Amount), // 保留符号
                Price:      tokenDetail.Price,
                USDValue:   tokenDetail.USDValue,  // 可以是负数
                ProtocolID: proto.ProtocolID,      // 标记来源
                IsDebt:     tokenDetail.IsDebt,
            })
        }
    }
}

// 插入到 tokens 表
s.tokenRepo.UpsertBatch(protocolTokens)
```

### 4. 前端展示

#### TypeScript 类型
```typescript
export interface Token {
  balance: string      // 可以是负数字符串
  usd_value: number    // 可以是负数
  protocol_id?: string // 标记来源协议
  is_debt?: boolean    // 债务标记
}
```

#### 展示逻辑
```vue
<template>
  <div v-for="token in address.tokens" :key="token.id" class="token-row">
    <img :src="token.logo_url" :alt="token.symbol" />

    <div class="token-info">
      <span class="symbol">{{ token.symbol }}</span>
      <span class="name">{{ token.name }}</span>

      <!-- 显示协议标签 -->
      <span v-if="token.protocol_id" class="protocol-tag">
        {{ token.protocol_id }}
      </span>
    </div>

    <div class="token-balance">
      {{ parseFloat(token.balance).toFixed(6) }}
    </div>

    <div class="token-value" :class="{ negative: token.is_debt }">
      ${{ token.usd_value.toLocaleString() }}
    </div>
  </div>
</template>

<style>
.token-value.negative {
  color: #dc2626; /* 红色显示负值 */
}
</style>
```

## Debt 代币价格处理

### 问题
img.png 显示 `variableDebtEthwstETH` 的价格为 $0.00，但实际上应该用基础资产的价格。

### 解决方案

你说："价格就取 eth 的价格"

DeBank API 已经返回了正确的价格：
```json
{
  "id": "0xae7ab96520de3a18e5e111b5eaab095312d7fe84",
  "symbol": "stETH",
  "price": 3605.37,  ← 已经有价格！
  "amount": -66034.33
}
```

所以我们的实现已经正确：
```go
USDValue: amount * token.Price  // -66034.33 * 3605.37 = -$237,978,304
```

## 货币本位切换

你提到："但切换本位的时候，还是要变哦"

### 当前状态
前端已经有 `selectedCurrency` 状态（USD, EUR, BTC 等）

### 需要做的
1. **获取汇率数据**（如 USD→EUR, USD→BTC）
2. **实时转换显示**

```typescript
const convertValue = (usdValue: number, targetCurrency: string) => {
  const rates = {
    'USD': 1,
    'EUR': 0.92,    // 1 USD = 0.92 EUR
    'BTC': 0.000026, // 1 USD = 0.000026 BTC
    // ...
  }
  return usdValue * (rates[targetCurrency] || 1)
}

// 使用
const displayValue = convertValue(token.usd_value, selectedCurrency)
```

## 测试验证

### 测试地址
`0x23a5e45f9556dc7ffb507db8a3cfb2589bc8adad`

### 期望结果

**钱包代币**: 0 个

**协议代币（Aave V3）**:
1. ETH (supplied)
   - Amount: 72,140.53
   - Price: $3,609.87
   - Value: $260,417,948.75 (正值)

2. stETH (borrowed)
   - Amount: -66,034.33
   - Price: $3,605.37
   - Value: -$238,077,937.68 (负值)

**净资产**: $260,417,948.75 - $238,077,937.68 = **$22,340,011.07**

### API 测试
```bash
# 添加地址
curl -X POST http://localhost:8080/api/v1/addresses \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_id": 1,
    "address": "0x23a5e45f9556dc7ffb507db8a3cfb2589bc8adad",
    "label": "Aave Position",
    "chain_type": "EVM"
  }'

# 等待同步或手动刷新
curl -X POST http://localhost:8080/api/v1/addresses/1/refresh

# 查询代币列表
curl http://localhost:8080/api/v1/addresses/1 | jq '.tokens'

# 应该看到两个代币（一个正值，一个负值）
```

## 优势

1. ✅ **统一展示**: 钱包代币和协议代币都在一个列表中
2. ✅ **完整信息**: 保留了供应和借出的详细信息
3. ✅ **正确计算**: 自动计算净值（正值 + 负值）
4. ✅ **可扩展**: 支持任何类型的协议（Lending, Staking, LP 等）
5. ✅ **可追溯**: 通过 `protocol_id` 可以追溯到具体协议

## 与之前方案的对比

### 之前的方案
```
展示：净值 = $22,284,543
问题：看不到具体持仓详情
```

### 新方案
```
展示：
  ETH (Aave): +$260,417,948
  stETH (Debt): -$238,077,937
  ──────────────────────────
  净值: $22,340,011
优势：完整透明的持仓信息
```

## 注意事项

1. **唯一键冲突**:
   - 协议代币的 `token_id` 可能与钱包代币冲突
   - 使用 `protocol_id` 区分

2. **负值处理**:
   - Balance 存储为字符串，支持负号
   - USDValue 使用 DECIMAL，支持负数

3. **过滤垃圾代币**:
   - Spam 过滤仍然有效
   - 不过滤协议代币

4. **性能**:
   - 协议代币数量通常不多（<10个/协议）
   - 对查询性能影响小

## 下一步

1. **重新初始化数据库**（删除现有数据，应用新表结构）
2. **启动后端测试**
3. **添加测试地址**
4. **验证代币列表**是否包含协议代币
5. **检查负值显示**是否正确
6. **实现货币本位切换**（如需要）

---

**核心要点**: 把协议代币当作普通代币处理，用 `protocol_id` 和 `is_debt` 标记来源，前端统一展示！
