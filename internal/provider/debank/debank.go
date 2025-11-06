package debank

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/miles/rotki-demo/internal/config"
	"github.com/miles/rotki-demo/internal/logger"
	"github.com/miles/rotki-demo/internal/provider"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// DeBankProvider implements the DataProvider interface using DeBank API
type DeBankProvider struct {
	config      *config.DeBankConfig
	httpClient  *http.Client
	rateLimiter *rate.Limiter
}

// NewDeBankProvider creates a new DeBank provider instance
func NewDeBankProvider(cfg *config.DeBankConfig) *DeBankProvider {
	// Create rate limiter: requests per second with burst capacity
	limiter := rate.NewLimiter(
		rate.Limit(cfg.RateLimit.RequestsPerSecond),
		cfg.RateLimit.Burst,
	)

	return &DeBankProvider{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.GetTimeout(),
		},
		rateLimiter: limiter,
	}
}

// GetName returns the provider name
func (d *DeBankProvider) GetName() string {
	return "debank"
}

// doRequest performs an HTTP request with rate limiting
func (d *DeBankProvider) doRequest(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	// Wait for rate limiter
	if err := d.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// Build URL with query parameters
	url := d.config.BaseURL + path
	if len(params) > 0 {
		queryParts := make([]string, 0, len(params))
		for k, v := range params {
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", k, v))
		}
		url = url + "?" + strings.Join(queryParts, "&")
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("AccessKey", d.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Log request
	logger.Debug("DeBank API request",
		zap.String("url", url),
		zap.Any("params", params),
	)

	// Execute request
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		logger.Error("DeBank API error",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(body)),
		)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetTotalBalance returns the total balance for an address
func (d *DeBankProvider) GetTotalBalance(ctx context.Context, address string) (*provider.TotalBalanceResponse, error) {
	body, err := d.doRequest(ctx, "/v1/user/total_balance", map[string]string{
		"id": address,
	})
	if err != nil {
		return nil, err
	}

	var response struct {
		TotalUSDValue float64 `json:"total_usd_value"`
		ChainList     []struct {
			ChainID  string  `json:"id"`
			USDValue float64 `json:"usd_value"`
		} `json:"chain_list"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	result := &provider.TotalBalanceResponse{
		TotalUSDValue: response.TotalUSDValue,
		ChainList:     make([]provider.ChainBalance, len(response.ChainList)),
	}

	for i, chain := range response.ChainList {
		result.ChainList[i] = provider.ChainBalance{
			ChainID:  chain.ChainID,
			USDValue: chain.USDValue,
		}
	}

	return result, nil
}

// GetTokenList returns token list for an address
func (d *DeBankProvider) GetTokenList(ctx context.Context, address string, chainIDs []string) ([]provider.TokenInfo, error) {
	params := map[string]string{
		"id":     address,
		"is_all": "true",
	}

	// If specific chains requested, use them
	if len(chainIDs) > 0 {
		params["chain_ids"] = strings.Join(chainIDs, ",")
	}

	body, err := d.doRequest(ctx, "/v1/user/all_token_list", params)
	if err != nil {
		return nil, err
	}

	// Custom type to handle both string and number for raw_amount
	type FlexibleString struct {
		Value string
	}

	var tokens []struct {
		ChainID    string          `json:"chain"`
		ID         string          `json:"id"`
		Symbol     string          `json:"symbol"`
		Name       string          `json:"name"`
		Decimals   int             `json:"decimals"`
		LogoURL    string          `json:"logo_url"`
		Amount     float64         `json:"amount"`
		RawAmount  json.RawMessage `json:"raw_amount"` // Use RawMessage to handle both string and number
		Price      float64         `json:"price"`
		IsCore     bool            `json:"is_core"`
		IsVerified bool            `json:"is_verified"`
		IsWallet   bool            `json:"is_wallet"`
		TimeAt     float64         `json:"time_at"`
	}

	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	result := make([]provider.TokenInfo, len(tokens))
	for i, token := range tokens {
		// Parse raw_amount which can be either string or number
		var rawAmount string
		if len(token.RawAmount) > 0 {
			// Try to unmarshal as string first
			if err := json.Unmarshal(token.RawAmount, &rawAmount); err != nil {
				// If failed, it might be a number, convert it
				var numAmount float64
				if err := json.Unmarshal(token.RawAmount, &numAmount); err == nil {
					rawAmount = fmt.Sprintf("%.0f", numAmount)
				}
			}
		}

		result[i] = provider.TokenInfo{
			ChainID:    token.ChainID,
			TokenID:    token.ID,
			Address:    token.ID,
			Symbol:     token.Symbol,
			Name:       token.Name,
			Decimals:   token.Decimals,
			LogoURL:    token.LogoURL,
			Balance:    fmt.Sprintf("%f", token.Amount),
			RawBalance: rawAmount,
			Price:      token.Price,
			USDValue:   token.Amount * token.Price,
			IsCore:     token.IsCore,
			IsVerified: token.IsVerified,
			IsWallet:   token.IsWallet,
			TimeAt:     time.Unix(int64(token.TimeAt), 0),
		}
	}

	return result, nil
}

// GetUsedChainList returns chains used by an address
func (d *DeBankProvider) GetUsedChainList(ctx context.Context, address string) ([]provider.ChainInfo, error) {
	body, err := d.doRequest(ctx, "/v1/user/used_chain_list", map[string]string{
		"id": address,
	})
	if err != nil {
		return nil, err
	}

	var chains []struct {
		ChainID       string `json:"id"`
		Name          string `json:"name"`
		LogoURL       string `json:"logo_url"`
		NativeTokenID string `json:"native_token_id"`
	}

	if err := json.Unmarshal(body, &chains); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	result := make([]provider.ChainInfo, len(chains))
	for i, chain := range chains {
		result[i] = provider.ChainInfo{
			ChainID:       chain.ChainID,
			Name:          chain.Name,
			LogoURL:       chain.LogoURL,
			NativeTokenID: chain.NativeTokenID,
		}
	}

	return result, nil
}

// GetProtocolList returns DeFi protocol positions
func (d *DeBankProvider) GetProtocolList(ctx context.Context, address string, chainIDs []string) ([]provider.ProtocolInfo, error) {
	params := map[string]string{
		"id": address,
	}

	if len(chainIDs) > 0 {
		params["chain_ids"] = strings.Join(chainIDs, ",")
	}

	body, err := d.doRequest(ctx, "/v1/user/all_complex_protocol_list", params)
	if err != nil {
		return nil, err
	}

	var protocols []struct {
		ProtocolID        string  `json:"id"`
		Name              string  `json:"name"`
		SiteURL           string  `json:"site_url"`
		LogoURL           string  `json:"logo_url"`
		Chain             string  `json:"chain"`
		NetUSDValue       float64 `json:"net_usd_value"`
		AssetUSDValue     float64 `json:"asset_usd_value"`
		DebtUSDValue      float64 `json:"debt_usd_value"`
		PortfolioItemList []struct {
			Name         string `json:"name"`
			PositionType string `json:"detail_types"`
			Stats        struct {
				NetUSDValue   float64 `json:"net_usd_value"`
				AssetUSDValue float64 `json:"asset_usd_value"`
				DebtUSDValue  float64 `json:"debt_usd_value"`
			} `json:"stats"`
		} `json:"portfolio_item_list"`
	}

	if err := json.Unmarshal(body, &protocols); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	result := make([]provider.ProtocolInfo, len(protocols))
	for i, proto := range protocols {
		portfolioItems := make([]provider.PortfolioItem, len(proto.PortfolioItemList))
		for j, item := range proto.PortfolioItemList {
			portfolioItems[j] = provider.PortfolioItem{
				Name:          item.Name,
				PositionType:  item.PositionType,
				NetUSDValue:   item.Stats.NetUSDValue,
				AssetUSDValue: item.Stats.AssetUSDValue,
				DebtUSDValue:  item.Stats.DebtUSDValue,
			}
		}

		result[i] = provider.ProtocolInfo{
			ProtocolID:     proto.ProtocolID,
			Name:           proto.Name,
			SiteURL:        proto.SiteURL,
			LogoURL:        proto.LogoURL,
			ChainID:        proto.Chain,
			NetUSDValue:    proto.NetUSDValue,
			AssetUSDValue:  proto.AssetUSDValue,
			DebtUSDValue:   proto.DebtUSDValue,
			PortfolioItems: portfolioItems,
		}
	}

	return result, nil
}
