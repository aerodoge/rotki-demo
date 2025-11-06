package provider

import (
	"context"
	"time"
)

// DataProvider defines the interface for blockchain data providers
// This abstraction allows switching between different data sources (DeBank, self-query, etc.)
type DataProvider interface {
	// GetTotalBalance returns the total USD value across all chains for an address
	GetTotalBalance(ctx context.Context, address string) (*TotalBalanceResponse, error)

	// GetTokenList returns all tokens for an address on specific chains
	GetTokenList(ctx context.Context, address string, chainIDs []string) ([]TokenInfo, error)

	// GetUsedChainList returns chains that the address has activity on
	GetUsedChainList(ctx context.Context, address string) ([]ChainInfo, error)

	// GetProtocolList returns DeFi protocol positions for an address
	GetProtocolList(ctx context.Context, address string, chainIDs []string) ([]ProtocolInfo, error)

	// GetName returns the provider name (e.g., "debank", "self-query")
	GetName() string
}

// TotalBalanceResponse represents total balance across all chains
type TotalBalanceResponse struct {
	TotalUSDValue float64        `json:"total_usd_value"`
	ChainList     []ChainBalance `json:"chain_list"`
}

// ChainBalance represents balance on a specific chain
type ChainBalance struct {
	ChainID  string  `json:"chain_id"`
	USDValue float64 `json:"usd_value"`
}

// TokenInfo represents a token with its balance and value
type TokenInfo struct {
	ChainID    string    `json:"chain_id"`
	TokenID    string    `json:"token_id"`
	Address    string    `json:"address"`
	Symbol     string    `json:"symbol"`
	Name       string    `json:"name"`
	Decimals   int       `json:"decimals"`
	LogoURL    string    `json:"logo_url"`
	Balance    string    `json:"balance"`     // Raw balance as string to avoid precision loss
	RawBalance string    `json:"raw_balance"` // Raw balance without decimals
	Price      float64   `json:"price"`
	USDValue   float64   `json:"usd_value"`
	IsCore     bool      `json:"is_core"`
	IsVerified bool      `json:"is_verified"`
	IsWallet   bool      `json:"is_wallet"`
	TimeAt     time.Time `json:"time_at"`
}

// ChainInfo represents blockchain information
type ChainInfo struct {
	ChainID       string    `json:"chain_id"`
	Name          string    `json:"name"`
	LogoURL       string    `json:"logo_url"`
	NativeTokenID string    `json:"native_token_id"`
	BornAt        time.Time `json:"born_at,omitempty"`
}

// ProtocolInfo represents DeFi protocol position
type ProtocolInfo struct {
	ProtocolID     string          `json:"protocol_id"`
	Name           string          `json:"name"`
	SiteURL        string          `json:"site_url"`
	LogoURL        string          `json:"logo_url"`
	ChainID        string          `json:"chain_id"`
	NetUSDValue    float64         `json:"net_usd_value"`
	AssetUSDValue  float64         `json:"asset_usd_value"`
	DebtUSDValue   float64         `json:"debt_usd_value"`
	PortfolioItems []PortfolioItem `json:"portfolio_items"`
}

// PortfolioItem represents a position in a protocol
type PortfolioItem struct {
	Name          string        `json:"name"`
	PositionType  string        `json:"position_type"` // deposit, borrow, stake, etc.
	NetUSDValue   float64       `json:"net_usd_value"`
	AssetUSDValue float64       `json:"asset_usd_value"`
	DebtUSDValue  float64       `json:"debt_usd_value"`
	TokenList     []TokenDetail `json:"token_list"`
}

// TokenDetail represents detailed token information in a position
type TokenDetail struct {
	Symbol   string  `json:"symbol"`
	Amount   float64 `json:"amount"`
	Price    float64 `json:"price"`
	USDValue float64 `json:"usd_value"`
}

// ProviderFactory creates data provider instances
type ProviderFactory interface {
	CreateProvider(providerType string) (DataProvider, error)
}
