package repository

import (
	"context"

	"github.com/rotki-demo/internal/models"
	"gorm.io/gorm"
)

// RPCNodeRepository 处理 RPC 节点数据访问
type RPCNodeRepository struct {
	db *gorm.DB
}

// NewRPCNodeRepository 创建一个新的 RPC 节点仓库
func NewRPCNodeRepository(db *gorm.DB) *RPCNodeRepository {
	return &RPCNodeRepository{db: db}
}

// Create 创建一个新的 RPC 节点
func (r *RPCNodeRepository) Create(ctx context.Context, node *models.RPCNode) error {
	return r.db.WithContext(ctx).Create(node).Error
}

// GetByID 根据 ID 获取 RPC 节点
func (r *RPCNodeRepository) GetByID(ctx context.Context, id uint) (*models.RPCNode, error) {
	var node models.RPCNode
	err := r.db.WithContext(ctx).Preload("Chain").First(&node, id).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// GetByChainID 获取特定链的所有 RPC 节点
func (r *RPCNodeRepository) GetByChainID(ctx context.Context, chainID string) ([]models.RPCNode, error) {
	var nodes []models.RPCNode
	err := r.db.WithContext(ctx).
		Preload("Chain").
		Where("chain_id = ?", chainID).
		Order("priority DESC, weight DESC").
		Find(&nodes).Error
	return nodes, err
}

// GetEnabledByChainID 获取特定链的所有已启用 RPC 节点
func (r *RPCNodeRepository) GetEnabledByChainID(ctx context.Context, chainID string) ([]models.RPCNode, error) {
	var nodes []models.RPCNode
	err := r.db.WithContext(ctx).
		Preload("Chain").
		Where("chain_id = ? AND is_enabled = ?", chainID, true).
		Order("priority DESC, weight DESC").
		Find(&nodes).Error
	return nodes, err
}

// GetAll 获取所有 RPC 节点
func (r *RPCNodeRepository) GetAll(ctx context.Context) ([]models.RPCNode, error) {
	var nodes []models.RPCNode
	err := r.db.WithContext(ctx).
		Preload("Chain").
		Order("chain_id, priority DESC, weight DESC").
		Find(&nodes).Error
	return nodes, err
}

// Update 更新 RPC 节点
func (r *RPCNodeRepository) Update(ctx context.Context, node *models.RPCNode) error {
	return r.db.WithContext(ctx).Save(node).Error
}

// UpdateConnectionStatus 更新 RPC 节点的连接状态
func (r *RPCNodeRepository) UpdateConnectionStatus(ctx context.Context, id uint, isConnected bool) error {
	now := gorm.Expr("NOW()")
	return r.db.WithContext(ctx).Model(&models.RPCNode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_connected": isConnected,
			"last_checked": now,
		}).Error
}

// Delete 删除 RPC 节点
func (r *RPCNodeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.RPCNode{}, id).Error
}

// GetGroupedByChain 获取按链分组的所有 RPC 节点
func (r *RPCNodeRepository) GetGroupedByChain(ctx context.Context) (map[string][]models.RPCNode, error) {
	var nodes []models.RPCNode
	err := r.db.WithContext(ctx).
		Preload("Chain").
		Order("chain_id, priority DESC, weight DESC").
		Find(&nodes).Error
	if err != nil {
		return nil, err
	}

	// 按链分组
	grouped := make(map[string][]models.RPCNode)
	for _, node := range nodes {
		grouped[node.ChainID] = append(grouped[node.ChainID], node)
	}

	return grouped, nil
}
