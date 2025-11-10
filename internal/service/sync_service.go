package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rotki-demo/internal/logger"
	"github.com/rotki-demo/internal/models"
	"github.com/rotki-demo/internal/provider"
	"github.com/rotki-demo/internal/repository"
	"go.uber.org/zap"
)

// SyncService 处理数据同步
type SyncService struct {
	dataProvider provider.DataProvider
	walletRepo   *repository.WalletRepository
	addressRepo  *repository.AddressRepository
	tokenRepo    *repository.TokenRepository
	protocolRepo *repository.ProtocolRepository
	chainRepo    *repository.ChainRepository
	syncInterval time.Duration
	batchSize    int
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

// NewSyncService 创建一个新的同步服务
func NewSyncService(
	dataProvider provider.DataProvider,
	walletRepo *repository.WalletRepository,
	addressRepo *repository.AddressRepository,
	tokenRepo *repository.TokenRepository,
	protocolRepo *repository.ProtocolRepository,
	chainRepo *repository.ChainRepository,
	syncInterval time.Duration,
	batchSize int,
) *SyncService {
	return &SyncService{
		dataProvider: dataProvider,
		walletRepo:   walletRepo,
		addressRepo:  addressRepo,
		tokenRepo:    tokenRepo,
		protocolRepo: protocolRepo,
		chainRepo:    chainRepo,
		syncInterval: syncInterval,
		batchSize:    batchSize,
		stopChan:     make(chan struct{}),
	}
}

// Start 启动后台同步进程
func (s *SyncService) Start() {
	s.wg.Add(1)
	go s.syncLoop()
	logger.Info("Sync service started", zap.Duration("interval", s.syncInterval))
}

// Stop 停止后台同步进程
func (s *SyncService) Stop() {
	close(s.stopChan)
	s.wg.Wait()
	logger.Info("Sync service stopped")
}

// syncLoop 运行周期性同步
func (s *SyncService) syncLoop() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	// 运行初始同步
	s.syncAll()

	for {
		select {
		case <-ticker.C:
			s.syncAll()
		case <-s.stopChan:
			return
		}
	}
}

// syncAll 同步所有需要更新的地址
func (s *SyncService) syncAll() {
	ctx := context.Background()

	// 获取需要同步的地址
	addresses, err := s.addressRepo.GetAllNeedingSync(s.syncInterval)
	if err != nil {
		logger.Error("Failed to get addresses for sync", zap.Error(err))
		return
	}

	if len(addresses) == 0 {
		logger.Debug("No addresses need syncing")
		return
	}

	logger.Info("Starting sync", zap.Int("address_count", len(addresses)))

	// 分批处理
	for i := 0; i < len(addresses); i += s.batchSize {
		end := i + s.batchSize
		if end > len(addresses) {
			end = len(addresses)
		}

		batch := addresses[i:end]
		s.syncBatch(ctx, batch)
	}

	logger.Info("Sync completed", zap.Int("address_count", len(addresses)))
}

// syncBatch 并发同步一批地址
func (s *SyncService) syncBatch(ctx context.Context, addresses []models.Address) {
	var wg sync.WaitGroup

	for _, addr := range addresses {
		wg.Add(1)
		go func(address models.Address) {
			defer wg.Done()
			if err := s.SyncAddress(ctx, address.ID); err != nil {
				logger.Error("Failed to sync address",
					zap.Uint("address_id", address.ID),
					zap.String("address", address.Address),
					zap.Error(err),
				)
			}
		}(addr)
	}

	wg.Wait()
}

// SyncAddress 同步特定地址的数据
func (s *SyncService) SyncAddress(ctx context.Context, addressID uint) error {
	// 获取地址详情
	address, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return fmt.Errorf("failed to get address: %w", err)
	}

	// 获取钱包以检查启用的链
	wallet, err := s.walletRepo.GetByID(address.WalletID)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	logger.Debug("Syncing address",
		zap.Uint("address_id", addressID),
		zap.String("address", address.Address),
		zap.Any("enabled_chains", wallet.EnabledChains),
	)

	// 确定要查询的链
	var chainIDsToQuery []string
	if len(wallet.EnabledChains) > 0 {
		chainIDsToQuery = wallet.EnabledChains
	}

	// 首先获取并更新插入链
	chains, err := s.dataProvider.GetUsedChainList(ctx, address.Address)
	if err != nil {
		return fmt.Errorf("failed to get chain list: %w", err)
	}

	// 根据启用的链过滤链
	if len(chainIDsToQuery) > 0 {
		enabledChainMap := make(map[string]bool)
		for _, chainID := range chainIDsToQuery {
			enabledChainMap[chainID] = true
		}

		filteredChains := make([]provider.ChainInfo, 0)
		for _, chain := range chains {
			if enabledChainMap[chain.ChainID] {
				filteredChains = append(filteredChains, chain)
			}
		}
		chains = filteredChains
	}

	// 转换为数据库模型
	dbChains := make([]models.Chain, 0, len(chains))
	for _, chain := range chains {
		dbChains = append(dbChains, models.Chain{
			ID:            chain.ChainID,
			Name:          chain.Name,
			LogoURL:       chain.LogoURL,
			NativeTokenID: chain.NativeTokenID,
		})
	}

	// 首先更新插入链
	if err := s.chainRepo.UpsertBatch(dbChains); err != nil {
		return fmt.Errorf("failed to upsert chains: %w", err)
	}

	// 从提供者获取代币列表，可选按链过滤
	tokens, err := s.dataProvider.GetTokenList(ctx, address.Address, chainIDsToQuery)
	if err != nil {
		return fmt.Errorf("failed to get token list: %w", err)
	}

	// 过滤掉垃圾/欺诈代币
	filteredTokens := filterSpamTokens(tokens)

	// 转换为数据库模型
	dbTokens := make([]models.Token, 0, len(filteredTokens))
	for _, token := range filteredTokens {
		dbTokens = append(dbTokens, models.Token{
			AddressID: addressID,
			ChainID:   token.ChainID,
			TokenID:   token.TokenID,
			Symbol:    token.Symbol,
			Name:      token.Name,
			Decimals:  token.Decimals,
			LogoURL:   token.LogoURL,
			Balance:   token.Balance,
			Price:     token.Price,
			USDValue:  token.USDValue,
		})
	}

	// 先删除旧的钱包代币（保留协议代币）
	if err := s.tokenRepo.DeleteWalletTokensByAddressID(addressID); err != nil {
		return fmt.Errorf("failed to delete old wallet tokens: %w", err)
	}

	// 插入新的钱包代币
	if err := s.tokenRepo.UpsertBatch(dbTokens); err != nil {
		return fmt.Errorf("failed to upsert tokens: %w", err)
	}

	// 获取并同步协议持仓
	protocols, err := s.dataProvider.GetProtocolList(ctx, address.Address, chainIDsToQuery)
	if err != nil {
		logger.Warn("Failed to get protocol list (non-fatal)",
			zap.Uint("address_id", addressID),
			zap.Error(err),
		)
		// 协议获取失败不应该阻止整个同步过程
	} else {
		// 转换为数据库模型
		dbProtocols := make([]models.Protocol, 0, len(protocols))
		protocolTokens := make([]models.Token, 0) // 收集所有协议代币

		for _, proto := range protocols {
			// 确定主要的仓位类型
			positionType := "unknown"
			if len(proto.PortfolioItems) > 0 {
				positionType = proto.PortfolioItems[0].PositionType
			}

			// 计算所有 portfolio items 的总净值、资产值和债务值
			var totalNetUSD, totalAssetUSD, totalDebtUSD float64
			for _, item := range proto.PortfolioItems {
				totalNetUSD += item.NetUSDValue
				totalAssetUSD += item.AssetUSDValue
				totalDebtUSD += item.DebtUSDValue
			}

			// 将原始数据存储为 JSON
			rawData := make(models.JSONMap)
			rawData["portfolio_items"] = proto.PortfolioItems

			dbProtocols = append(dbProtocols, models.Protocol{
				AddressID:     addressID,
				ProtocolID:    proto.ProtocolID,
				Name:          proto.Name,
				SiteURL:       proto.SiteURL,
				LogoURL:       proto.LogoURL,
				ChainID:       proto.ChainID,
				NetUSDValue:   totalNetUSD,
				AssetUSDValue: totalAssetUSD,
				DebtUSDValue:  totalDebtUSD,
				PositionType:  positionType,
				RawData:       rawData,
			})

			// 提取协议中的代币并添加到 tokens 表
			// 使用 AssetTokenList（包含正负值的完整列表）
			for _, item := range proto.PortfolioItems {
				for _, tokenDetail := range item.AssetTokenList {
					// 构造符合 Rotki 风格的名称
					tokenName := tokenDetail.Name
					if tokenDetail.IsDebt {
						// Debt 代币可能需要特殊前缀
						if !strings.Contains(tokenName, "debt") && !strings.Contains(tokenName, "Debt") {
							tokenName = "Debt " + tokenName
						}
					}

					protocolTokens = append(protocolTokens, models.Token{
						AddressID:  addressID,
						ChainID:    tokenDetail.ChainID,
						TokenID:    tokenDetail.TokenID,
						Symbol:     tokenDetail.Symbol,
						Name:       tokenName,
						Decimals:   tokenDetail.Decimals,
						LogoURL:    tokenDetail.LogoURL,
						Balance:    fmt.Sprintf("%.18f", tokenDetail.Amount), // 保留符号
						Price:      tokenDetail.Price,
						USDValue:   tokenDetail.USDValue, // 可以是负数
						ProtocolID: proto.ProtocolID,
						IsDebt:     tokenDetail.IsDebt,
					})
				}
			}
		}

		// 更新插入协议
		if err := s.protocolRepo.UpsertBatch(dbProtocols); err != nil {
			return fmt.Errorf("failed to upsert protocols: %w", err)
		}

		// 删除旧的协议代币
		if err := s.tokenRepo.DeleteProtocolTokensByAddressID(addressID); err != nil {
			return fmt.Errorf("failed to delete old protocol tokens: %w", err)
		}

		// 将协议代币插入到 tokens 表中
		if len(protocolTokens) > 0 {
			if err := s.tokenRepo.UpsertBatch(protocolTokens); err != nil {
				return fmt.Errorf("failed to upsert protocol tokens: %w", err)
			}
			logger.Debug("Protocol tokens synced",
				zap.Uint("address_id", addressID),
				zap.Int("token_count", len(protocolTokens)),
			)
		}

		logger.Debug("Protocols synced",
			zap.Uint("address_id", addressID),
			zap.Int("protocol_count", len(protocols)),
		)
	}

	// 更新最后同步时间戳
	if err := s.addressRepo.UpdateLastSynced(addressID); err != nil {
		return fmt.Errorf("failed to update last synced: %w", err)
	}

	logger.Debug("Address synced successfully",
		zap.Uint("address_id", addressID),
		zap.Int("token_count", len(tokens)),
		zap.Int("protocol_count", len(protocols)),
	)

	return nil
}

// SyncWallet 同步钱包中的所有地址
func (s *SyncService) SyncWallet(ctx context.Context, walletID uint) error {
	addresses, err := s.addressRepo.GetByWalletID(walletID)
	if err != nil {
		return fmt.Errorf("failed to get wallet addresses: %w", err)
	}

	for _, address := range addresses {
		if err := s.SyncAddress(ctx, address.ID); err != nil {
			logger.Error("Failed to sync address in wallet",
				zap.Uint("wallet_id", walletID),
				zap.Uint("address_id", address.ID),
				zap.Error(err),
			)
		}
	}

	return nil
}

// filterSpamTokens 根据常见模式过滤掉垃圾/欺诈代币和协议凭证代币
func filterSpamTokens(tokens []provider.TokenInfo) []provider.TokenInfo {
	spamKeywords := []string{
		"t.me/", "t.ly/", "fli.so/", "wr.do/", "www.",
		"claim", "swap", "redeem", "visit", "airdrop",
		"reward", "voucher", "distribution",
	}

	// 协议凭证代币前缀（aToken, cToken 等）
	// 这些代币会在协议的 asset_token_list 中展开为实际资产，所以要过滤掉
	protocolTokenPrefixes := []string{
		"aEth",         // Aave Ethereum (aEthUSDC, aEthWETH, etc.)
		"aBas",         // Aave Base
		"aArb",         // Aave Arbitrum
		"aOpt",         // Aave Optimism
		"aPol",         // Aave Polygon
		"aAva",         // Aave Avalanche
		"cToken",       // Compound
		"variableDebt", // Aave variable debt tokens
		"stableDebt",   // Aave stable debt tokens
	}

	filtered := make([]provider.TokenInfo, 0, len(tokens))
	for _, token := range tokens {
		isSpam := false
		isProtocolToken := false

		// 检查符号或名称是否包含垃圾关键词
		symbolLower := strings.ToLower(token.Symbol)
		nameLower := strings.ToLower(token.Name)

		for _, keyword := range spamKeywords {
			if strings.Contains(symbolLower, keyword) || strings.Contains(nameLower, keyword) {
				isSpam = true
				break
			}
		}

		// 检查是否是协议凭证代币
		for _, prefix := range protocolTokenPrefixes {
			if strings.HasPrefix(token.Symbol, prefix) {
				isProtocolToken = true
				break
			}
		}

		// 额外检查：价值为 $0 且具有可疑模式的代币
		if token.Price == 0 && token.USDValue == 0 {
			// 检查是否看起来像垃圾代币（包含特殊字符、表情符号等）
			if strings.Contains(token.Symbol, "✅") || strings.Contains(token.Name, "✅") {
				isSpam = true
			}
			if strings.Contains(token.Symbol, "$") && !strings.HasPrefix(token.Symbol, "$") {
				isSpam = true
			}
		}

		// 只保留非垃圾且非协议凭证的代币
		if !isSpam && !isProtocolToken {
			filtered = append(filtered, token)
		}
	}

	return filtered
}
