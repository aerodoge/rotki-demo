package repository

import (
	"github.com/miles/rotki-demo/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ChainRepository handles chain data access
type ChainRepository struct {
	db *gorm.DB
}

// NewChainRepository creates a new chain repository
func NewChainRepository(db *gorm.DB) *ChainRepository {
	return &ChainRepository{db: db}
}

// UpsertBatch inserts or updates multiple chains
func (r *ChainRepository) UpsertBatch(chains []models.Chain) error {
	if len(chains) == 0 {
		return nil
	}

	// Use ON DUPLICATE KEY UPDATE for MySQL
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "logo_url", "native_token_id"}),
	}).Create(&chains).Error
}

// GetByID retrieves a chain by ID
func (r *ChainRepository) GetByID(chainID string) (*models.Chain, error) {
	var chain models.Chain
	err := r.db.Where("id = ?", chainID).First(&chain).Error
	if err != nil {
		return nil, err
	}
	return &chain, nil
}

// List retrieves all chains
func (r *ChainRepository) List() ([]models.Chain, error) {
	var chains []models.Chain
	err := r.db.Find(&chains).Error
	return chains, err
}
