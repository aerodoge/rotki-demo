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
        const tokenValue = addr.tokens?.reduce((s, t) => s + (t.usd_value || 0), 0) || 0
        return sum + tokenValue
      }, 0)
    },

    getTotalValue: (state): number => {
      return state.addresses.reduce((sum, addr) => {
        const tokenValue = addr.tokens?.reduce((s, t) => s + (t.usd_value || 0), 0) || 0
        return sum + tokenValue
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
        // 从所有地址过滤垃圾代币
        this.addresses = response.data.map((address) => ({
          ...address,
          tokens: filterSpamTokens(address.tokens)
        }))
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
        const response = await addressesAPI.create(addressData)
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
          this.addresses[index] = {
            ...response.data,
            tokens: filterSpamTokens(response.data.tokens)
          }
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
          this.addresses[index] = {
            ...response.data,
            tokens: filterSpamTokens(response.data.tokens)
          }
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
