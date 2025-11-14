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
  balance: string  // 可以是负数字符串（debt代币）
  price: number
  usd_value: number  // 可以是负数
  logo_url?: string
  protocol_id?: string  // 如果来自协议，标记协议ID
  is_debt?: boolean     // 是否是债务代币
}

export interface Protocol {
  id: number
  address_id: number
  protocol_id: string
  name: string
  site_url?: string
  logo_url?: string
  chain_id: string
  chain?: Chain
  net_usd_value: number
  asset_usd_value: number
  debt_usd_value: number
  position_type: string
  raw_data?: any
  last_updated?: string
}

export interface Address {
  id: number
  wallet_id: number
  address: string
  chain_type?: string
  label?: string
  tags?: string[]
  last_synced_at?: string
  created_at?: string
  updated_at?: string
  tokens?: Token[]
  protocols?: Protocol[]
}

export interface Wallet {
  id: number
  name: string
  description?: string
  tags?: string[]
  enabled_chains?: string[]
  status?: string
  addresses?: Address[]
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
  tags?: string[]
  enabled_chains?: string[]
  status?: string
}

export interface UpdateWalletRequest {
  name?: string
  description?: string
  tags?: string[]
  enabled_chains?: string[]
  status?: string
}

export interface CreateAddressRequest {
  wallet_id: number | string
  address: string
  chain_type?: string
  label?: string
  tags?: string[]
}

export interface UpdateAddressRequest {
  label?: string
  tags?: string[]
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
