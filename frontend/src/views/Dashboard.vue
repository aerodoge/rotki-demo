<template>
  <div class="container mx-auto p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <h1 class="text-3xl font-bold tracking-tight">Dashboard</h1>
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input
            type="checkbox"
            v-model="hideSmallBalances"
            class="w-4 h-4 rounded border-gray-300 cursor-pointer"
          />
          <span class="text-muted-foreground">Hide (&lt;$10)</span>
        </label>
      </div>
      <Button @click="refreshBalances" variant="outline" size="sm">
        <span class="mr-2">🔄</span>
        Refresh
      </Button>
    </div>

    <!-- Balance Cards -->
    <div class="grid gap-4 md:grid-cols-3">
      <Card v-for="currency in currencies" :key="currency">
        <CardHeader class="pb-2">
          <CardDescription>Total Balance</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ formatBalance(getTotalBalance(), currency) }}
            <span class="text-muted-foreground ml-1">{{ currencySymbols[currency] }}</span>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Blockchain Balances -->
    <Card>
      <CardHeader>
        <CardTitle>Blockchain Balances</CardTitle>
        <CardDescription>View your assets across different blockchains</CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="flex items-center justify-center py-8">
          <div class="text-muted-foreground">Loading...</div>
        </div>
        <div v-else-if="chainBalances.length === 0" class="flex items-center justify-center py-8">
          <div class="text-muted-foreground">No blockchain balances found</div>
        </div>
        <div v-else class="space-y-3">
          <div
            v-for="chain in filteredChainBalances"
            :key="chain.chain_id"
            class="flex items-center justify-between p-4 rounded-lg border bg-card hover:bg-accent/50 transition-colors"
          >
            <div class="flex items-center gap-3">
              <div class="relative w-8 h-8">
                <img
                  :src="getChainLogo(chain.chain_id)"
                  :alt="chain.name"
                  class="w-8 h-8 rounded-full"
                  @error="(e) => handleImageError(e, chain)"
                />
                <div
                  v-if="imageErrors.has(chain.chain_id)"
                  class="absolute inset-0 w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center"
                >
                  <span class="text-sm font-semibold">{{ chain.name?.substring(0, 1) }}</span>
                </div>
              </div>
              <div>
                <div class="font-medium">{{ chain.name || chain.chain_id }}</div>
                <div class="text-sm text-muted-foreground">{{ chain.chain_id }}</div>
              </div>
            </div>
            <div class="text-right">
              <div class="font-semibold">{{ formatBalance(chain.balance, 'USD') }} $</div>
              <div class="text-sm text-muted-foreground">
                {{ formatBalance(chain.balance, 'ETH') }} Ξ /
                {{ formatBalance(chain.balance, 'BTC') }} ₿
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useWalletStore } from '@/stores/wallet'
import { storeToRefs } from 'pinia'
import { useCurrency } from '@/composables/useCurrency'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

const walletStore = useWalletStore()
const { addresses, loading } = storeToRefs(walletStore)

// 跟踪图片加载失败的链
const imageErrors = ref(new Set<string>())

// 隐藏小额资产
const hideSmallBalances = ref(false)

const { currencies, currencySymbols, selectedCurrency, exchangeRates, updateExchangeRates } =
  useCurrency()

// 从地址计算链余额
const chainBalances = computed(() => {
  const balanceMap = new Map()

  addresses.value.forEach((address) => {
    // 1. 累加钱包代币（不属于协议的代币）
    if (address.tokens && address.tokens.length > 0) {
      address.tokens.forEach((token) => {
        // 只统计钱包代币，协议代币的价值已经包含在协议净值中
        if (!token.protocol_id) {
          const chainId = token.chain_id
          const chain = token.chain

          if (!balanceMap.has(chainId)) {
            balanceMap.set(chainId, {
              chain_id: chainId,
              name: chain?.name || chainId,
              logo_url: chain?.logo_url,
              balance: 0
            })
          }

          const chainBalance = balanceMap.get(chainId)
          chainBalance.balance += token.usd_value || 0
        }
      })
    }

    // 2. 累加协议净值
    if (address.protocols && address.protocols.length > 0) {
      address.protocols.forEach((protocol) => {
        const chainId = protocol.chain_id

        if (!balanceMap.has(chainId)) {
          // 如果该链还没有记录，尝试从协议的链信息获取名称和logo
          balanceMap.set(chainId, {
            chain_id: chainId,
            name: chainId, // 可以从 chain_info 表获取，暂时用 chain_id
            logo_url: undefined,
            balance: 0
          })
        }

        const chainBalance = balanceMap.get(chainId)
        chainBalance.balance += protocol.net_usd_value || 0
      })
    }
  })

  return Array.from(balanceMap.values()).sort((a, b) => b.balance - a.balance)
})

// 过滤小额资产的链余额
const filteredChainBalances = computed(() => {
  if (!hideSmallBalances.value) {
    return chainBalances.value
  }
  return chainBalances.value.filter(chain => chain.balance >= 10)
})

// 以USD计算的总余额
const getTotalBalance = () => {
  return addresses.value.reduce((total, address) => {
    // 只计算钱包代币（不属于任何协议的代币）
    const walletTokenValue =
      address.tokens?.reduce((sum, token) => {
        if (!token.protocol_id) {
          return sum + (token.usd_value || 0)
        }
        return sum
      }, 0) || 0
    // 协议净值已经包含了协议代币的价值
    const protocolValue =
      address.protocols?.reduce((sum, protocol) => sum + (protocol.net_usd_value || 0), 0) || 0
    return total + walletTokenValue + protocolValue
  }, 0)
}

// 将余额格式化为选定的货币（带千位分隔符）
const formatBalance = (usdValue: number, currency: string | null = null): string => {
  const targetCurrency = currency || selectedCurrency.value
  const rate = exchangeRates.value[targetCurrency]

  if (!rate || rate === 0) {
    return '0.00'
  }

  const value = usdValue / rate

  if (targetCurrency === 'BTC') {
    return value.toLocaleString('en-US', { minimumFractionDigits: 4, maximumFractionDigits: 4 })
  } else if (targetCurrency === 'ETH') {
    return value.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
  } else {
    return value.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
  }
}

// 获取链的 logo，优先使用本地图片
const getChainLogo = (chainId: string): string => {
  // 优先使用本地图片
  return `/images/chains/${chainId}.png`
}

// 图片加载失败时的处理
const handleImageError = (event: Event, chain: any) => {
  imageErrors.value.add(chain.chain_id)
  const img = event.target as HTMLImageElement
  img.style.display = 'none'
}

const refreshBalances = async () => {
  await walletStore.fetchAddresses()
  await updateExchangeRates(addresses.value)
}

onMounted(async () => {
  await walletStore.fetchAddresses()
  await updateExchangeRates(addresses.value)
})
</script>
