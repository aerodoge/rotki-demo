package repository

import (
	"context"

	"github.com/miles/rotki-demo/internal/models"
	"gorm.io/gorm"
)

// RPCNodeRepository handles RPC node data access
type RPCNodeRepository struct {
	db *gorm.DB
}

// NewRPCNodeRepository creates a new RPC node repository
func NewRPCNodeRepository(db *gorm.DB) *RPCNodeRepository {
	return &RPCNodeRepository{db: db}
}

// Create creates a new RPC node
func (r *RPCNodeRepository) Create(ctx context.Context, node *models.RPCNode) error {
	return r.db.WithContext(ctx).Create(node).Error
}

// GetByID retrieves an RPC node by ID
func (r *RPCNodeRepository) GetByID(ctx context.Context, id uint) (*models.RPCNode, error) {
	var node models.RPCNode
	err := r.db.WithContext(ctx).Preload("Chain").First(&node, id).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// GetByChainID retrieves all RPC nodes for a specific chain
func (r *RPCNodeRepository) GetByChainID(ctx context.Context, chainID string) ([]models.RPCNode, error) {
	var nodes []models.RPCNode
	err := r.db.WithContext(ctx).
		Preload("Chain").
		Where("chain_id = ?", chainID).
		Order("priority DESC, weight DESC").
		Find(&nodes).Error
	return nodes, err
}

// GetEnabledByChainID retrieves all enabled RPC nodes for a specific chain
func (r *RPCNodeRepository) GetEnabledByChainID(ctx context.Context, chainID string) ([]models.RPCNode, error) {
	var nodes []models.RPCNode
	err := r.db.WithContext(ctx).
		Preload("Chain").
		Where("chain_id = ? AND is_enabled = ?", chainID, true).
		Order("priority DESC, weight DESC").
		Find(&nodes).Error
	return nodes, err
}

// GetAll retrieves all RPC nodes
func (r *RPCNodeRepository) GetAll(ctx context.Context) ([]models.RPCNode, error) {
	var nodes []models.RPCNode
	err := r.db.WithContext(ctx).
		Preload("Chain").
		Order("chain_id, priority DESC, weight DESC").
		Find(&nodes).Error
	return nodes, err
}

// Update updates an RPC node
func (r *RPCNodeRepository) Update(ctx context.Context, node *models.RPCNode) error {
	return r.db.WithContext(ctx).Save(node).Error
}

// UpdateConnectionStatus updates the connection status of an RPC node
func (r *RPCNodeRepository) UpdateConnectionStatus(ctx context.Context, id uint, isConnected bool) error {
	now := gorm.Expr("NOW()")
	return r.db.WithContext(ctx).Model(&models.RPCNode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_connected": isConnected,
			"last_checked": now,
		}).Error
}

// Delete deletes an RPC node
func (r *RPCNodeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.RPCNode{}, id).Error
}

// GetGroupedByChain retrieves all RPC nodes grouped by chain
func (r *RPCNodeRepository) GetGroupedByChain(ctx context.Context) (map[string][]models.RPCNode, error) {
	var nodes []models.RPCNode
	err := r.db.WithContext(ctx).
		Preload("Chain").
		Order("chain_id, priority DESC, weight DESC").
		Find(&nodes).Error
	if err != nil {
		return nil, err
	}

	// Group by chain
	grouped := make(map[string][]models.RPCNode)
	for _, node := range nodes {
		grouped[node.ChainID] = append(grouped[node.ChainID], node)
	}

	return grouped, nil
}
