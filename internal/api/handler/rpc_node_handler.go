package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rotki-demo/internal/models"
	"github.com/rotki-demo/internal/service"
	"go.uber.org/zap"
)

// RPCNodeHandler 处理 RPC 节点相关的 HTTP 请求
type RPCNodeHandler struct {
	service *service.RPCNodeService
	logger  *zap.Logger
}

// NewRPCNodeHandler 创建一个新的 RPC 节点处理器
func NewRPCNodeHandler(service *service.RPCNodeService, logger *zap.Logger) *RPCNodeHandler {
	return &RPCNodeHandler{
		service: service,
		logger:  logger,
	}
}

// CreateRPCNode 处理创建新的 RPC 节点
// @Summary      创建 RPC 节点
// @Description  创建一个新的 RPC 节点
// @Tags         rpc-nodes
// @Accept       json
// @Produce      json
// @Param        node  body      github_com_rotki-demo_internal_models.RPCNode  true  "RPC 节点数据"
// @Success      201   {object}  github_com_rotki-demo_internal_models.RPCNode
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /rpc-nodes [post]
func (h *RPCNodeHandler) CreateRPCNode(c *gin.Context) {
	var req models.RPCNode
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果未提供则设置默认值
	if req.Timeout == 0 {
		req.Timeout = 30
	}
	if req.Weight == 0 {
		req.Weight = 100
	}

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		h.logger.Error("Failed to create RPC node", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create RPC node"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// GetRPCNode 处理获取单个 RPC 节点
// @Summary      获取 RPC 节点
// @Description  根据 ID 获取 RPC 节点
// @Tags         rpc-nodes
// @Produce      json
// @Param        id   path      int  true  "RPC 节点 ID"
// @Success      200  {object}  github_com_rotki-demo_internal_models.RPCNode
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /rpc-nodes/{id} [get]
func (h *RPCNodeHandler) GetRPCNode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	node, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get RPC node", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "RPC node not found"})
		return
	}

	c.JSON(http.StatusOK, node)
}

// ListRPCNodes 处理获取所有 RPC 节点
// @Summary      获取 RPC 节点列表
// @Description  获取所有 RPC 节点，可按链 ID 过滤
// @Tags         rpc-nodes
// @Produce      json
// @Param        chain_id  query     string  false  "按链 ID 过滤"
// @Success      200       {array}   github_com_rotki-demo_internal_models.RPCNode
// @Failure      500       {object}  map[string]string
// @Router       /rpc-nodes [get]
func (h *RPCNodeHandler) ListRPCNodes(c *gin.Context) {
	chainID := c.Query("chain_id")

	var nodes []models.RPCNode
	var err error

	if chainID != "" {
		nodes, err = h.service.GetByChainID(c.Request.Context(), chainID)
	} else {
		nodes, err = h.service.GetAll(c.Request.Context())
	}

	if err != nil {
		h.logger.Error("Failed to list RPC nodes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve RPC nodes"})
		return
	}

	c.JSON(http.StatusOK, nodes)
}

// GetRPCNodesByChain 处理获取按链分组的 RPC 节点
// @Summary Get RPC nodes grouped by chain
// @Tags rpc-nodes
// @Produce json
// @Success 200 {object} map[string][]models.RPCNode
// @Router /api/v1/rpc-nodes/grouped [get]
func (h *RPCNodeHandler) GetRPCNodesByChain(c *gin.Context) {
	grouped, err := h.service.GetGroupedByChain(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get grouped RPC nodes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve RPC nodes"})
		return
	}

	c.JSON(http.StatusOK, grouped)
}

// UpdateRPCNode 处理更新 RPC 节点
// @Summary Update RPC node
// @Tags rpc-nodes
// @Accept json
// @Produce json
// @Param id path int true "RPC node ID"
// @Param node body models.RPCNode true "Updated RPC node data"
// @Success 200 {object} models.RPCNode
// @Router /api/v1/rpc-nodes/{id} [put]
func (h *RPCNodeHandler) UpdateRPCNode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// 首先获取现有节点
	existing, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get RPC node", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "RPC node not found"})
		return
	}

	var req models.RPCNode
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保留 ID 和时间戳
	req.ID = uint(id)
	req.CreatedAt = existing.CreatedAt

	if err := h.service.Update(c.Request.Context(), &req); err != nil {
		h.logger.Error("Failed to update RPC node", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update RPC node"})
		return
	}

	c.JSON(http.StatusOK, req)
}

// DeleteRPCNode 处理删除 RPC 节点
// @Summary Delete RPC node
// @Tags rpc-nodes
// @Param id path int true "RPC node ID"
// @Success 204
// @Router /api/v1/rpc-nodes/{id} [delete]
func (h *RPCNodeHandler) DeleteRPCNode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete RPC node", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete RPC node"})
		return
	}

	c.Status(http.StatusNoContent)
}

// CheckRPCNodeConnection 处理测试特定 RPC 节点的连接
// @Summary Check RPC node connection
// @Tags rpc-nodes
// @Produce json
// @Param id path int true "RPC node ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rpc-nodes/{id}/check [post]
func (h *RPCNodeHandler) CheckRPCNodeConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	isConnected, err := h.service.CheckConnection(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Warn("Connection check failed",
			zap.Uint64("node_id", id),
			zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"connected": false,
			"error":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"connected": isConnected,
	})
}

// CheckAllRPCNodeConnections 处理测试所有 RPC 节点的连接
// @Summary Check all RPC node connections
// @Tags rpc-nodes
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rpc-nodes/check-all [post]
func (h *RPCNodeHandler) CheckAllRPCNodeConnections(c *gin.Context) {
	if err := h.service.CheckAllConnections(c.Request.Context()); err != nil {
		h.logger.Error("Failed to check all connections", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check connections"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Connection checks completed",
	})
}
