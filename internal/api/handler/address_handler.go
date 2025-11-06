package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/miles/rotki-demo/internal/models"
	"github.com/miles/rotki-demo/internal/repository"
	"github.com/miles/rotki-demo/internal/service"
)

// AddressHandler handles address-related HTTP requests
type AddressHandler struct {
	addressRepo *repository.AddressRepository
	tokenRepo   *repository.TokenRepository
	syncService *service.SyncService
}

// NewAddressHandler creates a new address handler
func NewAddressHandler(
	addressRepo *repository.AddressRepository,
	tokenRepo *repository.TokenRepository,
	syncService *service.SyncService,
) *AddressHandler {
	return &AddressHandler{
		addressRepo: addressRepo,
		tokenRepo:   tokenRepo,
		syncService: syncService,
	}
}

// CreateAddressRequest represents the request to create an address
type CreateAddressRequest struct {
	WalletID  uint   `json:"wallet_id" binding:"required"`
	Address   string `json:"address" binding:"required"`
	ChainType string `json:"chain_type"`
	Label     string `json:"label"`
}

// UpdateAddressRequest represents the request to update an address
type UpdateAddressRequest struct {
	Label string            `json:"label"`
	Tags  models.StringSlice `json:"tags"`
}

// CreateAddress creates a new address
// POST /api/v1/addresses
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ChainType == "" {
		req.ChainType = "EVM"
	}

	address := &models.Address{
		WalletID:  req.WalletID,
		Address:   req.Address,
		ChainType: req.ChainType,
		Label:     req.Label,
	}

	if err := h.addressRepo.Create(address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	// Trigger immediate sync for the new address in background
	// Use background context instead of request context to avoid cancellation
	go func() {
		ctx := context.Background()
		_ = h.syncService.SyncAddress(ctx, address.ID)
	}()

	c.JSON(http.StatusCreated, address)
}

// GetAddress retrieves an address by ID
// GET /api/v1/addresses/:id
func (h *AddressHandler) GetAddress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	address, err := h.addressRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	// Get tokens for this address
	tokens, err := h.tokenRepo.GetByAddressID(address.ID)
	if err == nil {
		address.Tokens = tokens
	}

	c.JSON(http.StatusOK, address)
}

// ListAddresses retrieves all addresses
// GET /api/v1/addresses
func (h *AddressHandler) ListAddresses(c *gin.Context) {
	// Check if filtering by wallet
	walletIDStr := c.Query("wallet_id")

	var addresses []models.Address
	var err error

	if walletIDStr != "" {
		walletID, parseErr := strconv.ParseUint(walletIDStr, 10, 32)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
			return
		}
		addresses, err = h.addressRepo.GetByWalletID(uint(walletID))
	} else {
		addresses, err = h.addressRepo.List()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve addresses"})
		return
	}

	// Enrich with tokens
	for i := range addresses {
		tokens, err := h.tokenRepo.GetByAddressID(addresses[i].ID)
		if err == nil {
			addresses[i].Tokens = tokens
		}
	}

	c.JSON(http.StatusOK, addresses)
}

// UpdateAddress updates an address label
// PUT /api/v1/addresses/:id
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	address, err := h.addressRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	var req UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address.Label = req.Label
	address.Tags = req.Tags

	if err := h.addressRepo.Update(address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		return
	}

	// Get tokens for response
	tokens, err := h.tokenRepo.GetByAddressID(address.ID)
	if err == nil {
		address.Tokens = tokens
	}

	c.JSON(http.StatusOK, address)
}

// DeleteAddress deletes an address
// DELETE /api/v1/addresses/:id
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	if err := h.addressRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}

// RefreshAddress triggers a sync for a specific address
// POST /api/v1/addresses/:id/refresh
func (h *AddressHandler) RefreshAddress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	if err := h.syncService.SyncAddress(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh address"})
		return
	}

	// Get updated address data
	address, err := h.addressRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve address"})
		return
	}

	// Get tokens
	tokens, err := h.tokenRepo.GetByAddressID(address.ID)
	if err == nil {
		address.Tokens = tokens
	}

	c.JSON(http.StatusOK, address)
}

// RefreshWallet triggers a sync for all addresses in a wallet
// POST /api/v1/wallets/:id/refresh
func (h *AddressHandler) RefreshWallet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	if err := h.syncService.SyncWallet(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet refreshed successfully"})
}
