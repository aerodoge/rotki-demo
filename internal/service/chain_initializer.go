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

// ChainInitializer 处理链数据的初始化
type ChainInitializer struct {
	chainRepo *repository.ChainRepository
}

// NewChainInitializer 创建一个新的链初始化器
func NewChainInitializer(chainRepo *repository.ChainRepository) *ChainInitializer {
	return &ChainInitializer{
		chainRepo: chainRepo,
	}
}

// DeBankChainInfo 表示从 DeBank API 获取的链数据结构
type DeBankChainInfo struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	LogoURL       string  `json:"logo_url"`
	TokenID       string  `json:"token_id"`
	TokenSymbol   string  `json:"token_symbol"`
	NetworkID     int     `json:"network_id"`
	BlockInterval float64 `json:"block_interval"`
}

// InitializeAllChains 从 chains.json 文件加载所有链并填充数据库
func (ci *ChainInitializer) InitializeAllChains(chainsFilePath string) error {
	logger.Info("Initializing chains from file", zap.String("path", chainsFilePath))

	// 读取 chains.json 文件
	data, err := os.ReadFile(chainsFilePath)
	if err != nil {
		return fmt.Errorf("failed to read chains file: %w", err)
	}

	// 解析 JSON
	var debankChains []DeBankChainInfo
	if err := json.Unmarshal(data, &debankChains); err != nil {
		return fmt.Errorf("failed to parse chains JSON: %w", err)
	}

	logger.Info("Found chains in file", zap.Int("count", len(debankChains)))

	// 转换为我们的链模型并批量插入
	chains := make([]models.Chain, 0, len(debankChains))
	for _, dc := range debankChains {
		chain := models.Chain{
			ID:            dc.ID,
			Name:          dc.Name,
			ChainType:     "EVM", // DeBank 中的所有链都是 EVM 兼容的
			LogoURL:       dc.LogoURL,
			NativeTokenID: dc.TokenID,
			IsActive:      true,
		}
		chains = append(chains, chain)
	}

	// 批量更新或插入所有链
	if err := ci.chainRepo.UpsertBatch(chains); err != nil {
		return fmt.Errorf("failed to upsert chains: %w", err)
	}

	logger.Info("Successfully initialized chains", zap.Int("count", len(chains)))
	return nil
}

// InitializeAllChainsFromDefault 从默认位置加载链
func (ci *ChainInitializer) InitializeAllChainsFromDefault() error {
	// 尝试多个可能的位置
	possiblePaths := []string{
		"frontend/public/images/chains/chains.json",
		"../frontend/public/images/chains/chains.json",
		"../../frontend/public/images/chains/chains.json",
		"/app/frontend/public/images/chains/chains.json", // Docker 路径
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
