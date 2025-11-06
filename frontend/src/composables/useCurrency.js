import { ref, computed } from 'vue'

const currencies = ['USD', 'ETH', 'BTC']
const currencySymbols = {
  USD: '$',
  ETH: 'Ξ',
  BTC: '₿'
}

// Shared state
const selectedCurrency = ref('USD')
const exchangeRates = ref({
  USD: 1,
  ETH: 0,
  BTC: 0
})

export function useCurrency() {
  // Update exchange rates from token prices
  const updateExchangeRates = (addresses) => {
    let ethPrice = null
    let btcPrice = null

    addresses.forEach((addr) => {
      if (addr.tokens && addr.tokens.length > 0) {
        addr.tokens.forEach((token) => {
          // Look for ETH (native token symbol is usually ETH or WETH)
          if (
            (token.symbol === 'ETH' || token.symbol === 'WETH' || token.token_id === 'eth') &&
            token.price &&
            !ethPrice
          ) {
            ethPrice = token.price
          }
          // Look for BTC (wrapped BTC tokens)
          if (
            (token.symbol === 'WBTC' || token.symbol === 'BTC' || token.symbol.includes('BTC')) &&
            token.price &&
            !btcPrice
          ) {
            btcPrice = token.price
          }
        })
      }
    })

    // Update rates (1 ETH = X USD, so to convert USD to ETH: divide by ethPrice)
    if (ethPrice) {
      exchangeRates.value.ETH = ethPrice
      console.log('ETH price found:', ethPrice)
    } else {
      exchangeRates.value.ETH = 3500
      console.log('ETH price not found, using default:', 3500)
    }

    if (btcPrice) {
      exchangeRates.value.BTC = btcPrice
      console.log('BTC price found:', btcPrice)
    } else {
      exchangeRates.value.BTC = 95000
      console.log('BTC price not found, using default:', 95000)
    }

    console.log('Exchange rates updated:', exchangeRates.value)
  }

  // Format value based on selected currency
  const formatValue = (value) => {
    const rate = exchangeRates.value[selectedCurrency.value]
    if (rate === 0 || !rate) {
      return value.toFixed(2)
    }
    const convertedValue = value / rate
    return convertedValue.toFixed(selectedCurrency.value === 'USD' ? 2 : 6)
  }

  // Format token price based on selected currency
  const formatTokenPrice = (usdPrice) => {
    const rate = exchangeRates.value[selectedCurrency.value]
    if (rate === 0 || !rate) {
      return usdPrice.toFixed(2)
    }
    const convertedPrice = usdPrice / rate
    return convertedPrice.toFixed(selectedCurrency.value === 'USD' ? 2 : 6)
  }

  // Format token value based on selected currency
  const formatTokenValue = (usdValue) => {
    const rate = exchangeRates.value[selectedCurrency.value]
    if (rate === 0 || !rate) {
      return usdValue.toFixed(2)
    }
    const convertedValue = usdValue / rate
    return convertedValue.toFixed(selectedCurrency.value === 'USD' ? 2 : 6)
  }

  // Currency selection function
  const selectCurrency = (currency) => {
    selectedCurrency.value = currency
    localStorage.setItem('selectedCurrency', currency)
  }

  // Restore saved currency preference
  const restoreSavedCurrency = () => {
    const savedCurrency = localStorage.getItem('selectedCurrency')
    if (savedCurrency && currencies.includes(savedCurrency)) {
      selectedCurrency.value = savedCurrency
    }
  }

  return {
    currencies,
    currencySymbols,
    selectedCurrency,
    exchangeRates,
    updateExchangeRates,
    formatValue,
    formatTokenPrice,
    formatTokenValue,
    selectCurrency,
    restoreSavedCurrency
  }
}
