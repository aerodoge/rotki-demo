package debank

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rotki-demo/internal/config"
	"github.com/rotki-demo/internal/logger"
	"github.com/rotki-demo/internal/provider"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// DeBankProvider 使用 DeBank API 实现 DataProvider 接口
type DeBankProvider struct {
	config      *config.DeBankConfig
	httpClient  *http.Client
	rateLimiter *rate.Limiter
}

// NewDeBankProvider 创建一个新的 DeBank 提供者实例
func NewDeBankProvider(cfg *config.DeBankConfig) *DeBankProvider {
	// 创建速率限制器：每秒请求数和突发容量
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

// GetName 返回提供者名称
func (d *DeBankProvider) GetName() string {
	return "debank"
}

// doRequest 执行带速率限制的 HTTP 请求
func (d *DeBankProvider) doRequest(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	// 等待速率限制器
	if err := d.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// 构建带查询参数的 URL
	url := d.config.BaseURL + path
	if len(params) > 0 {
		queryParts := make([]string, 0, len(params))
		for k, v := range params {
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", k, v))
		}
		url = url + "?" + strings.Join(queryParts, "&")
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 添加请求头
	req.Header.Set("AccessKey", d.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// 记录请求
	logger.Debug("DeBank API request",
		zap.String("url", url),
		zap.Any("params", params),
	)

	// 执行请求
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		logger.Error("DeBank API error",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(body)),
		)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetTotalBalance 返回地址的总余额
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

// GetTokenList 返回地址的代币列表
func (d *DeBankProvider) GetTokenList(ctx context.Context, address string, chainIDs []string) ([]provider.TokenInfo, error) {
	params := map[string]string{
		"id":     address,
		"is_all": "false",
	}

	// 如果请求特定链，使用它们
	if len(chainIDs) > 0 {
		params["chain_ids"] = strings.Join(chainIDs, ",")
	}

	body, err := d.doRequest(ctx, "/v1/user/all_token_list", params)
	if err != nil {
		return nil, err
	}

	// 自定义类型以处理 raw_amount 的字符串和数字
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
		RawAmount  json.RawMessage `json:"raw_amount"` // 使用 RawMessage 处理字符串和数字
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
		// 解析 raw_amount，可以是字符串或数字
		var rawAmount string
		if len(token.RawAmount) > 0 {
			// 首先尝试作为字符串解析
			if err := json.Unmarshal(token.RawAmount, &rawAmount); err != nil {
				// 如果失败，可能是数字，转换它
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

// GetUsedChainList 返回地址使用的链
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

// GetProtocolList 返回 DeFi 协议持仓
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

	type debankToken struct {
		ID       string  `json:"id"`
		Chain    string  `json:"chain"`
		Name     string  `json:"name"`
		Symbol   string  `json:"symbol"`
		Decimals int     `json:"decimals"`
		LogoURL  string  `json:"logo_url"`
		Price    float64 `json:"price"`
		Amount   float64 `json:"amount"`
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
			Name         string   `json:"name"`
			PositionType []string `json:"detail_types"`
			Stats        struct {
				NetUSDValue   float64 `json:"net_usd_value"`
				AssetUSDValue float64 `json:"asset_usd_value"`
				DebtUSDValue  float64 `json:"debt_usd_value"`
			} `json:"stats"`
			AssetTokenList []debankToken `json:"asset_token_list"`
			Detail         struct {
				SupplyTokenList []debankToken `json:"supply_token_list"`
				BorrowTokenList []debankToken `json:"borrow_token_list"`
				HealthRate      float64       `json:"health_rate"`
			} `json:"detail"`
		} `json:"portfolio_item_list"`
	}

	if err := json.Unmarshal(body, &protocols); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 辅助函数：转换 DeBank token 到 provider.TokenDetail
	convertToken := func(t debankToken, chainID string, isDebt bool) provider.TokenDetail {
		amount := t.Amount
		if isDebt && amount > 0 {
			amount = -amount // 确保 debt 是负数
		}
		return provider.TokenDetail{
			TokenID:  t.ID,
			ChainID:  chainID,
			Symbol:   t.Symbol,
			Name:     t.Name,
			Decimals: t.Decimals,
			LogoURL:  t.LogoURL,
			Amount:   amount,
			Price:    t.Price,
			USDValue: amount * t.Price,
			IsDebt:   isDebt,
		}
	}

	result := make([]provider.ProtocolInfo, len(protocols))
	for i, proto := range protocols {
		portfolioItems := make([]provider.PortfolioItem, len(proto.PortfolioItemList))
		for j, item := range proto.PortfolioItemList {
			// 转换 asset_token_list（包含正负值）
			assetTokens := make([]provider.TokenDetail, len(item.AssetTokenList))
			for k, token := range item.AssetTokenList {
				assetTokens[k] = convertToken(token, proto.Chain, token.Amount < 0)
			}

			// 转换 supply_token_list
			supplyTokens := make([]provider.TokenDetail, len(item.Detail.SupplyTokenList))
			for k, token := range item.Detail.SupplyTokenList {
				supplyTokens[k] = convertToken(token, proto.Chain, false)
			}

			// 转换 borrow_token_list (这些应该是负值)
			borrowTokens := make([]provider.TokenDetail, len(item.Detail.BorrowTokenList))
			for k, token := range item.Detail.BorrowTokenList {
				borrowTokens[k] = convertToken(token, proto.Chain, true)
			}

			// 提取主要的 position type
			positionType := "unknown"
			if len(item.PositionType) > 0 {
				positionType = item.PositionType[0]
			}

			portfolioItems[j] = provider.PortfolioItem{
				Name:            item.Name,
				PositionType:    positionType,
				NetUSDValue:     item.Stats.NetUSDValue,
				AssetUSDValue:   item.Stats.AssetUSDValue,
				DebtUSDValue:    item.Stats.DebtUSDValue,
				AssetTokenList:  assetTokens,
				SupplyTokenList: supplyTokens,
				BorrowTokenList: borrowTokens,
				HealthRate:      item.Detail.HealthRate,
			}
		}

		// 计算所有 portfolio items 的总净值、资产值和债务值
		var totalNetUSD, totalAssetUSD, totalDebtUSD float64
		for _, item := range portfolioItems {
			totalNetUSD += item.NetUSDValue
			totalAssetUSD += item.AssetUSDValue
			totalDebtUSD += item.DebtUSDValue
		}

		result[i] = provider.ProtocolInfo{
			ProtocolID:     proto.ProtocolID,
			Name:           proto.Name,
			SiteURL:        proto.SiteURL,
			LogoURL:        proto.LogoURL,
			ChainID:        proto.Chain,
			NetUSDValue:    totalNetUSD,
			AssetUSDValue:  totalAssetUSD,
			DebtUSDValue:   totalDebtUSD,
			PortfolioItems: portfolioItems,
		}
	}

	return result, nil
}
