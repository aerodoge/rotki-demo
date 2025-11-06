package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/miles/rotki-demo/internal/models"
	"github.com/miles/rotki-demo/internal/repository"
)

// WalletHandler handles wallet-related HTTP requests
type WalletHandler struct {
	walletRepo *repository.WalletRepository
}

// NewWalletHandler creates a new wallet handler
func NewWalletHandler(walletRepo *repository.WalletRepository) *WalletHandler {
	return &WalletHandler{
		walletRepo: walletRepo,
	}
}

// CreateWalletRequest represents the request to create a wallet
type CreateWalletRequest struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	EnabledChains []string `json:"enabled_chains"`
}

// UpdateWalletRequest represents the request to update a wallet
type UpdateWalletRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	EnabledChains []string `json:"enabled_chains"`
}

// CreateWallet creates a new wallet
// POST /api/v1/wallets
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet := &models.Wallet{
		Name:          req.Name,
		Description:   req.Description,
		Tags:          models.StringSlice(req.Tags),
		EnabledChains: models.StringSlice(req.EnabledChains),
	}

	if err := h.walletRepo.Create(wallet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

// GetWallet retrieves a wallet by ID
// GET /api/v1/wallets/:id
func (h *WalletHandler) GetWallet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	wallet, err := h.walletRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// ListWallets retrieves all wallets
// GET /api/v1/wallets
func (h *WalletHandler) ListWallets(c *gin.Context) {
	wallets, err := h.walletRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wallets"})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

// UpdateWallet updates a wallet
// PUT /api/v1/wallets/:id
func (h *WalletHandler) UpdateWallet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	wallet, err := h.walletRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	var req UpdateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the update request for debugging
	c.Request.Context().Value("logger")

	if req.Name != "" {
		wallet.Name = req.Name
	}
	// Allow clearing description by checking if it was provided in request
	wallet.Description = req.Description

	if req.Tags != nil {
		wallet.Tags = models.StringSlice(req.Tags)
	}

	if req.EnabledChains != nil {
		wallet.EnabledChains = models.StringSlice(req.EnabledChains)
	}

	if err := h.walletRepo.Update(wallet); err != nil {
		// Log the actual error for debugging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// DeleteWallet deletes a wallet
// DELETE /api/v1/wallets/:id
func (h *WalletHandler) DeleteWallet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	if err := h.walletRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet deleted successfully"})
}
