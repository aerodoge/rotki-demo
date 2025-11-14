package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// StringSlice 是用于将字符串数组存储为 JSON 的自定义类型
type StringSlice []string

// Scan 实现 sql.Scanner 接口
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

// Value 实现 driver.Valuer 接口
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil || len(s) == 0 {
		return nil, nil
	}
	return json.Marshal(s)
}

// Wallet 表示可以包含多个地址的钱包
type Wallet struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	Name          string      `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description   string      `gorm:"type:text" json:"description"`
	Tags          StringSlice `gorm:"type:json" json:"tags"`                                     // 用户定义的标签
	EnabledChains StringSlice `gorm:"type:json" json:"enabled_chains"`                           // 启用的链 ID 列表，空表示所有链
	Status        string      `gorm:"type:varchar(20);not null;default:'Enabled'" json:"status"` // Enabled 或 Disabled
	Addresses     []Address   `gorm:"foreignKey:WalletID" json:"addresses,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// Address 表示区块链地址
type Address struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	WalletID     uint        `gorm:"not null;index" json:"wallet_id"`
	Address      string      `gorm:"type:varchar(255);not null;index" json:"address"`
	ChainType    string      `gorm:"type:varchar(50);not null;default:'EVM'" json:"chain_type"`
	Label        string      `gorm:"type:varchar(255)" json:"label"`
	Tags         StringSlice `gorm:"type:json" json:"tags"` // 用户定义的标签
	LastSyncedAt *time.Time  `json:"last_synced_at,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`

	// 关系
	Wallet         *Wallet         `gorm:"foreignKey:WalletID" json:"wallet,omitempty"`
	Tokens         []Token         `gorm:"foreignKey:AddressID" json:"tokens,omitempty"`
	Protocols      []Protocol      `gorm:"foreignKey:AddressID" json:"protocols,omitempty"`
	AssetSnapshots []AssetSnapshot `gorm:"foreignKey:AddressID" json:"asset_snapshots,omitempty"`
}

// JSONMap 是用于存储 JSON 数据的自定义类型
type JSONMap map[string]interface{}

// Scan 实现 sql.Scanner 接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现 driver.Valuer 接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// AssetSnapshot 存储资产的定期快照
type AssetSnapshot struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AddressID     uint      `gorm:"not null;index:idx_address_time" json:"address_id"`
	SnapshotTime  time.Time `gorm:"not null;index:idx_address_time;index" json:"snapshot_time"`
	TotalUSDValue float64   `gorm:"type:decimal(30,6)" json:"total_usd_value"`
	DataSource    string    `gorm:"not null;default:'debank'" json:"data_source"`
	RawData       JSONMap   `gorm:"type:json" json:"raw_data,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// Chain 表示区块链网络
type Chain struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	ChainType     string    `gorm:"not null;default:'EVM';index" json:"chain_type"`
	LogoURL       string    `json:"logo_url"`
	NativeTokenID string    `json:"native_token_id"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Token 表示地址的代币余额
type Token struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AddressID   uint      `gorm:"not null;index;uniqueIndex:uk_address_chain_token" json:"address_id"`
	ChainID     string    `gorm:"not null;index;uniqueIndex:uk_address_chain_token" json:"chain_id"`
	TokenID     string    `gorm:"not null;uniqueIndex:uk_address_chain_token" json:"token_id"`
	Symbol      string    `json:"symbol"`
	Name        string    `json:"name"`
	Decimals    int       `json:"decimals"`
	LogoURL     string    `json:"logo_url"`
	Balance     string    `gorm:"type:decimal(40,18)" json:"balance"` // 可以是负数（debt代币）
	Price       float64   `gorm:"type:decimal(30,6)" json:"price"`
	USDValue    float64   `gorm:"type:decimal(30,6)" json:"usd_value"`                  // 可以是负数
	ProtocolID  string    `gorm:"type:varchar(100);index" json:"protocol_id,omitempty"` // 如果来自协议，记录协议ID
	IsDebt      bool      `gorm:"default:false" json:"is_debt"`                         // 是否是债务代币
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`

	// 关系
	Address *Address `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	Chain   *Chain   `gorm:"foreignKey:ChainID" json:"chain,omitempty"`
}

// SyncJob 跟踪后台同步操作
type SyncJob struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	AddressID    *uint      `gorm:"index" json:"address_id,omitempty"`
	WalletID     *uint      `gorm:"index" json:"wallet_id,omitempty"`
	JobType      string     `gorm:"not null" json:"job_type"`     // full_sync、token_sync、protocol_sync
	Status       string     `gorm:"not null;index" json:"status"` // pending、running、completed、failed
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	ErrorMessage string     `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// RPCNode 表示特定链的 RPC 节点配置
type RPCNode struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ChainID     string     `gorm:"type:varchar(50);not null;index" json:"chain_id"` // eth、bsc、polygon 等
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`          // 显示名称，如 "0xRPC"、"PublicNode"
	URL         string     `gorm:"type:varchar(500);not null" json:"url"`           // RPC 端点 URL
	Weight      int        `gorm:"not null;default:100" json:"weight"`              // 负载均衡权重（0-100）
	IsEnabled   bool       `gorm:"not null;default:true" json:"is_enabled"`         // 此节点是否处于活动状态
	IsConnected bool       `gorm:"not null;default:false" json:"is_connected"`      // 连接状态
	LastChecked *time.Time `json:"last_checked,omitempty"`                          // 最后连接检查时间
	Priority    int        `gorm:"not null;default:0" json:"priority"`              // 优先级更高的节点优先
	Timeout     int        `gorm:"not null;default:30" json:"timeout"`              // 请求超时时间（秒）
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关系
	Chain *Chain `gorm:"foreignKey:ChainID" json:"chain,omitempty"`
}

// Protocol 表示 DeFi 协议持仓
type Protocol struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AddressID     uint      `gorm:"not null;index;uniqueIndex:uk_address_protocol" json:"address_id"`
	ProtocolID    string    `gorm:"not null;uniqueIndex:uk_address_protocol" json:"protocol_id"`
	Name          string    `json:"name"`
	SiteURL       string    `json:"site_url"`
	LogoURL       string    `json:"logo_url"`
	ChainID       string    `gorm:"not null;index" json:"chain_id"`
	NetUSDValue   float64   `gorm:"type:decimal(30,6)" json:"net_usd_value"`
	AssetUSDValue float64   `gorm:"type:decimal(30,6)" json:"asset_usd_value"`
	DebtUSDValue  float64   `gorm:"type:decimal(30,6)" json:"debt_usd_value"`
	PositionType  string    `json:"position_type"` // lending, staking, liquidity, etc.
	RawData       JSONMap   `gorm:"type:json" json:"raw_data,omitempty"`
	LastUpdated   time.Time `gorm:"autoUpdateTime" json:"last_updated"`

	// 关系
	Address *Address `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	Chain   *Chain   `gorm:"foreignKey:ChainID" json:"chain,omitempty"`
}

// TableName 覆盖表名
func (Wallet) TableName() string        { return "wallets" }
func (Address) TableName() string       { return "addresses" }
func (AssetSnapshot) TableName() string { return "asset_snapshots" }
func (Chain) TableName() string         { return "chains" }
func (Token) TableName() string         { return "tokens" }
func (SyncJob) TableName() string       { return "sync_jobs" }
func (RPCNode) TableName() string       { return "rpc_nodes" }
func (Protocol) TableName() string      { return "protocols" }
