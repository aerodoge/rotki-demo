package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// StringSlice is a custom type for storing string arrays as JSON
type StringSlice []string

// Scan implements the sql.Scanner interface
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

// Value implements the driver.Valuer interface
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil || len(s) == 0 {
		return nil, nil
	}
	return json.Marshal(s)
}

// Wallet represents a wallet that can contain multiple addresses
type Wallet struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	Name          string      `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description   string      `gorm:"type:text" json:"description"`
	Tags          StringSlice `gorm:"type:json" json:"tags"`           // User-defined tags
	EnabledChains StringSlice `gorm:"type:json" json:"enabled_chains"` // List of enabled chain IDs, empty means all chains
	Addresses     []Address   `gorm:"foreignKey:WalletID" json:"addresses,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// Address represents a blockchain address
type Address struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	WalletID     uint        `gorm:"not null;index" json:"wallet_id"`
	Address      string      `gorm:"type:varchar(255);not null;index" json:"address"`
	ChainType    string      `gorm:"type:varchar(50);not null;default:'EVM'" json:"chain_type"`
	Label        string      `gorm:"type:varchar(255)" json:"label"`
	Tags         StringSlice `gorm:"type:json" json:"tags"` // User-defined tags
	LastSyncedAt *time.Time  `json:"last_synced_at,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`

	// Relations
	Wallet         *Wallet         `gorm:"foreignKey:WalletID" json:"wallet,omitempty"`
	Tokens         []Token         `gorm:"foreignKey:AddressID" json:"tokens,omitempty"`
	AssetSnapshots []AssetSnapshot `gorm:"foreignKey:AddressID" json:"asset_snapshots,omitempty"`
}

// JSONMap is a custom type for storing JSON data
type JSONMap map[string]interface{}

// Scan implements the sql.Scanner interface
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

// Value implements the driver.Valuer interface
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// AssetSnapshot stores periodic snapshots of assets
type AssetSnapshot struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AddressID     uint      `gorm:"not null;index:idx_address_time" json:"address_id"`
	SnapshotTime  time.Time `gorm:"not null;index:idx_address_time;index" json:"snapshot_time"`
	TotalUSDValue float64   `gorm:"type:decimal(30,6)" json:"total_usd_value"`
	DataSource    string    `gorm:"not null;default:'debank'" json:"data_source"`
	RawData       JSONMap   `gorm:"type:json" json:"raw_data,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// Chain represents a blockchain network
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

// Token represents a token balance for an address
type Token struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AddressID   uint      `gorm:"not null;index;uniqueIndex:uk_address_chain_token" json:"address_id"`
	ChainID     string    `gorm:"not null;index;uniqueIndex:uk_address_chain_token" json:"chain_id"`
	TokenID     string    `gorm:"not null;uniqueIndex:uk_address_chain_token" json:"token_id"`
	Symbol      string    `json:"symbol"`
	Name        string    `json:"name"`
	Decimals    int       `json:"decimals"`
	LogoURL     string    `json:"logo_url"`
	Balance     string    `gorm:"type:decimal(40,18)" json:"balance"`
	Price       float64   `gorm:"type:decimal(30,6)" json:"price"`
	USDValue    float64   `gorm:"type:decimal(30,6)" json:"usd_value"`
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`

	// Relations
	Address *Address `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	Chain   *Chain   `gorm:"foreignKey:ChainID" json:"chain,omitempty"`
}

// SyncJob tracks background sync operations
type SyncJob struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	AddressID    *uint      `gorm:"index" json:"address_id,omitempty"`
	WalletID     *uint      `gorm:"index" json:"wallet_id,omitempty"`
	JobType      string     `gorm:"not null" json:"job_type"`     // full_sync, token_sync, protocol_sync
	Status       string     `gorm:"not null;index" json:"status"` // pending, running, completed, failed
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	ErrorMessage string     `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// RPCNode represents an RPC node configuration for a specific chain
type RPCNode struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ChainID     string     `gorm:"type:varchar(50);not null;index" json:"chain_id"` // eth, bsc, polygon, etc.
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`          // Display name like "0xRPC", "PublicNode"
	URL         string     `gorm:"type:varchar(500);not null" json:"url"`           // RPC endpoint URL
	Weight      int        `gorm:"not null;default:100" json:"weight"`              // Weight for load balancing (0-100)
	IsEnabled   bool       `gorm:"not null;default:true" json:"is_enabled"`         // Whether this node is active
	IsConnected bool       `gorm:"not null;default:false" json:"is_connected"`      // Connection status
	LastChecked *time.Time `json:"last_checked,omitempty"`                          // Last connectivity check time
	Priority    int        `gorm:"not null;default:0" json:"priority"`              // Higher priority nodes are preferred
	Timeout     int        `gorm:"not null;default:30" json:"timeout"`              // Request timeout in seconds
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	Chain *Chain `gorm:"foreignKey:ChainID" json:"chain,omitempty"`
}

// TableName overrides
func (Wallet) TableName() string        { return "wallets" }
func (Address) TableName() string       { return "addresses" }
func (AssetSnapshot) TableName() string { return "asset_snapshots" }
func (Chain) TableName() string         { return "chains" }
func (Token) TableName() string         { return "tokens" }
func (SyncJob) TableName() string       { return "sync_jobs" }
func (RPCNode) TableName() string       { return "rpc_nodes" }
