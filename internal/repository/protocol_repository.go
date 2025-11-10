package repository

import (
	"github.com/rotki-demo/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProtocolRepository 处理协议持仓的数据库操作
type ProtocolRepository struct {
	db *gorm.DB
}

// NewProtocolRepository 创建一个新的协议仓库
func NewProtocolRepository(db *gorm.DB) *ProtocolRepository {
	return &ProtocolRepository{db: db}
}

// GetByAddressID 根据地址 ID 获取所有协议
func (r *ProtocolRepository) GetByAddressID(addressID uint) ([]models.Protocol, error) {
	var protocols []models.Protocol
	if err := r.db.
		Preload("Chain").
		Where("address_id = ?", addressID).
		Order("net_usd_value DESC").
		Find(&protocols).Error; err != nil {
		return nil, err
	}
	return protocols, nil
}

// UpsertBatch 批量更新或插入协议
func (r *ProtocolRepository) UpsertBatch(protocols []models.Protocol) error {
	if len(protocols) == 0 {
		return nil
	}

	// 使用 ON CONFLICT 更新或插入
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "address_id"},
			{Name: "protocol_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"site_url",
			"logo_url",
			"chain_id",
			"net_usd_value",
			"asset_usd_value",
			"debt_usd_value",
			"position_type",
			"raw_data",
			"last_updated",
		}),
	}).Create(&protocols).Error
}

// DeleteByAddressID 删除地址的所有协议（在同步前清理）
func (r *ProtocolRepository) DeleteByAddressID(addressID uint) error {
	return r.db.Where("address_id = ?", addressID).Delete(&models.Protocol{}).Error
}

// GetTotalValueByAddress 获取地址的协议总价值
func (r *ProtocolRepository) GetTotalValueByAddress(addressID uint) (float64, error) {
	var total float64
	err := r.db.Model(&models.Protocol{}).
		Where("address_id = ?", addressID).
		Select("COALESCE(SUM(net_usd_value), 0)").
		Scan(&total).Error
	return total, err
}
