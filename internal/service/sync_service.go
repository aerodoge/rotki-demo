package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/miles/rotki-demo/internal/logger"
	"github.com/miles/rotki-demo/internal/models"
	"github.com/miles/rotki-demo/internal/provider"
	"github.com/miles/rotki-demo/internal/repository"
	"go.uber.org/zap"
)

// SyncService handles data synchronization
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

// NewSyncService creates a new sync service
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

// Start starts the background sync process
func (s *SyncService) Start() {
	s.wg.Add(1)
	go s.syncLoop()
	logger.Info("Sync service started", zap.Duration("interval", s.syncInterval))
}

// Stop stops the background sync process
func (s *SyncService) Stop() {
	close(s.stopChan)
	s.wg.Wait()
	logger.Info("Sync service stopped")
}

// syncLoop runs periodic sync
func (s *SyncService) syncLoop() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	// Run initial sync
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

// syncAll syncs all addresses that need updating
func (s *SyncService) syncAll() {
	ctx := context.Background()

	// Get addresses that need syncing
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

	// Process in batches
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

// syncBatch syncs a batch of addresses concurrently
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

// SyncAddress syncs data for a specific address
func (s *SyncService) SyncAddress(ctx context.Context, addressID uint) error {
	// Get address details
	address, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return fmt.Errorf("failed to get address: %w", err)
	}

	// Get wallet to check enabled chains
	wallet, err := s.walletRepo.GetByID(address.WalletID)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	logger.Debug("Syncing address",
		zap.Uint("address_id", addressID),
		zap.String("address", address.Address),
		zap.Any("enabled_chains", wallet.EnabledChains),
	)

	// Determine which chains to query
	var chainIDsToQuery []string
	if len(wallet.EnabledChains) > 0 {
		chainIDsToQuery = wallet.EnabledChains
	}

	// First, get and upsert chains
	chains, err := s.dataProvider.GetUsedChainList(ctx, address.Address)
	if err != nil {
		return fmt.Errorf("failed to get chain list: %w", err)
	}

	// Filter chains based on enabled chains
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

	// Convert to database models
	dbChains := make([]models.Chain, 0, len(chains))
	for _, chain := range chains {
		dbChains = append(dbChains, models.Chain{
			ID:            chain.ChainID,
			Name:          chain.Name,
			LogoURL:       chain.LogoURL,
			NativeTokenID: chain.NativeTokenID,
		})
	}

	// Upsert chains first
	if err := s.chainRepo.UpsertBatch(dbChains); err != nil {
		return fmt.Errorf("failed to upsert chains: %w", err)
	}

	// Get token list from provider, optionally filtered by chains
	tokens, err := s.dataProvider.GetTokenList(ctx, address.Address, chainIDsToQuery)
	if err != nil {
		return fmt.Errorf("failed to get token list: %w", err)
	}

	// Filter out spam/scam tokens
	filteredTokens := filterSpamTokens(tokens)

	// Convert to database models
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

	// Upsert tokens
	if err := s.tokenRepo.UpsertBatch(dbTokens); err != nil {
		return fmt.Errorf("failed to upsert tokens: %w", err)
	}

	// Update last synced timestamp
	if err := s.addressRepo.UpdateLastSynced(addressID); err != nil {
		return fmt.Errorf("failed to update last synced: %w", err)
	}

	logger.Debug("Address synced successfully",
		zap.Uint("address_id", addressID),
		zap.Int("token_count", len(tokens)),
	)

	return nil
}

// SyncWallet syncs all addresses in a wallet
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

// filterSpamTokens filters out spam/scam tokens based on common patterns
func filterSpamTokens(tokens []provider.TokenInfo) []provider.TokenInfo {
	spamKeywords := []string{
		"t.me/", "t.ly/", "fli.so/", "wr.do/", "www.",
		"claim", "swap", "redeem", "visit", "airdrop",
		"reward", "voucher", "distribution",
	}

	filtered := make([]provider.TokenInfo, 0, len(tokens))
	for _, token := range tokens {
		isSpam := false

		// Check if symbol or name contains spam keywords
		symbolLower := strings.ToLower(token.Symbol)
		nameLower := strings.ToLower(token.Name)

		for _, keyword := range spamKeywords {
			if strings.Contains(symbolLower, keyword) || strings.Contains(nameLower, keyword) {
				isSpam = true
				break
			}
		}

		// Additional check: tokens with $0 value and suspicious patterns
		if token.Price == 0 && token.USDValue == 0 {
			// Check if it looks like a spam token (contains special characters, emojis, etc.)
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
