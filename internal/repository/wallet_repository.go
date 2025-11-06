package repository

import (
	"github.com/miles/rotki-demo/internal/models"
	"gorm.io/gorm"
)

// WalletRepository handles wallet data operations
type WalletRepository struct {
	db *gorm.DB
}

// NewWalletRepository creates a new wallet repository
func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// Create creates a new wallet
func (r *WalletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

// GetByID retrieves a wallet by ID
func (r *WalletRepository) GetByID(id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Preload("Addresses").First(&wallet, id).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// GetByName retrieves a wallet by name
func (r *WalletRepository) GetByName(name string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("name = ?", name).Preload("Addresses").First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// List retrieves all wallets
func (r *WalletRepository) List() ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := r.db.Preload("Addresses").Find(&wallets).Error
	return wallets, err
}

// Update updates a wallet
func (r *WalletRepository) Update(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

// Delete deletes a wallet
func (r *WalletRepository) Delete(id uint) error {
	return r.db.Delete(&models.Wallet{}, id).Error
}
