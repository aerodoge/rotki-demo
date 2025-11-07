package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/miles/rotki-demo/internal/models"
	"github.com/miles/rotki-demo/internal/repository"
	"go.uber.org/zap"
)

// RPCNodeService 处理 RPC 节点业务逻辑
type RPCNodeService struct {
	repo   *repository.RPCNodeRepository
	logger *zap.Logger
}

// NewRPCNodeService 创建一个新的 RPC 节点服务
func NewRPCNodeService(repo *repository.RPCNodeRepository, logger *zap.Logger) *RPCNodeService {
	return &RPCNodeService{
		repo:   repo,
		logger: logger,
	}
}

// TestConnection 测试与 RPC 节点的连接
func (s *RPCNodeService) TestConnection(ctx context.Context, url string, timeout int) (bool, error) {
	// 创建带超时的 HTTP 客户端
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// 准备 JSON-RPC 请求（eth_blockNumber 是一个简单的测试）
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 发起请求
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应是否成功
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	// 检查响应是否包含结果或错误
	if _, hasError := result["error"]; hasError {
		return false, fmt.Errorf("RPC error: %v", result["error"])
	}

	if _, hasResult := result["result"]; !hasResult {
		return false, fmt.Errorf("no result in response")
	}

	return true, nil
}

// Create 创建一个新的 RPC 节点并测试其连接
func (s *RPCNodeService) Create(ctx context.Context, node *models.RPCNode) error {
	// 在创建前测试连接
	isConnected, err := s.TestConnection(ctx, node.URL, node.Timeout)
	if err != nil {
		s.logger.Warn("RPC node connection test failed",
			zap.String("url", node.URL),
			zap.Error(err))
	}

	node.IsConnected = isConnected
	now := time.Now()
	node.LastChecked = &now

	return s.repo.Create(ctx, node)
}

// GetByID 根据 ID 获取 RPC 节点
func (s *RPCNodeService) GetByID(ctx context.Context, id uint) (*models.RPCNode, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByChainID 获取特定链的所有 RPC 节点
func (s *RPCNodeService) GetByChainID(ctx context.Context, chainID string) ([]models.RPCNode, error) {
	return s.repo.GetByChainID(ctx, chainID)
}

// GetAll 获取所有 RPC 节点
func (s *RPCNodeService) GetAll(ctx context.Context) ([]models.RPCNode, error) {
	return s.repo.GetAll(ctx)
}

// GetGroupedByChain 获取按链分组的所有 RPC 节点
func (s *RPCNodeService) GetGroupedByChain(ctx context.Context) (map[string][]models.RPCNode, error) {
	return s.repo.GetGroupedByChain(ctx)
}

// Update 更新 RPC 节点
func (s *RPCNodeService) Update(ctx context.Context, node *models.RPCNode) error {
	return s.repo.Update(ctx, node)
}

// Delete 删除 RPC 节点
func (s *RPCNodeService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// CheckConnection 测试特定节点的连接并更新其状态
func (s *RPCNodeService) CheckConnection(ctx context.Context, id uint) (bool, error) {
	node, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to get node: %w", err)
	}

	isConnected, err := s.TestConnection(ctx, node.URL, node.Timeout)
	if err != nil {
		s.logger.Warn("Connection check failed",
			zap.Uint("node_id", id),
			zap.String("url", node.URL),
			zap.Error(err))
	}

	// 更新连接状态
	if err := s.repo.UpdateConnectionStatus(ctx, id, isConnected); err != nil {
		s.logger.Error("Failed to update connection status",
			zap.Uint("node_id", id),
			zap.Error(err))
		return isConnected, fmt.Errorf("failed to update status: %w", err)
	}

	return isConnected, nil
}

// CheckAllConnections 测试所有 RPC 节点并更新其状态
func (s *RPCNodeService) CheckAllConnections(ctx context.Context) error {
	nodes, err := s.repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}

	for _, node := range nodes {
		if !node.IsEnabled {
			continue
		}

		_, err := s.CheckConnection(ctx, node.ID)
		if err != nil {
			s.logger.Error("Failed to check connection",
				zap.Uint("node_id", node.ID),
				zap.String("chain_id", node.ChainID),
				zap.Error(err))
		}
	}

	return nil
}

// GetEnabledNodesByChain 返回链的已启用和已连接节点
// 数据提供者可以使用此方法根据权重选择 RPC 节点
func (s *RPCNodeService) GetEnabledNodesByChain(ctx context.Context, chainID string) ([]models.RPCNode, error) {
	return s.repo.GetEnabledByChainID(ctx, chainID)
}
