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

// RPCNodeService handles RPC node business logic
type RPCNodeService struct {
	repo   *repository.RPCNodeRepository
	logger *zap.Logger
}

// NewRPCNodeService creates a new RPC node service
func NewRPCNodeService(repo *repository.RPCNodeRepository, logger *zap.Logger) *RPCNodeService {
	return &RPCNodeService{
		repo:   repo,
		logger: logger,
	}
}

// TestConnection tests the connection to an RPC node
func (s *RPCNodeService) TestConnection(ctx context.Context, url string, timeout int) (bool, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Prepare JSON-RPC request (eth_blockNumber is a simple test)
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

	// Make request
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

	// Check if response is successful
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if response contains result or error
	if _, hasError := result["error"]; hasError {
		return false, fmt.Errorf("RPC error: %v", result["error"])
	}

	if _, hasResult := result["result"]; !hasResult {
		return false, fmt.Errorf("no result in response")
	}

	return true, nil
}

// Create creates a new RPC node and tests its connection
func (s *RPCNodeService) Create(ctx context.Context, node *models.RPCNode) error {
	// Test connection before creating
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

// GetByID retrieves an RPC node by ID
func (s *RPCNodeService) GetByID(ctx context.Context, id uint) (*models.RPCNode, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByChainID retrieves all RPC nodes for a specific chain
func (s *RPCNodeService) GetByChainID(ctx context.Context, chainID string) ([]models.RPCNode, error) {
	return s.repo.GetByChainID(ctx, chainID)
}

// GetAll retrieves all RPC nodes
func (s *RPCNodeService) GetAll(ctx context.Context) ([]models.RPCNode, error) {
	return s.repo.GetAll(ctx)
}

// GetGroupedByChain retrieves all RPC nodes grouped by chain
func (s *RPCNodeService) GetGroupedByChain(ctx context.Context) (map[string][]models.RPCNode, error) {
	return s.repo.GetGroupedByChain(ctx)
}

// Update updates an RPC node
func (s *RPCNodeService) Update(ctx context.Context, node *models.RPCNode) error {
	return s.repo.Update(ctx, node)
}

// Delete deletes an RPC node
func (s *RPCNodeService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// CheckConnection tests the connection for a specific node and updates its status
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

	// Update connection status
	if err := s.repo.UpdateConnectionStatus(ctx, id, isConnected); err != nil {
		s.logger.Error("Failed to update connection status",
			zap.Uint("node_id", id),
			zap.Error(err))
		return isConnected, fmt.Errorf("failed to update status: %w", err)
	}

	return isConnected, nil
}

// CheckAllConnections tests all RPC nodes and updates their status
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

// GetEnabledNodesByChain returns enabled and connected nodes for a chain
// This can be used by the data provider to select RPC nodes based on weight
func (s *RPCNodeService) GetEnabledNodesByChain(ctx context.Context, chainID string) ([]models.RPCNode, error) {
	return s.repo.GetEnabledByChainID(ctx, chainID)
}
