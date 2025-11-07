package repository

import (
	"github.com/miles/rotki-demo/internal/models"
	"gorm.io/gorm"
)

// WalletRepository 处理钱包数据操作
type WalletRepository struct {
	db *gorm.DB
}

// NewWalletRepository 创建一个新的钱包仓库
func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// Create 创建一个新钱包
func (r *WalletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

// GetByID 根据 ID 获取钱包
func (r *WalletRepository) GetByID(id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Preload("Addresses").First(&wallet, id).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// GetByName 根据名称获取钱包
func (r *WalletRepository) GetByName(name string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("name = ?", name).Preload("Addresses").First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// List 获取所有钱包
func (r *WalletRepository) List() ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := r.db.Preload("Addresses").Find(&wallets).Error
	return wallets, err
}

// Update 更新钱包
func (r *WalletRepository) Update(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

// Delete 删除钱包
func (r *WalletRepository) Delete(id uint) error {
	return r.db.Delete(&models.Wallet{}, id).Error
}
