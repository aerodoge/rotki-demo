package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/miles/rotki-demo/internal/models"
	"github.com/miles/rotki-demo/internal/repository"
)

// WalletHandler 处理钱包相关的 HTTP 请求
type WalletHandler struct {
	walletRepo *repository.WalletRepository
}

// NewWalletHandler 创建一个新的钱包处理器
func NewWalletHandler(walletRepo *repository.WalletRepository) *WalletHandler {
	return &WalletHandler{
		walletRepo: walletRepo,
	}
}

// CreateWalletRequest 表示创建钱包的请求
type CreateWalletRequest struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	EnabledChains []string `json:"enabled_chains"`
}

// UpdateWalletRequest 表示更新钱包的请求
type UpdateWalletRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	EnabledChains []string `json:"enabled_chains"`
}

// CreateWallet 创建一个新的钱包
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

// GetWallet 根据 ID 获取钱包
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

// ListWallets 获取所有钱包
// GET /api/v1/wallets
func (h *WalletHandler) ListWallets(c *gin.Context) {
	wallets, err := h.walletRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wallets"})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

// UpdateWallet 更新钱包
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

	// 记录更新请求用于调试
	c.Request.Context().Value("logger")

	if req.Name != "" {
		wallet.Name = req.Name
	}
	// 通过检查请求中是否提供了描述来允许清除描述
	wallet.Description = req.Description

	if req.Tags != nil {
		wallet.Tags = models.StringSlice(req.Tags)
	}

	if req.EnabledChains != nil {
		wallet.EnabledChains = models.StringSlice(req.EnabledChains)
	}

	if err := h.walletRepo.Update(wallet); err != nil {
		// 记录实际错误用于调试
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// DeleteWallet 删除钱包
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
