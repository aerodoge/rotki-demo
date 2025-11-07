package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miles/rotki-demo/internal/repository"
)

// ChainHandler 处理链相关的 HTTP 请求
type ChainHandler struct {
	chainRepo *repository.ChainRepository
}

// NewChainHandler 创建一个新的链处理器
func NewChainHandler(chainRepo *repository.ChainRepository) *ChainHandler {
	return &ChainHandler{
		chainRepo: chainRepo,
	}
}

// ListChains 获取所有链
// GET /api/v1/chains
func (h *ChainHandler) ListChains(c *gin.Context) {
	chains, err := h.chainRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chains"})
		return
	}

	c.JSON(http.StatusOK, chains)
}
