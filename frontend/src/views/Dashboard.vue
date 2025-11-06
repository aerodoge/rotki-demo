<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1>Dashboard</h1>
    </div>

    <!-- Total Balance Cards for all 3 currencies -->
    <div class="balance-cards">
      <div class="balance-card" v-for="currency in currencies" :key="currency">
        <div class="balance-label">Total Balance</div>
        <div class="balance-value">
          <span class="amount">{{ formatBalance(getTotalBalance(), currency) }}</span>
          <span class="currency-symbol">{{ currencySymbols[currency] }}</span>
        </div>
      </div>
    </div>

    <!-- Blockchain Balances -->
    <div class="section-card">
      <div class="section-header">
        <h2>Blockchain Balances</h2>
        <button class="refresh-btn" @click="refreshBalances">🔄</button>
      </div>
      <div class="blockchain-list">
        <div v-if="loading" class="loading">Loading...</div>
        <div v-else-if="chainBalances.length === 0" class="empty-state">
          No blockchain balances found
        </div>
        <div v-else class="chain-item" v-for="chain in chainBalances" :key="chain.chain_id">
          <div class="chain-info">
            <img v-if="chain.logo_url" :src="chain.logo_url" :alt="chain.name" class="chain-icon" />
            <span v-else class="chain-icon-text">{{ chain.name?.substring(0, 1) }}</span>
            <span class="chain-name">{{ chain.name || chain.chain_id }}</span>
          </div>
          <div class="chain-balance">
            {{ formatBalance(chain.balance, 'USD') }} $ /
            {{ formatBalance(chain.balance, 'ETH') }} Ξ /
            {{ formatBalance(chain.balance, 'BTC') }} ₿
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useWalletStore } from '../stores/wallet'
import { storeToRefs } from 'pinia'
import { useCurrency } from '../composables/useCurrency'

const walletStore = useWalletStore()
const { addresses, loading } = storeToRefs(walletStore)

const { currencies, currencySymbols, selectedCurrency, exchangeRates, updateExchangeRates } =
  useCurrency()

// Chain balances computed from addresses
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

// Total balance in USD
const getTotalBalance = () => {
  return addresses.value.reduce((total, address) => {
    if (address.tokens && address.tokens.length > 0) {
      return total + address.tokens.reduce((sum, token) => sum + (token.usd_value || 0), 0)
    }
    return total
  }, 0)
}

// Format balance in selected currency
const formatBalance = (usdValue, currency = null) => {
  const targetCurrency = currency || selectedCurrency.value
  const rate = exchangeRates.value[targetCurrency]

  if (!rate || rate === 0) {
    return usdValue.toFixed(2)
  }

  const convertedValue = usdValue / rate
  return convertedValue.toFixed(targetCurrency === 'USD' ? 2 : 6)
}

const refreshBalances = async () => {
  await walletStore.fetchAddresses()
  updateExchangeRates(addresses.value)
}

onMounted(async () => {
  await walletStore.fetchWallets()
  await walletStore.fetchAddresses()
  updateExchangeRates(addresses.value)
})
</script>

<style scoped>
.dashboard {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.dashboard-header h1 {
  font-size: 32px;
  font-weight: 600;
  color: #1f2937;
}

.balance-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.balance-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.balance-label {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 12px;
}

.balance-value {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.currency-symbol {
  font-size: 32px;
  font-weight: 700;
  color: #4f46e5;
}

.amount {
  font-size: 36px;
  font-weight: 700;
  color: #1f2937;
}

.section-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
}

.section-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.refresh-btn {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  padding: 4px;
  transition: transform 0.2s;
}

.refresh-btn:hover {
  transform: rotate(90deg);
}

.blockchain-list,
.asset-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chain-item,
.asset-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-radius: 8px;
  transition: background 0.2s;
}

.chain-item:hover,
.asset-item:hover {
  background: #f9fafb;
}

.chain-info,
.asset-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.chain-icon,
.asset-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
}

.chain-icon-text,
.asset-icon-text {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #e0e7ff 0%, #c7d2fe 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  color: #4f46e5;
  font-size: 14px;
  text-transform: uppercase;
}

.chain-name {
  font-size: 15px;
  font-weight: 500;
  color: #1f2937;
}

.chain-balance {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  font-family: monospace;
  white-space: nowrap;
}

.asset-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.asset-symbol {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.asset-name {
  font-size: 12px;
  color: #6b7280;
}

.asset-balance {
  text-align: right;
}

.balance-amount {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  font-family: monospace;
}

.balance-value {
  font-size: 12px;
  color: #6b7280;
  font-family: monospace;
}

.loading,
.empty-state {
  text-align: center;
  padding: 40px;
  color: #6b7280;
}
</style>
