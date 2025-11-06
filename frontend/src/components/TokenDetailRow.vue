<template>
  <tr class="token-detail-row">
    <td colspan="6">
      <div class="token-detail">
        <div class="token-icon-wrapper">
          <div class="token-main-icon">
            <img
              v-if="token.logo_url && token.logo_url.length > 0 && !failedImage"
              :src="token.logo_url"
              :alt="token.symbol"
              class="token-logo"
              @error="handleImageError"
            />
            <span v-else class="token-symbol-text">
              {{ token.symbol.substring(0, 2) }}
            </span>
          </div>
          <div class="token-chain-badge" :title="token.chain?.name || token.chain_id">
            <img
              v-if="token.chain?.logo_url"
              :src="token.chain.logo_url"
              :alt="token.chain_id"
              class="chain-logo-small"
            />
            <span v-else class="chain-text-small">{{ getChainIcon(token.chain_id) }}</span>
          </div>
        </div>
        <div class="token-info-detail">
          <div class="token-symbol-main">{{ token.symbol }}</div>
          <div class="token-name-sub">{{ token.name }}</div>
        </div>
        <div class="token-location">
          <span class="location-label">{{ token.chain?.name || token.chain_id }}</span>
        </div>
        <div class="token-price">
          <div class="price-amount">
            {{ currencySymbols[selectedCurrency] }}{{ formatTokenPrice(token.price) }}
          </div>
        </div>
        <div class="token-balance">
          <div class="balance-amount">{{ parseFloat(token.balance).toFixed(4) }}</div>
        </div>
        <div class="token-value">
          <div class="value-amount">
            {{ currencySymbols[selectedCurrency] }}{{ formatTokenValue(token.usd_value) }}
          </div>
        </div>
      </div>
    </td>
  </tr>
</template>

<script setup>
import { ref } from 'vue'
import { useCurrency } from '../composables/useCurrency'

const props = defineProps({
  token: {
    type: Object,
    required: true
  }
})

const { currencySymbols, selectedCurrency, formatTokenPrice, formatTokenValue } = useCurrency()
const failedImage = ref(false)

const handleImageError = () => {
  failedImage.value = true
}

const getChainIcon = (chainId) => {
  const icons = {
    eth: '⟠',
    bsc: 'B',
    polygon: 'P',
    arbitrum: 'A',
    optimism: 'O'
  }
  return icons[chainId] || '⛓️'
}
</script>

<style scoped>
.token-detail-row {
  background: #fafafa;
}

.token-detail {
  display: grid;
  grid-template-columns: 50px 1fr 150px 150px 150px 150px;
  gap: 16px;
  padding: 12px 16px;
  align-items: center;
}

.token-icon-wrapper {
  position: relative;
  width: 40px;
  height: 40px;
}

.token-main-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  background: white;
  border: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: center;
}

.token-logo {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.token-symbol-text {
  font-size: 14px;
  font-weight: 700;
  color: #4f46e5;
  text-transform: uppercase;
}

.token-chain-badge {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: white;
  border: 2px solid white;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.chain-logo-small {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.chain-text-small {
  font-size: 8px;
}

.token-info-detail {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.token-symbol-main {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.token-name-sub {
  font-size: 12px;
  color: #6b7280;
}

.token-location {
  display: flex;
  align-items: center;
}

.location-label {
  padding: 4px 12px;
  background: #f3f4f6;
  border-radius: 12px;
  font-size: 12px;
  color: #6b7280;
}

.token-price,
.token-balance,
.token-value {
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.price-amount,
.balance-amount,
.value-amount {
  font-size: 13px;
  color: #1f2937;
  font-weight: 500;
}
</style>
