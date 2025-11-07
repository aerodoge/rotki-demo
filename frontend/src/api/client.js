import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 钱包 API
export const walletsAPI = {
  list: () => apiClient.get('/wallets'),
  get: (id) => apiClient.get(`/wallets/${id}`),
  create: (data) => apiClient.post('/wallets', data),
  update: (id, data) => apiClient.put(`/wallets/${id}`, data),
  delete: (id) => apiClient.delete(`/wallets/${id}`),
  refresh: (id) => apiClient.post(`/wallets/${id}/refresh`)
}

// 地址 API
export const addressesAPI = {
  list: (walletId) => {
    const params = walletId ? { wallet_id: walletId } : {}
    return apiClient.get('/addresses', { params })
  },
  get: (id) => apiClient.get(`/addresses/${id}`),
  create: (data) => apiClient.post('/addresses', data),
  update: (id, data) => apiClient.put(`/addresses/${id}`, data),
  delete: (id) => apiClient.delete(`/addresses/${id}`),
  refresh: (id) => apiClient.post(`/addresses/${id}/refresh`)
}

// 链 API
export const chainsAPI = {
  list: () => apiClient.get('/chains')
}

// RPC节点 API
export const rpcNodesAPI = {
  list: (chainId) => {
    const params = chainId ? { chain_id: chainId } : {}
    return apiClient.get('/rpc-nodes', { params })
  },
  grouped: () => apiClient.get('/rpc-nodes/grouped'),
  get: (id) => apiClient.get(`/rpc-nodes/${id}`),
  create: (data) => apiClient.post('/rpc-nodes', data),
  update: (id, data) => apiClient.put(`/rpc-nodes/${id}`, data),
  delete: (id) => apiClient.delete(`/rpc-nodes/${id}`),
  checkConnection: (id) => apiClient.post(`/rpc-nodes/${id}/check`),
  checkAllConnections: () => apiClient.post('/rpc-nodes/check-all')
}

export default apiClient
