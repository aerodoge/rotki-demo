import { ref, computed } from 'vue'

const currencies = ['USD', 'ETH', 'BTC']
const currencySymbols = {
  USD: '$',
  ETH: 'Ξ',
  BTC: '₿'
}

// 共享状态
const selectedCurrency = ref('USD')
const exchangeRates = ref({
  USD: 1,
  ETH: 0,
  BTC: 0
})

export function useCurrency() {
  // 从代币价格更新汇率
  const updateExchangeRates = (addresses) => {
    let ethPrice = null
    let btcPrice = null

    addresses.forEach((addr) => {
      if (addr.tokens && addr.tokens.length > 0) {
        addr.tokens.forEach((token) => {
          // 查找ETH（原生代币符号通常为ETH或WETH）
          if (
            (token.symbol === 'ETH' || token.symbol === 'WETH' || token.token_id === 'eth') &&
            token.price &&
            !ethPrice
          ) {
            ethPrice = token.price
          }
          // 查找BTC（包装BTC代币）
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

    // 更新汇率（1 ETH = X USD，所以将USD转换为ETH：除以ethPrice）
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

  // 根据选定的货币格式化价值
  const formatValue = (value) => {
    const rate = exchangeRates.value[selectedCurrency.value]
    if (rate === 0 || !rate) {
      return value.toFixed(2)
    }
    const convertedValue = value / rate
    return convertedValue.toFixed(selectedCurrency.value === 'USD' ? 2 : 6)
  }

  // 根据选定的货币格式化代币价格
  const formatTokenPrice = (usdPrice) => {
    const rate = exchangeRates.value[selectedCurrency.value]
    if (rate === 0 || !rate) {
      return usdPrice.toFixed(2)
    }
    const convertedPrice = usdPrice / rate
    return convertedPrice.toFixed(selectedCurrency.value === 'USD' ? 2 : 6)
  }

  // 根据选定的货币格式化代币价值
  const formatTokenValue = (usdValue) => {
    const rate = exchangeRates.value[selectedCurrency.value]
    if (rate === 0 || !rate) {
      return usdValue.toFixed(2)
    }
    const convertedValue = usdValue / rate
    return convertedValue.toFixed(selectedCurrency.value === 'USD' ? 2 : 6)
  }

  // 货币选择函数
  const selectCurrency = (currency) => {
    selectedCurrency.value = currency
    localStorage.setItem('selectedCurrency', currency)
  }

  // 恢复保存的货币偏好设置
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
