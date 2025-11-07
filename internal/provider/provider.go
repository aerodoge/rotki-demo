package provider

import (
	"context"
	"time"
)

// DataProvider 定义区块链数据提供者的接口
// 此抽象允许在不同数据源之间切换（DeBank、自查询等）
type DataProvider interface {
	// GetTotalBalance 返回地址在所有链上的总 USD 价值
	GetTotalBalance(ctx context.Context, address string) (*TotalBalanceResponse, error)

	// GetTokenList 返回地址在特定链上的所有代币
	GetTokenList(ctx context.Context, address string, chainIDs []string) ([]TokenInfo, error)

	// GetUsedChainList 返回地址有活动的链
	GetUsedChainList(ctx context.Context, address string) ([]ChainInfo, error)

	// GetProtocolList 返回地址的 DeFi 协议持仓
	GetProtocolList(ctx context.Context, address string, chainIDs []string) ([]ProtocolInfo, error)

	// GetName 返回提供者名称（例如 "debank"、"self-query"）
	GetName() string
}

// TotalBalanceResponse 表示所有链的总余额
type TotalBalanceResponse struct {
	TotalUSDValue float64        `json:"total_usd_value"`
	ChainList     []ChainBalance `json:"chain_list"`
}

// ChainBalance 表示特定链上的余额
type ChainBalance struct {
	ChainID  string  `json:"chain_id"`
	USDValue float64 `json:"usd_value"`
}

// TokenInfo 表示带有余额和价值的代币
type TokenInfo struct {
	ChainID    string    `json:"chain_id"`
	TokenID    string    `json:"token_id"`
	Address    string    `json:"address"`
	Symbol     string    `json:"symbol"`
	Name       string    `json:"name"`
	Decimals   int       `json:"decimals"`
	LogoURL    string    `json:"logo_url"`
	Balance    string    `json:"balance"`     // 原始余额作为字符串以避免精度损失
	RawBalance string    `json:"raw_balance"` // 不带小数的原始余额
	Price      float64   `json:"price"`
	USDValue   float64   `json:"usd_value"`
	IsCore     bool      `json:"is_core"`
	IsVerified bool      `json:"is_verified"`
	IsWallet   bool      `json:"is_wallet"`
	TimeAt     time.Time `json:"time_at"`
}

// ChainInfo 表示区块链信息
type ChainInfo struct {
	ChainID       string    `json:"chain_id"`
	Name          string    `json:"name"`
	LogoURL       string    `json:"logo_url"`
	NativeTokenID string    `json:"native_token_id"`
	BornAt        time.Time `json:"born_at,omitempty"`
}

// ProtocolInfo 表示 DeFi 协议持仓
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

// PortfolioItem 表示协议中的持仓
type PortfolioItem struct {
	Name          string        `json:"name"`
	PositionType  string        `json:"position_type"` // deposit、borrow、stake 等
	NetUSDValue   float64       `json:"net_usd_value"`
	AssetUSDValue float64       `json:"asset_usd_value"`
	DebtUSDValue  float64       `json:"debt_usd_value"`
	TokenList     []TokenDetail `json:"token_list"`
}

// TokenDetail 表示持仓中的详细代币信息
type TokenDetail struct {
	Symbol   string  `json:"symbol"`
	Amount   float64 `json:"amount"`
	Price    float64 `json:"price"`
	USDValue float64 `json:"usd_value"`
}

// ProviderFactory 创建数据提供者实例
type ProviderFactory interface {
	CreateProvider(providerType string) (DataProvider, error)
}
