package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotki-demo/internal/repository"
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
// @Summary      获取链列表
// @Description  获取所有支持的区块链列表
// @Tags         chains
// @Accept       json
// @Produce      json
// @Success      200  {array}   github_com_rotki-demo_internal_models.Chain
// @Failure      500  {object}  map[string]string
// @Router       /chains [get]
func (h *ChainHandler) ListChains(c *gin.Context) {
	chains, err := h.chainRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chains"})
		return
	}

	c.JSON(http.StatusOK, chains)
}
