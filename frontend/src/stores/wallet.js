import { defineStore } from 'pinia'
import { walletsAPI, addressesAPI } from '../api/client'

// Helper function to filter spam tokens
const filterSpamTokens = (tokens) => {
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

    // Check for spam keywords
    for (const keyword of spamKeywords) {
      if (symbolLower.includes(keyword) || nameLower.includes(keyword)) {
        return false
      }
    }

    // Check for suspicious patterns in tokens with $0 value
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

export const useWalletStore = defineStore('wallet', {
  state: () => ({
    wallets: [],
    addresses: [],
    selectedWallet: null,
    loading: false,
    error: null
  }),

  getters: {
    getAddressesByWallet: (state) => (walletId) => {
      return state.addresses.filter((addr) => addr.wallet_id === walletId)
    },

    getTotalValueByWallet: (state) => (walletId) => {
      const addresses = state.addresses.filter((addr) => addr.wallet_id === walletId)
      return addresses.reduce((sum, addr) => {
        const tokenValue = addr.tokens?.reduce((s, t) => s + (t.usd_value || 0), 0) || 0
        return sum + tokenValue
      }, 0)
    },

    getTotalValue: (state) => {
      return state.addresses.reduce((sum, addr) => {
        const tokenValue = addr.tokens?.reduce((s, t) => s + (t.usd_value || 0), 0) || 0
        return sum + tokenValue
      }, 0)
    }
  },

  actions: {
    async fetchWallets() {
      this.loading = true
      this.error = null
      try {
        const response = await walletsAPI.list()
        this.wallets = response.data
      } catch (error) {
        this.error = error.message
        console.error('Failed to fetch wallets:', error)
      } finally {
        this.loading = false
      }
    },

    async fetchAddresses(walletId = null) {
      this.loading = true
      this.error = null
      try {
        const response = await addressesAPI.list(walletId)
        // Filter spam tokens from all addresses
        this.addresses = response.data.map((address) => ({
          ...address,
          tokens: filterSpamTokens(address.tokens)
        }))
      } catch (error) {
        this.error = error.message
        console.error('Failed to fetch addresses:', error)
      } finally {
        this.loading = false
      }
    },

    async createWallet(walletData) {
      try {
        const response = await walletsAPI.create(walletData)
        this.wallets.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async updateWallet(walletId, walletData) {
      try {
        const response = await walletsAPI.update(walletId, walletData)
        const index = this.wallets.findIndex((w) => w.id === walletId)
        if (index !== -1) {
          this.wallets[index] = response.data
        }
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async createAddress(addressData) {
      try {
        const response = await addressesAPI.create(addressData)
        this.addresses.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async updateAddress(addressId, addressData) {
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
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async deleteWallet(walletId) {
      try {
        await walletsAPI.delete(walletId)
        this.wallets = this.wallets.filter((w) => w.id !== walletId)
        this.addresses = this.addresses.filter((a) => a.wallet_id !== walletId)
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async deleteAddress(addressId) {
      try {
        await addressesAPI.delete(addressId)
        this.addresses = this.addresses.filter((a) => a.id !== addressId)
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async refreshWallet(walletId) {
      this.loading = true
      try {
        await walletsAPI.refresh(walletId)
        await this.fetchAddresses(walletId)
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async refreshAddress(addressId) {
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
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    selectWallet(wallet) {
      this.selectedWallet = wallet
    }
  }
})
