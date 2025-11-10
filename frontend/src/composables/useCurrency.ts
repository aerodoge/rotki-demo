import { ref, type Ref } from 'vue'
import type { Address } from '@/types'

type Currency = 'USD' | 'ETH' | 'BTC'

interface CurrencySymbols {
  [key: string]: string
}

interface ExchangeRates {
  [key: string]: number
}

const currencies: Currency[] = ['USD', 'ETH', 'BTC']
const currencySymbols: CurrencySymbols = {
  USD: '$',
  ETH: 'Ξ',
  BTC: '₿'
}

// 共享状态
const selectedCurrency: Ref<Currency> = ref('USD')
const exchangeRates: Ref<ExchangeRates> = ref({
  USD: 1,
  ETH: 0,
  BTC: 0
})

export function useCurrency() {
  // 从外部 API 获取 BTC 价格
  const fetchBTCPrice = async (): Promise<number | null> => {
    try {
      // 使用 CoinGecko 免费 API
      const response = await fetch('https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd')
      const data = await response.json()
      if (data.bitcoin && data.bitcoin.usd) {
        return data.bitcoin.usd
      }
    } catch (error) {
      console.error('Failed to fetch BTC price from CoinGecko:', error)
    }
    return null
  }

  // 从代币价格更新汇率
  const updateExchangeRates = async (addresses: Address[]): Promise<void> => {
    let ethPrice: number | null = null
    let btcPrice: number | null = null

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

    // 更新 ETH 汇率
    if (ethPrice) {
      exchangeRates.value.ETH = ethPrice
      console.log('ETH price found:', ethPrice)
    } else {
      exchangeRates.value.ETH = 3500
      console.log('ETH price not found, using default:', 3500)
    }

    // 更新 BTC 汇率
    if (btcPrice) {
      exchangeRates.value.BTC = btcPrice
      console.log('BTC price found from tokens:', btcPrice)
    } else {
      // 如果从代币中找不到 BTC 价格，从外部 API 获取
      console.log('BTC price not found in tokens, fetching from API...')
      const apiBtcPrice = await fetchBTCPrice()
      if (apiBtcPrice) {
        exchangeRates.value.BTC = apiBtcPrice
        console.log('BTC price fetched from API:', apiBtcPrice)
      } else {
        exchangeRates.value.BTC = 95000
        console.log('BTC price API failed, using default:', 95000)
      }
    }

    console.log('Exchange rates updated:', exchangeRates.value)
  }

  // 根据选定的货币格式化价值（带千位分隔符）
  const formatValue = (value: number): string => {
    const rate = exchangeRates.value[selectedCurrency.value]
    if (rate === 0 || !rate) {
      return value.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    }
    const convertedValue = value / rate
    const decimals = selectedCurrency.value === 'USD' ? 2 : 6
    return convertedValue.toLocaleString('en-US', { minimumFractionDigits: decimals, maximumFractionDigits: decimals })
  }

  // 根据选定的货币格式化代币价格（带千位分隔符）
  const formatTokenPrice = (usdPrice: number): string => {
    const rate = exchangeRates.value[selectedCurrency.value]
    if (rate === 0 || !rate) {
      return usdPrice.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    }
    const convertedPrice = usdPrice / rate
    const decimals = selectedCurrency.value === 'USD' ? 2 : 6
    return convertedPrice.toLocaleString('en-US', { minimumFractionDigits: decimals, maximumFractionDigits: decimals })
  }

  // 根据选定的货币格式化代币价值（带千位分隔符）
  const formatTokenValue = (usdValue: number): string => {
    const rate = exchangeRates.value[selectedCurrency.value]
    if (rate === 0 || !rate) {
      return usdValue.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    }
    const convertedValue = usdValue / rate
    const decimals = selectedCurrency.value === 'USD' ? 2 : 6
    return convertedValue.toLocaleString('en-US', { minimumFractionDigits: decimals, maximumFractionDigits: decimals })
  }

  // 货币选择函数
  const selectCurrency = (currency: Currency): void => {
    selectedCurrency.value = currency
    localStorage.setItem('selectedCurrency', currency)
  }

  // 恢复保存的货币偏好设置
  const restoreSavedCurrency = (): void => {
    const savedCurrency = localStorage.getItem('selectedCurrency')
    if (savedCurrency && currencies.includes(savedCurrency as Currency)) {
      selectedCurrency.value = savedCurrency as Currency
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
