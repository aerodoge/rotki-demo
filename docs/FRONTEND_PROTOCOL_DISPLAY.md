# 前端协议展示建议

## 当前状态

✅ **后端已完成**: 协议数据已正确同步和返回
✅ **前端计算已修复**: 总资产已包含协议净值
⏳ **待实现**: 前端 UI 需要展示协议详情

## 推荐展示方式

### 1. 地址详情页 - 协议列表

参考 DeBank 的展示方式，在地址详情中添加协议部分：

```vue
<template>
  <div class="address-detail">
    <!-- 总资产 -->
    <div class="total-value">
      <h2>Total Value</h2>
      <p class="value">${{ formatNumber(getTotalValue(address)) }}</p>
    </div>

    <!-- 代币列表 -->
    <section class="tokens-section">
      <h3>Wallet Tokens</h3>
      <p class="section-value">${{ formatNumber(getTokenValue(address)) }}</p>
      <TokenList :tokens="address.tokens" />
    </section>

    <!-- 协议列表 -->
    <section v-if="address.protocols && address.protocols.length > 0" class="protocols-section">
      <h3>DeFi Protocols</h3>
      <p class="section-value">${{ formatNumber(getProtocolValue(address)) }}</p>

      <div v-for="protocol in address.protocols" :key="protocol.id" class="protocol-card">
        <div class="protocol-header">
          <img :src="protocol.logo_url" :alt="protocol.name" class="protocol-logo" />
          <div class="protocol-info">
            <h4>{{ protocol.name }}</h4>
            <span class="protocol-type">{{ formatProtocolType(protocol.position_type) }}</span>
          </div>
          <div class="protocol-value">
            <span class="net-value">${{ formatNumber(protocol.net_usd_value) }}</span>
          </div>
        </div>

        <!-- 可展开的详情 -->
        <div v-if="expandedProtocol === protocol.id" class="protocol-details">
          <div class="detail-row">
            <span class="label">Supplied / Assets:</span>
            <span class="value">${{ formatNumber(protocol.asset_usd_value) }}</span>
          </div>
          <div class="detail-row">
            <span class="label">Borrowed / Debt:</span>
            <span class="value negative">${{ formatNumber(protocol.debt_usd_value) }}</span>
          </div>
          <div class="detail-row total">
            <span class="label">Net Value:</span>
            <span class="value">${{ formatNumber(protocol.net_usd_value) }}</span>
          </div>

          <!-- 健康因子（如果适用） -->
          <div v-if="protocol.position_type === 'Lending'" class="health-factor">
            <span class="label">Health Factor:</span>
            <span class="value" :class="getHealthFactorClass(protocol)">
              {{ calculateHealthFactor(protocol) }}
            </span>
          </div>
        </div>

        <button @click="toggleProtocol(protocol.id)" class="expand-btn">
          {{ expandedProtocol === protocol.id ? 'Hide Details' : 'Show Details' }}
        </button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Address, Protocol } from '@/types'

const props = defineProps<{
  address: Address
}>()

const expandedProtocol = ref<number | null>(null)

const getTotalValue = (address: Address) => {
  const tokenValue = getTokenValue(address)
  const protocolValue = getProtocolValue(address)
  return tokenValue + protocolValue
}

const getTokenValue = (address: Address) => {
  return address.tokens?.reduce((sum, token) => sum + (token.usd_value || 0), 0) || 0
}

const getProtocolValue = (address: Address) => {
  return address.protocols?.reduce((sum, protocol) => sum + (protocol.net_usd_value || 0), 0) || 0
}

const formatProtocolType = (type: string) => {
  const typeMap: Record<string, string> = {
    'Lending': 'Lending',
    'Staked': 'Staking',
    'Liquidity Pool': 'LP',
    'Vesting': 'Vesting',
    'Farming': 'Farming'
  }
  return typeMap[type] || type
}

const formatNumber = (value: number) => {
  return value.toLocaleString('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

const calculateHealthFactor = (protocol: Protocol) => {
  if (protocol.debt_usd_value === 0) return '∞'
  const healthFactor = protocol.asset_usd_value / protocol.debt_usd_value
  return healthFactor.toFixed(2)
}

const getHealthFactorClass = (protocol: Protocol) => {
  const hf = parseFloat(calculateHealthFactor(protocol))
  if (hf > 2) return 'healthy'
  if (hf > 1.5) return 'warning'
  return 'danger'
}

const toggleProtocol = (protocolId: number) => {
  expandedProtocol.value = expandedProtocol.value === protocolId ? null : protocolId
}
</script>

<style scoped>
.protocol-card {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 12px;
  background: white;
}

.protocol-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.protocol-logo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}

.protocol-info {
  flex: 1;
}

.protocol-type {
  display: inline-block;
  padding: 2px 8px;
  background: #f3f4f6;
  border-radius: 4px;
  font-size: 12px;
  color: #6b7280;
}

.protocol-value .net-value {
  font-size: 18px;
  font-weight: 600;
  color: #059669;
}

.protocol-details {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
}

.detail-row.total {
  border-top: 1px solid #e5e7eb;
  margin-top: 8px;
  padding-top: 12px;
  font-weight: 600;
}

.value.negative {
  color: #dc2626;
}

.health-factor {
  margin-top: 12px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 6px;
  display: flex;
  justify-content: space-between;
}

.health-factor .value.healthy {
  color: #059669;
}

.health-factor .value.warning {
  color: #d97706;
}

.health-factor .value.danger {
  color: #dc2626;
}

.expand-btn {
  margin-top: 12px;
  width: 100%;
  padding: 8px;
  background: #f3f4f6;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  color: #4b5563;
  font-size: 14px;
}

.expand-btn:hover {
  background: #e5e7eb;
}
</style>
```

## 2. Dashboard 概览

在 Dashboard 上添加协议统计：

```vue
<template>
  <div class="dashboard">
    <!-- 总资产卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <h3>Total Assets</h3>
        <p class="value">${{ formatNumber(getTotalBalance()) }}</p>
      </div>

      <div class="stat-card">
        <h3>Wallet Tokens</h3>
        <p class="value">${{ formatNumber(getTotalTokenValue()) }}</p>
      </div>

      <div class="stat-card">
        <h3>DeFi Protocols</h3>
        <p class="value">${{ formatNumber(getTotalProtocolValue()) }}</p>
        <p class="count">{{ getTotalProtocolCount() }} protocols</p>
      </div>
    </div>

    <!-- 协议分布图表 -->
    <div class="protocol-distribution">
      <h3>Protocol Distribution</h3>
      <div v-for="item in getProtocolDistribution()" :key="item.protocol_id" class="protocol-item">
        <div class="protocol-name">{{ item.name }}</div>
        <div class="protocol-bar">
          <div class="bar-fill" :style="{ width: item.percentage + '%' }"></div>
        </div>
        <div class="protocol-amount">${{ formatNumber(item.value) }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useWalletStore } from '@/stores/wallet'

const walletStore = useWalletStore()

const getTotalTokenValue = () => {
  return walletStore.addresses.reduce((sum, addr) => {
    return sum + (addr.tokens?.reduce((s, t) => s + (t.usd_value || 0), 0) || 0)
  }, 0)
}

const getTotalProtocolValue = () => {
  return walletStore.addresses.reduce((sum, addr) => {
    return sum + (addr.protocols?.reduce((s, p) => s + (p.net_usd_value || 0), 0) || 0)
  }, 0)
}

const getTotalProtocolCount = () => {
  return walletStore.addresses.reduce((sum, addr) => {
    return sum + (addr.protocols?.length || 0)
  }, 0)
}

const getProtocolDistribution = () => {
  const protocolMap = new Map()
  const totalValue = getTotalProtocolValue()

  walletStore.addresses.forEach(addr => {
    addr.protocols?.forEach(protocol => {
      const existing = protocolMap.get(protocol.protocol_id)
      if (existing) {
        existing.value += protocol.net_usd_value
      } else {
        protocolMap.set(protocol.protocol_id, {
          protocol_id: protocol.protocol_id,
          name: protocol.name,
          value: protocol.net_usd_value
        })
      }
    })
  })

  return Array.from(protocolMap.values())
    .map(item => ({
      ...item,
      percentage: totalValue > 0 ? (item.value / totalValue) * 100 : 0
    }))
    .sort((a, b) => b.value - a.value)
}
</script>
```

## 3. EVMAccounts 页面增强

在地址列表中显示协议图标：

```vue
<template>
  <tr class="address-row">
    <td>{{ address.address }}</td>
    <td>
      <div class="protocols-preview">
        <img
          v-for="protocol in address.protocols?.slice(0, 3)"
          :key="protocol.id"
          :src="protocol.logo_url"
          :alt="protocol.name"
          :title="protocol.name"
          class="protocol-icon"
        />
        <span v-if="(address.protocols?.length || 0) > 3" class="more-count">
          +{{ address.protocols!.length - 3 }}
        </span>
      </div>
    </td>
    <td class="value-cell">
      <div class="value-breakdown">
        <span class="total">${{ formatNumber(getAddressValue(address)) }}</span>
        <span class="breakdown">
          Tokens: ${{ formatNumber(getTokenValue(address)) }}
          <span v-if="getProtocolValue(address) > 0">
            · Protocols: ${{ formatNumber(getProtocolValue(address)) }}
          </span>
        </span>
      </div>
    </td>
  </tr>
</template>

<style scoped>
.protocols-preview {
  display: flex;
  align-items: center;
  gap: 4px;
}

.protocol-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid white;
  margin-left: -8px;
}

.protocol-icon:first-child {
  margin-left: 0;
}

.more-count {
  margin-left: 4px;
  font-size: 12px;
  color: #6b7280;
}

.value-breakdown {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.value-breakdown .total {
  font-size: 16px;
  font-weight: 600;
}

.value-breakdown .breakdown {
  font-size: 12px;
  color: #6b7280;
}
</style>
```

## 4. 数据刷新指示器

添加协议数据同步状态：

```vue
<template>
  <div class="sync-status">
    <div v-if="isSyncing" class="syncing">
      <Loader class="spinner" />
      <span>Syncing protocols...</span>
    </div>
    <div v-else class="last-sync">
      <Check class="icon" />
      <span>Last synced: {{ formatLastSync(address.last_synced_at) }}</span>
    </div>
  </div>
</template>
```

## 5. 协议类型过滤

允许用户按协议类型筛选：

```vue
<template>
  <div class="protocol-filters">
    <button
      v-for="type in protocolTypes"
      :key="type"
      :class="{ active: selectedType === type }"
      @click="selectedType = type"
    >
      {{ type }}
    </button>
  </div>

  <div class="protocols-list">
    <ProtocolCard
      v-for="protocol in filteredProtocols"
      :key="protocol.id"
      :protocol="protocol"
    />
  </div>
</template>

<script setup lang="ts">
const protocolTypes = computed(() => {
  const types = new Set(['All'])
  walletStore.addresses.forEach(addr => {
    addr.protocols?.forEach(p => {
      if (p.position_type) types.add(p.position_type)
    })
  })
  return Array.from(types)
})

const selectedType = ref('All')

const filteredProtocols = computed(() => {
  const allProtocols = walletStore.addresses.flatMap(addr => addr.protocols || [])
  if (selectedType.value === 'All') return allProtocols
  return allProtocols.filter(p => p.position_type === selectedType.value)
})
</script>
```

## 实现优先级

### 第一阶段 (必须)
1. ✅ 修复总资产计算（已完成）
2. 🔄 在地址列表显示协议图标和总值
3. 🔄 地址详情显示协议列表

### 第二阶段 (推荐)
1. 协议详情展开/折叠
2. Dashboard 协议统计
3. 健康因子显示（Lending 类型）

### 第三阶段 (可选)
1. 协议类型过滤
2. 协议分布图表
3. 历史数据追踪

## 注意事项

1. **使用 net_usd_value**: 始终使用净值，不要用 asset_usd_value
2. **处理空数据**: 协议数据可能为空或同步中
3. **性能优化**: 大量协议时考虑虚拟滚动
4. **错误处理**: 显示同步失败状态
5. **实时更新**: 考虑 WebSocket 实时推送更新

## 参考资源

- DeBank UI: https://debank.com/profile/0x23a5e45f9556dc7ffb507db8a3cfb2589bc8adad
- Zerion UI: https://app.zerion.io/
- Zapper UI: https://zapper.xyz/

这些都是优秀的 DeFi 资产管理界面参考！
