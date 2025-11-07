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

// AddressHandler 处理地址相关的 HTTP 请求
type AddressHandler struct {
	addressRepo *repository.AddressRepository
	tokenRepo   *repository.TokenRepository
	syncService *service.SyncService
}

// NewAddressHandler 创建一个新的地址处理器
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

// CreateAddressRequest 表示创建地址的请求
type CreateAddressRequest struct {
	WalletID  uint   `json:"wallet_id" binding:"required"`
	Address   string `json:"address" binding:"required"`
	ChainType string `json:"chain_type"`
	Label     string `json:"label"`
}

// UpdateAddressRequest 表示更新地址的请求
type UpdateAddressRequest struct {
	Label string             `json:"label"`
	Tags  models.StringSlice `json:"tags"`
}

// CreateAddress 创建一个新的地址
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

	// 在后台触发新地址的即时同步
	// 使用后台上下文而不是请求上下文以避免取消
	go func() {
		ctx := context.Background()
		_ = h.syncService.SyncAddress(ctx, address.ID)
	}()

	c.JSON(http.StatusCreated, address)
}

// GetAddress 根据 ID 获取地址
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

	// 获取此地址的代币
	tokens, err := h.tokenRepo.GetByAddressID(address.ID)
	if err == nil {
		address.Tokens = tokens
	}

	c.JSON(http.StatusOK, address)
}

// ListAddresses 获取所有地址
// GET /api/v1/addresses
func (h *AddressHandler) ListAddresses(c *gin.Context) {
	// 检查是否按钱包过滤
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

	// 补充代币信息
	for i := range addresses {
		tokens, err := h.tokenRepo.GetByAddressID(addresses[i].ID)
		if err == nil {
			addresses[i].Tokens = tokens
		}
	}

	c.JSON(http.StatusOK, addresses)
}

// UpdateAddress 更新地址标签
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

	// 获取代币用于响应
	tokens, err := h.tokenRepo.GetByAddressID(address.ID)
	if err == nil {
		address.Tokens = tokens
	}

	c.JSON(http.StatusOK, address)
}

// DeleteAddress 删除地址
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

// RefreshAddress 触发特定地址的同步
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

	// 获取更新后的地址数据
	address, err := h.addressRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve address"})
		return
	}

	// 获取代币
	tokens, err := h.tokenRepo.GetByAddressID(address.ID)
	if err == nil {
		address.Tokens = tokens
	}

	c.JSON(http.StatusOK, address)
}

// RefreshWallet 触发钱包中所有地址的同步
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
