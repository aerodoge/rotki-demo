#!/bin/bash

# RPC Nodes API Test Script
# Tests all CRUD operations for RPC nodes

API_BASE="http://localhost:8080/api/v1"

echo "=== RPC Nodes API Test ==="
echo ""

# Test 1: Create RPC Node for Ethereum
echo "1. Creating RPC node for Ethereum..."
CREATE_RESPONSE=$(curl -s -X POST "$API_BASE/rpc-nodes" \
  -H "Content-Type: application/json" \
  -d '{
    "chain_id": "eth",
    "name": "0xRPC",
    "url": "https://eth.llamarpc.com",
    "weight": 100,
    "priority": 1,
    "timeout": 30,
    "is_enabled": true
  }')

NODE_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
echo "Created node with ID: $NODE_ID"
echo ""

# Test 2: List all nodes
echo "2. Listing all RPC nodes..."
curl -s "$API_BASE/rpc-nodes" | jq '.'
echo ""

# Test 3: Get specific node
echo "3. Getting node $NODE_ID..."
curl -s "$API_BASE/rpc-nodes/$NODE_ID" | jq '.'
echo ""

# Test 4: Create another node for BSC
echo "4. Creating RPC node for BSC..."
curl -s -X POST "$API_BASE/rpc-nodes" \
  -H "Content-Type: application/json" \
  -d '{
    "chain_id": "bsc",
    "name": "PublicNode BSC",
    "url": "https://bsc-rpc.publicnode.com",
    "weight": 80,
    "priority": 0,
    "timeout": 30,
    "is_enabled": true
  }' | jq '.'
echo ""

# Test 5: Get grouped nodes
echo "5. Getting nodes grouped by chain..."
curl -s "$API_BASE/rpc-nodes/grouped" | jq '.'
echo ""

# Test 6: Filter nodes by chain
echo "6. Filtering nodes by chain (eth)..."
curl -s "$API_BASE/rpc-nodes?chain_id=eth" | jq '.'
echo ""

# Test 7: Update node
echo "7. Updating node $NODE_ID..."
curl -s -X PUT "$API_BASE/rpc-nodes/$NODE_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "chain_id": "eth",
    "name": "0xRPC Updated",
    "url": "https://eth.llamarpc.com",
    "weight": 90,
    "priority": 2,
    "timeout": 30,
    "is_enabled": true
  }' | jq '.'
echo ""

# Test 8: Check connection
echo "8. Checking connection for node $NODE_ID..."
curl -s -X POST "$API_BASE/rpc-nodes/$NODE_ID/check" | jq '.'
echo ""

# Test 9: Check all connections
echo "9. Checking all connections..."
curl -s -X POST "$API_BASE/rpc-nodes/check-all" | jq '.'
echo ""

# Test 10: Delete node
echo "10. Deleting node $NODE_ID..."
curl -s -X DELETE "$API_BASE/rpc-nodes/$NODE_ID"
echo "Node deleted (HTTP 204 No Content expected)"
echo ""

# Test 11: Verify deletion
echo "11. Verifying deletion - listing all nodes again..."
curl -s "$API_BASE/rpc-nodes" | jq '.'
echo ""

echo "=== Test Complete ==="
