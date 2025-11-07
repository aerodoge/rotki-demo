package repository

import (
	"time"

	"github.com/rotki-demo/internal/models"
	"gorm.io/gorm"
)

// AddressRepository 处理地址数据操作
type AddressRepository struct {
	db *gorm.DB
}

// NewAddressRepository 创建一个新的地址仓库
func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

// Create 创建一个新地址
func (r *AddressRepository) Create(address *models.Address) error {
	return r.db.Create(address).Error
}

// GetByID 根据 ID 获取地址
func (r *AddressRepository) GetByID(id uint) (*models.Address, error) {
	var address models.Address
	err := r.db.Preload("Wallet").Preload("Tokens").First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// GetByAddress 根据地址字符串和链类型获取地址
func (r *AddressRepository) GetByAddress(addr string, chainType string) (*models.Address, error) {
	var address models.Address
	err := r.db.Where("address = ? AND chain_type = ?", addr, chainType).
		Preload("Wallet").
		Preload("Tokens").
		First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// GetByWalletID 获取钱包的所有地址
func (r *AddressRepository) GetByWalletID(walletID uint) ([]models.Address, error) {
	var addresses []models.Address
	err := r.db.Where("wallet_id = ?", walletID).
		Preload("Tokens").
		Find(&addresses).Error
	return addresses, err
}

// List 获取所有地址
func (r *AddressRepository) List() ([]models.Address, error) {
	var addresses []models.Address
	err := r.db.Preload("Wallet").Preload("Tokens").Find(&addresses).Error
	return addresses, err
}

// Update 更新地址
func (r *AddressRepository) Update(address *models.Address) error {
	return r.db.Save(address).Error
}

// UpdateLastSynced 更新最后同步时间戳
func (r *AddressRepository) UpdateLastSynced(id uint) error {
	now := time.Now()
	return r.db.Model(&models.Address{}).Where("id = ?", id).Update("last_synced_at", now).Error
}

// Delete 删除地址
func (r *AddressRepository) Delete(id uint) error {
	return r.db.Delete(&models.Address{}, id).Error
}

// GetAllNeedingSync 返回最近未同步的地址
func (r *AddressRepository) GetAllNeedingSync(olderThan time.Duration) ([]models.Address, error) {
	var addresses []models.Address
	cutoffTime := time.Now().Add(-olderThan)
	err := r.db.Where("last_synced_at IS NULL OR last_synced_at < ?", cutoffTime).
		Preload("Wallet").
		Find(&addresses).Error
	return addresses, err
}
