<template>
  <div class="currency-selector">
    <button class="currency-btn" @click="showMenu = !showMenu">
      <span class="currency-symbol">{{ currencySymbols[selectedCurrency] }}</span>
      <span class="currency-name">{{ selectedCurrency }}</span>
      <span class="dropdown-icon">▼</span>
    </button>
    <div v-if="showMenu" class="currency-menu">
      <div
        v-for="currency in currencies"
        :key="currency"
        class="currency-option"
        :class="{ active: selectedCurrency === currency }"
        @click="handleSelect(currency)"
      >
        <span class="currency-symbol-large">{{ currencySymbols[currency] }}</span>
        <div class="currency-info">
          <div class="currency-code">{{ currency }}</div>
          <div class="currency-description">Select as the main currency</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useCurrency } from '../composables/useCurrency'

const { currencies, currencySymbols, selectedCurrency, selectCurrency } = useCurrency()
const showMenu = ref(false)

const handleSelect = (currency) => {
  selectCurrency(currency)
  showMenu.value = false
}
</script>

<style scoped>
.currency-selector {
  position: relative;
}

.currency-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.currency-btn:hover {
  border-color: #4f46e5;
  background: #f9fafb;
}

.currency-symbol {
  font-size: 18px;
  font-weight: 600;
  color: #4f46e5;
}

.currency-name {
  font-weight: 500;
  color: #374151;
}

.dropdown-icon {
  font-size: 10px;
  color: #6b7280;
  transition: transform 0.2s;
}

.currency-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  min-width: 280px;
  z-index: 1000;
  overflow: hidden;
}

.currency-option {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  cursor: pointer;
  transition: background 0.2s;
  border-bottom: 1px solid #f3f4f6;
}

.currency-option:last-child {
  border-bottom: none;
}

.currency-option:hover {
  background: #f9fafb;
}

.currency-option.active {
  background: #eff6ff;
  border-left: 3px solid #4f46e5;
}

.currency-symbol-large {
  font-size: 28px;
  font-weight: 700;
  color: #4f46e5;
  width: 40px;
  text-align: center;
}

.currency-info {
  flex: 1;
}

.currency-code {
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 2px;
}

.currency-description {
  font-size: 12px;
  color: #6b7280;
}
</style>
