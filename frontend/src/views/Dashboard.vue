<template>
  <div class="container mx-auto p-6 space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold tracking-tight">Dashboard</h1>
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
            v-for="chain in chainBalances"
            :key="chain.chain_id"
            class="flex items-center justify-between p-4 rounded-lg border bg-card hover:bg-accent/50 transition-colors"
          >
            <div class="flex items-center gap-3">
              <img
                v-if="chain.logo_url"
                :src="chain.logo_url"
                :alt="chain.name"
                class="w-8 h-8 rounded-full"
              />
              <div
                v-else
                class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center"
              >
                <span class="text-sm font-semibold">{{ chain.name?.substring(0, 1) }}</span>
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
import { computed, onMounted } from 'vue'
import { useWalletStore } from '@/stores/wallet'
import { storeToRefs } from 'pinia'
import { useCurrency } from '@/composables/useCurrency'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

const walletStore = useWalletStore()
const { addresses, loading } = storeToRefs(walletStore)

const { currencies, currencySymbols, selectedCurrency, exchangeRates, updateExchangeRates } =
  useCurrency()

// 从地址计算链余额
const chainBalances = computed(() => {
  const balanceMap = new Map()

  addresses.value.forEach((address) => {
    if (address.tokens && address.tokens.length > 0) {
      address.tokens.forEach((token) => {
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
      })
    }
  })

  return Array.from(balanceMap.values()).sort((a, b) => b.balance - a.balance)
})

// 以USD计算的总余额
const getTotalBalance = () => {
  return addresses.value.reduce((total, address) => {
    if (address.tokens && address.tokens.length > 0) {
      return total + address.tokens.reduce((sum, token) => sum + (token.usd_value || 0), 0)
    }
    return total
  }, 0)
}

// 将余额格式化为选定的货币
const formatBalance = (usdValue: number, currency: string | null = null): string => {
  const targetCurrency = currency || selectedCurrency.value
  const rate = exchangeRates.value[targetCurrency]

  if (!rate || rate === 0) {
    return '0.00'
  }

  const value = usdValue / rate

  if (targetCurrency === 'BTC') {
    return value.toFixed(8)
  } else if (targetCurrency === 'ETH') {
    return value.toFixed(6)
  } else {
    return value.toFixed(2)
  }
}

const refreshBalances = async () => {
  await walletStore.fetchAddresses()
  updateExchangeRates(addresses.value)
}

onMounted(async () => {
  await walletStore.fetchAddresses()
  updateExchangeRates(addresses.value)
})
</script>
