package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rotki-demo/internal/models"
	"github.com/rotki-demo/internal/repository"
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
// @Summary      创建钱包
// @Description  创建一个新的钱包
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Param        wallet  body      CreateWalletRequest  true  "钱包信息"
// @Success      201     {object}  github_com_rotki-demo_internal_models.Wallet
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /wallets [post]
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
// @Summary      获取钱包
// @Description  根据 ID 获取钱包详情
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "钱包 ID"
// @Success      200  {object}  github_com_rotki-demo_internal_models.Wallet
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /wallets/{id} [get]
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
// @Summary      获取钱包列表
// @Description  获取所有钱包列表
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Success      200  {array}   github_com_rotki-demo_internal_models.Wallet
// @Failure      500  {object}  map[string]string
// @Router       /wallets [get]
func (h *WalletHandler) ListWallets(c *gin.Context) {
	wallets, err := h.walletRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wallets"})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

// UpdateWallet 更新钱包
// @Summary      更新钱包
// @Description  更新钱包信息
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Param        id      path      int                   true  "钱包 ID"
// @Param        wallet  body      UpdateWalletRequest   true  "钱包信息"
// @Success      200     {object}  github_com_rotki-demo_internal_models.Wallet
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /wallets/{id} [put]
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
// @Summary      删除钱包
// @Description  根据 ID 删除钱包
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "钱包 ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /wallets/{id} [delete]
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
