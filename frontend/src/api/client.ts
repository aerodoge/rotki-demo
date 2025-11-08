import axios, { AxiosInstance, AxiosResponse } from 'axios'
import type {
  Wallet,
  Address,
  Chain,
  RPCNode,
  GroupedRPCNodes,
  CreateWalletRequest,
  UpdateWalletRequest,
  CreateAddressRequest,
  UpdateAddressRequest,
  CreateRPCNodeRequest,
  UpdateRPCNodeRequest
} from '@/types'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 钱包 API
export const walletsAPI = {
  list: (): Promise<AxiosResponse<Wallet[]>> => apiClient.get('/wallets'),
  get: (id: number): Promise<AxiosResponse<Wallet>> => apiClient.get(`/wallets/${id}`),
  create: (data: CreateWalletRequest): Promise<AxiosResponse<Wallet>> =>
    apiClient.post('/wallets', data),
  update: (id: number, data: UpdateWalletRequest): Promise<AxiosResponse<Wallet>> =>
    apiClient.put(`/wallets/${id}`, data),
  delete: (id: number): Promise<AxiosResponse<void>> => apiClient.delete(`/wallets/${id}`),
  refresh: (id: number): Promise<AxiosResponse<void>> => apiClient.post(`/wallets/${id}/refresh`)
}

// 地址 API
export const addressesAPI = {
  list: (walletId?: number): Promise<AxiosResponse<Address[]>> => {
    const params = walletId ? { wallet_id: walletId } : {}
    return apiClient.get('/addresses', { params })
  },
  get: (id: number): Promise<AxiosResponse<Address>> => apiClient.get(`/addresses/${id}`),
  create: (data: CreateAddressRequest): Promise<AxiosResponse<Address>> =>
    apiClient.post('/addresses', data),
  update: (id: number, data: UpdateAddressRequest): Promise<AxiosResponse<Address>> =>
    apiClient.put(`/addresses/${id}`, data),
  delete: (id: number): Promise<AxiosResponse<void>> => apiClient.delete(`/addresses/${id}`),
  refresh: (id: number): Promise<AxiosResponse<Address>> =>
    apiClient.post(`/addresses/${id}/refresh`)
}

// 链 API
export const chainsAPI = {
  list: (): Promise<AxiosResponse<Chain[]>> => apiClient.get('/chains')
}

// RPC节点 API
export const rpcNodesAPI = {
  list: (chainId?: string): Promise<AxiosResponse<RPCNode[]>> => {
    const params = chainId ? { chain_id: chainId } : {}
    return apiClient.get('/rpc-nodes', { params })
  },
  grouped: (): Promise<AxiosResponse<GroupedRPCNodes>> => apiClient.get('/rpc-nodes/grouped'),
  get: (id: number): Promise<AxiosResponse<RPCNode>> => apiClient.get(`/rpc-nodes/${id}`),
  create: (data: CreateRPCNodeRequest): Promise<AxiosResponse<RPCNode>> =>
    apiClient.post('/rpc-nodes', data),
  update: (id: number, data: UpdateRPCNodeRequest): Promise<AxiosResponse<RPCNode>> =>
    apiClient.put(`/rpc-nodes/${id}`, data),
  delete: (id: number): Promise<AxiosResponse<void>> => apiClient.delete(`/rpc-nodes/${id}`),
  checkConnection: (id: number): Promise<AxiosResponse<{ connected: boolean }>> =>
    apiClient.post(`/rpc-nodes/${id}/check`),
  checkAllConnections: (): Promise<AxiosResponse<void>> => apiClient.post('/rpc-nodes/check-all')
}

export default apiClient
