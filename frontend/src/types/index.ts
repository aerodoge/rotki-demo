// API响应类型
export interface Chain {
  id: string
  name: string
  logo_url?: string
  rpc_url?: string
  explorer_url?: string
}

export interface Token {
  id?: number
  address_id?: number
  chain_id: string
  chain?: Chain
  token_id: string
  name: string
  symbol: string
  decimals: number
  balance: string
  price: number
  usd_value: number
  logo_url?: string
}

export interface Address {
  id: number
  wallet_id: number
  chain_id: string
  address: string
  label?: string
  created_at?: string
  updated_at?: string
  tokens?: Token[]
}

export interface Wallet {
  id: number
  name: string
  description?: string
  created_at?: string
  updated_at?: string
}

export interface RPCNode {
  id: number
  chain_id: string
  name: string
  url: string
  is_active: boolean
  priority: number
  created_at?: string
  updated_at?: string
}

export interface GroupedRPCNodes {
  [chainId: string]: RPCNode[]
}

// API请求类型
export interface CreateWalletRequest {
  name: string
  description?: string
}

export interface UpdateWalletRequest {
  name?: string
  description?: string
}

export interface CreateAddressRequest {
  wallet_id: number
  chain_id: string
  address: string
  label?: string
}

export interface UpdateAddressRequest {
  label?: string
}

export interface CreateRPCNodeRequest {
  chain_id: string
  name: string
  url: string
  is_active?: boolean
  priority?: number
}

export interface UpdateRPCNodeRequest {
  name?: string
  url?: string
  is_active?: boolean
  priority?: number
}
