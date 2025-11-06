package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miles/rotki-demo/internal/repository"
)

// ChainHandler handles chain-related HTTP requests
type ChainHandler struct {
	chainRepo *repository.ChainRepository
}

// NewChainHandler creates a new chain handler
func NewChainHandler(chainRepo *repository.ChainRepository) *ChainHandler {
	return &ChainHandler{
		chainRepo: chainRepo,
	}
}

// ListChains retrieves all chains
// GET /api/v1/chains
func (h *ChainHandler) ListChains(c *gin.Context) {
	chains, err := h.chainRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chains"})
		return
	}

	c.JSON(http.StatusOK, chains)
}
