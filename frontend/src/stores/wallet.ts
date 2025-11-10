import { defineStore } from 'pinia'
import { walletsAPI, addressesAPI } from '@/api/client'
import type { Wallet, Address, Token, CreateWalletRequest, UpdateWalletRequest, CreateAddressRequest, UpdateAddressRequest } from '@/types'

// 过滤垃圾代币的辅助函数
const filterSpamTokens = (tokens: Token[] | undefined): Token[] => {
  if (!tokens || !Array.isArray(tokens)) return []

  const spamKeywords = [
    't.me/',
    't.ly/',
    'fli.so/',
    'wr.do/',
    'www.',
    'claim',
    'swap',
    'redeem',
    'visit',
    'airdrop',
    'reward',
    'voucher',
    'distribution'
  ]

  return tokens.filter((token) => {
    const symbolLower = (token.symbol || '').toLowerCase()
    const nameLower = (token.name || '').toLowerCase()

    // 检查垃圾关键词
    for (const keyword of spamKeywords) {
      if (symbolLower.includes(keyword) || nameLower.includes(keyword)) {
        return false
      }
    }

    // 检查价值为$0的代币中的可疑模式
    if (token.price === 0 && token.usd_value === 0) {
      if (symbolLower.includes('✅') || nameLower.includes('✅')) {
        return false
      }
      if (symbolLower.includes('$') && !symbolLower.startsWith('$')) {
        return false
      }
    }

    return true
  })
}

interface WalletState {
  wallets: Wallet[]
  addresses: Address[]
  selectedWallet: Wallet | null
  loading: boolean
  error: string | null
}

export const useWalletStore = defineStore('wallet', {
  state: (): WalletState => ({
    wallets: [],
    addresses: [],
    selectedWallet: null,
    loading: false,
    error: null
  }),

  getters: {
    getAddressesByWallet: (state) => (walletId: number): Address[] => {
      return state.addresses.filter((addr) => addr.wallet_id === walletId)
    },

    getTotalValueByWallet: (state) => (walletId: number): number => {
      const addresses = state.addresses.filter((addr) => addr.wallet_id === walletId)
      return addresses.reduce((sum, addr) => {
        // 只计算钱包代币（不属于任何协议的代币）
        const walletTokenValue = addr.tokens?.reduce((s, t) => {
          if (!t.protocol_id) {
            return s + (t.usd_value || 0)
          }
          return s
        }, 0) || 0
        // 协议净值已经包含了协议代币的价值
        const protocolValue = addr.protocols?.reduce((s, p) => s + (p.net_usd_value || 0), 0) || 0
        return sum + walletTokenValue + protocolValue
      }, 0)
    },

    getTotalValue: (state): number => {
      return state.addresses.reduce((sum, addr) => {
        // 只计算钱包代币（不属于任何协议的代币）
        const walletTokenValue = addr.tokens?.reduce((s, t) => {
          if (!t.protocol_id) {
            return s + (t.usd_value || 0)
          }
          return s
        }, 0) || 0
        // 协议净值已经包含了协议代币的价值
        const protocolValue = addr.protocols?.reduce((s, p) => s + (p.net_usd_value || 0), 0) || 0
        return sum + walletTokenValue + protocolValue
      }, 0)
    }
  },

  actions: {
    async fetchWallets(): Promise<void> {
      this.loading = true
      this.error = null
      try {
        const response = await walletsAPI.list()
        this.wallets = response.data
      } catch (error: any) {
        this.error = error.message
        console.error('Failed to fetch wallets:', error)
      } finally {
        this.loading = false
      }
    },

    async fetchAddresses(walletId: number | null = null): Promise<void> {
      this.loading = true
      this.error = null
      try {
        const response = await addressesAPI.list(walletId || undefined)
        const fetchedAddresses = response.data.map((address) => ({
          ...address,
          tokens: filterSpamTokens(address.tokens)
        }))

        if (walletId) {
          // 如果指定了钱包ID，只更新该钱包的地址，保留其他钱包的地址
          // 先移除该钱包的旧地址
          const otherAddresses = this.addresses.filter(addr => addr.wallet_id !== walletId)
          // 合并其他钱包的地址和新获取的地址
          this.addresses = [...otherAddresses, ...fetchedAddresses]
        } else {
          // 如果没有指定钱包ID，替换所有地址
          this.addresses = fetchedAddresses
        }
      } catch (error: any) {
        this.error = error.message
        console.error('Failed to fetch addresses:', error)
      } finally {
        this.loading = false
      }
    },

    async createWallet(walletData: CreateWalletRequest): Promise<Wallet> {
      try {
        const response = await walletsAPI.create(walletData)
        this.wallets.push(response.data)
        return response.data
      } catch (error: any) {
        this.error = error.message
        throw error
      }
    },

    async updateWallet(walletId: number, walletData: UpdateWalletRequest): Promise<Wallet> {
      try {
        const response = await walletsAPI.update(walletId, walletData)
        const index = this.wallets.findIndex((w) => w.id === walletId)
        if (index !== -1) {
          this.wallets[index] = response.data
        }
        return response.data
      } catch (error: any) {
        this.error = error.message
        throw error
      }
    },

    async createAddress(addressData: CreateAddressRequest): Promise<Address> {
      try {
        // 确保 wallet_id 是数字类型
        const payload = {
          ...addressData,
          wallet_id: typeof addressData.wallet_id === 'string'
            ? parseInt(addressData.wallet_id, 10)
            : addressData.wallet_id
        }
        const response = await addressesAPI.create(payload)
        this.addresses.push(response.data)
        return response.data
      } catch (error: any) {
        this.error = error.message
        throw error
      }
    },

    async updateAddress(addressId: number, addressData: UpdateAddressRequest): Promise<Address> {
      try {
        const response = await addressesAPI.update(addressId, addressData)
        const index = this.addresses.findIndex((a) => a.id === addressId)
        if (index !== -1) {
          // 使用 splice 确保 Vue 能检测到数组变化
          this.addresses.splice(index, 1, {
            ...response.data,
            tokens: filterSpamTokens(response.data.tokens)
          })
        }
        return response.data
      } catch (error: any) {
        this.error = error.message
        throw error
      }
    },

    async deleteWallet(walletId: number): Promise<void> {
      try {
        await walletsAPI.delete(walletId)
        this.wallets = this.wallets.filter((w) => w.id !== walletId)
        this.addresses = this.addresses.filter((a) => a.wallet_id !== walletId)
      } catch (error: any) {
        this.error = error.message
        throw error
      }
    },

    async deleteAddress(addressId: number): Promise<void> {
      try {
        await addressesAPI.delete(addressId)
        this.addresses = this.addresses.filter((a) => a.id !== addressId)
      } catch (error: any) {
        this.error = error.message
        throw error
      }
    },

    async refreshWallet(walletId: number): Promise<void> {
      this.loading = true
      try {
        await walletsAPI.refresh(walletId)
        await this.fetchAddresses(walletId)
      } catch (error: any) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async refreshAddress(addressId: number): Promise<void> {
      this.loading = true
      try {
        const response = await addressesAPI.refresh(addressId)
        const index = this.addresses.findIndex((a) => a.id === addressId)
        if (index !== -1) {
          // 使用 splice 确保 Vue 能检测到数组变化
          this.addresses.splice(index, 1, {
            ...response.data,
            tokens: filterSpamTokens(response.data.tokens)
          })
        }
      } catch (error: any) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    selectWallet(wallet: Wallet): void {
      this.selectedWallet = wallet
    }
  }
})
