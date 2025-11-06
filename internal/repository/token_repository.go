package repository

import (
	"github.com/miles/rotki-demo/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TokenRepository handles token data operations
type TokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// UpsertBatch inserts or updates tokens in batch
func (r *TokenRepository) UpsertBatch(tokens []models.Token) error {
	if len(tokens) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "address_id"},
			{Name: "chain_id"},
			{Name: "token_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"symbol", "name", "decimals", "logo_url",
			"balance", "price", "usd_value", "last_updated",
		}),
	}).Create(&tokens).Error
}

// GetByAddressID retrieves all tokens for an address
func (r *TokenRepository) GetByAddressID(addressID uint) ([]models.Token, error) {
	var tokens []models.Token
	err := r.db.Where("address_id = ?", addressID).
		Preload("Chain").
		Order("usd_value DESC").
		Find(&tokens).Error
	return tokens, err
}

// DeleteByAddressID deletes all tokens for an address
func (r *TokenRepository) DeleteByAddressID(addressID uint) error {
	return r.db.Where("address_id = ?", addressID).Delete(&models.Token{}).Error
}

// GetTotalValueByAddressID calculates total USD value for an address
func (r *TokenRepository) GetTotalValueByAddressID(addressID uint) (float64, error) {
	var total float64
	err := r.db.Model(&models.Token{}).
		Where("address_id = ?", addressID).
		Select("COALESCE(SUM(usd_value), 0)").
		Scan(&total).Error
	return total, err
}
