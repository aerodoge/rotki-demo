package repository

import (
	"github.com/rotki-demo/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TokenRepository 处理代币数据操作
type TokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository 创建一个新的代币仓库
func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// UpsertBatch 批量插入或更新代币
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

// GetByAddressID 获取地址的所有代币
func (r *TokenRepository) GetByAddressID(addressID uint) ([]models.Token, error) {
	var tokens []models.Token
	err := r.db.Where("address_id = ?", addressID).
		Preload("Chain").
		Order("usd_value DESC").
		Find(&tokens).Error
	return tokens, err
}

// DeleteByAddressID 删除地址的所有代币
func (r *TokenRepository) DeleteByAddressID(addressID uint) error {
	return r.db.Where("address_id = ?", addressID).Delete(&models.Token{}).Error
}

// DeleteWalletTokensByAddressID 删除地址的钱包代币（不包括协议代币）
func (r *TokenRepository) DeleteWalletTokensByAddressID(addressID uint) error {
	return r.db.Where("address_id = ? AND (protocol_id IS NULL OR protocol_id = '')", addressID).Delete(&models.Token{}).Error
}

// DeleteProtocolTokensByAddressID 删除地址的协议代币
func (r *TokenRepository) DeleteProtocolTokensByAddressID(addressID uint) error {
	return r.db.Where("address_id = ? AND protocol_id IS NOT NULL AND protocol_id != ''", addressID).Delete(&models.Token{}).Error
}

// GetTotalValueByAddressID 计算地址的总 USD 价值
func (r *TokenRepository) GetTotalValueByAddressID(addressID uint) (float64, error) {
	var total float64
	err := r.db.Model(&models.Token{}).
		Where("address_id = ?", addressID).
		Select("COALESCE(SUM(usd_value), 0)").
		Scan(&total).Error
	return total, err
}
