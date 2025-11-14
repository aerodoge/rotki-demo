package repository

import (
	"github.com/rotki-demo/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ChainRepository 处理链数据访问
type ChainRepository struct {
	db *gorm.DB
}

// NewChainRepository 创建一个新的链仓库
func NewChainRepository(db *gorm.DB) *ChainRepository {
	return &ChainRepository{db: db}
}

// UpsertBatch 插入或更新多个链
func (r *ChainRepository) UpsertBatch(chains []models.Chain) error {
	if len(chains) == 0 {
		return nil
	}

	// 为 MySQL 使用 ON DUPLICATE KEY UPDATE
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "logo_url", "native_token_id"}),
	}).Create(&chains).Error
}

// GetByID 根据 ID 获取链
func (r *ChainRepository) GetByID(chainID string) (*models.Chain, error) {
	var chain models.Chain
	err := r.db.Where("id = ?", chainID).First(&chain).Error
	if err != nil {
		return nil, err
	}
	return &chain, nil
}

// List 获取所有链
func (r *ChainRepository) List() ([]models.Chain, error) {
	// 只返回我们支持的链
	supportedChainIDs := []string{
		"eth",    // Ethereum
		"arb",    // Arbitrum
		"op",     // Optimism
		"base",   // Base
		"uni",    // Unichain
		"plasma", // Plasma
		"scrl",   // Scroll
		"plume",  // Plume
		"matic",  // Polygon
		"ink",    // Ink
		"hyper",  // HyperEVM
		"bsc",    // BSC
		"bera",   // Berachain
	}

	var chains []models.Chain
	err := r.db.Where("id IN ?", supportedChainIDs).Find(&chains).Error
	return chains, err
}
