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
	chainRepo *repository.ChainRepository,
	syncInterval time.Duration,
	batchSize int,
) *SyncService {
	return &SyncService{
		dataProvider: dataProvider,
		walletRepo:   walletRepo,
		addressRepo:  addressRepo,
		tokenRepo:    tokenRepo,
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

	// 更新插入代币
	if err := s.tokenRepo.UpsertBatch(dbTokens); err != nil {
		return fmt.Errorf("failed to upsert tokens: %w", err)
	}

	// 更新最后同步时间戳
	if err := s.addressRepo.UpdateLastSynced(addressID); err != nil {
		return fmt.Errorf("failed to update last synced: %w", err)
	}

	logger.Debug("Address synced successfully",
		zap.Uint("address_id", addressID),
		zap.Int("token_count", len(tokens)),
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

// filterSpamTokens 根据常见模式过滤掉垃圾/欺诈代币
func filterSpamTokens(tokens []provider.TokenInfo) []provider.TokenInfo {
	spamKeywords := []string{
		"t.me/", "t.ly/", "fli.so/", "wr.do/", "www.",
		"claim", "swap", "redeem", "visit", "airdrop",
		"reward", "voucher", "distribution",
	}

	filtered := make([]provider.TokenInfo, 0, len(tokens))
	for _, token := range tokens {
		isSpam := false

		// 检查符号或名称是否包含垃圾关键词
		symbolLower := strings.ToLower(token.Symbol)
		nameLower := strings.ToLower(token.Name)

		for _, keyword := range spamKeywords {
			if strings.Contains(symbolLower, keyword) || strings.Contains(nameLower, keyword) {
				isSpam = true
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

		if !isSpam {
			filtered = append(filtered, token)
		}
	}

	return filtered
}
