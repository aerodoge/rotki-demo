package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/miles/rotki-demo/internal/logger"
	"github.com/miles/rotki-demo/internal/models"
	"github.com/miles/rotki-demo/internal/repository"
	"go.uber.org/zap"
)

// ChainInitializer handles initialization of chain data
type ChainInitializer struct {
	chainRepo *repository.ChainRepository
}

// NewChainInitializer creates a new chain initializer
func NewChainInitializer(chainRepo *repository.ChainRepository) *ChainInitializer {
	return &ChainInitializer{
		chainRepo: chainRepo,
	}
}

// DeBankChainInfo represents the chain data structure from DeBank API
type DeBankChainInfo struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	LogoURL       string  `json:"logo_url"`
	TokenID       string  `json:"token_id"`
	TokenSymbol   string  `json:"token_symbol"`
	NetworkID     int     `json:"network_id"`
	BlockInterval float64 `json:"block_interval"`
}

// InitializeAllChains loads all chains from chains.json and populates the database
func (ci *ChainInitializer) InitializeAllChains(chainsFilePath string) error {
	logger.Info("Initializing chains from file", zap.String("path", chainsFilePath))

	// Read chains.json file
	data, err := os.ReadFile(chainsFilePath)
	if err != nil {
		return fmt.Errorf("failed to read chains file: %w", err)
	}

	// Parse JSON
	var debankChains []DeBankChainInfo
	if err := json.Unmarshal(data, &debankChains); err != nil {
		return fmt.Errorf("failed to parse chains JSON: %w", err)
	}

	logger.Info("Found chains in file", zap.Int("count", len(debankChains)))

	// Convert to our chain model and batch insert
	chains := make([]models.Chain, 0, len(debankChains))
	for _, dc := range debankChains {
		chain := models.Chain{
			ID:            dc.ID,
			Name:          dc.Name,
			ChainType:     "EVM", // All chains in DeBank are EVM-compatible
			LogoURL:       dc.LogoURL,
			NativeTokenID: dc.TokenID,
			IsActive:      true,
		}
		chains = append(chains, chain)
	}

	// Batch upsert all chains
	if err := ci.chainRepo.UpsertBatch(chains); err != nil {
		return fmt.Errorf("failed to upsert chains: %w", err)
	}

	logger.Info("Successfully initialized chains", zap.Int("count", len(chains)))
	return nil
}

// InitializeAllChainsFromDefault loads chains from the default location
func (ci *ChainInitializer) InitializeAllChainsFromDefault() error {
	// Try multiple possible locations
	possiblePaths := []string{
		"frontend/public/images/chains/chains.json",
		"../frontend/public/images/chains/chains.json",
		"../../frontend/public/images/chains/chains.json",
		"/app/frontend/public/images/chains/chains.json", // Docker path
	}

	for _, path := range possiblePaths {
		absPath, _ := filepath.Abs(path)
		if _, err := os.Stat(absPath); err == nil {
			logger.Info("Found chains file", zap.String("path", absPath))
			return ci.InitializeAllChains(absPath)
		}
	}

	logger.Warn("Could not find chains.json file in any default location, chains will be populated dynamically")
	return nil
}
