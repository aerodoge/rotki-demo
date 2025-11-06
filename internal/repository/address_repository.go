package repository

import (
	"time"

	"github.com/miles/rotki-demo/internal/models"
	"gorm.io/gorm"
)

// AddressRepository handles address data operations
type AddressRepository struct {
	db *gorm.DB
}

// NewAddressRepository creates a new address repository
func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

// Create creates a new address
func (r *AddressRepository) Create(address *models.Address) error {
	return r.db.Create(address).Error
}

// GetByID retrieves an address by ID
func (r *AddressRepository) GetByID(id uint) (*models.Address, error) {
	var address models.Address
	err := r.db.Preload("Wallet").Preload("Tokens").First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// GetByAddress retrieves an address by address string and chain type
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

// GetByWalletID retrieves all addresses for a wallet
func (r *AddressRepository) GetByWalletID(walletID uint) ([]models.Address, error) {
	var addresses []models.Address
	err := r.db.Where("wallet_id = ?", walletID).
		Preload("Tokens").
		Find(&addresses).Error
	return addresses, err
}

// List retrieves all addresses
func (r *AddressRepository) List() ([]models.Address, error) {
	var addresses []models.Address
	err := r.db.Preload("Wallet").Preload("Tokens").Find(&addresses).Error
	return addresses, err
}

// Update updates an address
func (r *AddressRepository) Update(address *models.Address) error {
	return r.db.Save(address).Error
}

// UpdateLastSynced updates the last synced timestamp
func (r *AddressRepository) UpdateLastSynced(id uint) error {
	now := time.Now()
	return r.db.Model(&models.Address{}).Where("id = ?", id).Update("last_synced_at", now).Error
}

// Delete deletes an address
func (r *AddressRepository) Delete(id uint) error {
	return r.db.Delete(&models.Address{}, id).Error
}

// GetAllNeedingSync returns addresses that haven't been synced recently
func (r *AddressRepository) GetAllNeedingSync(olderThan time.Duration) ([]models.Address, error) {
	var addresses []models.Address
	cutoffTime := time.Now().Add(-olderThan)
	err := r.db.Where("last_synced_at IS NULL OR last_synced_at < ?", cutoffTime).
		Preload("Wallet").
		Find(&addresses).Error
	return addresses, err
}
