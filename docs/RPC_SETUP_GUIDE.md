# RPC Nodes Setup Guide

## Quick Start

### 1. Run Database Migration

Apply the RPC nodes table migration:

```bash
# Option 1: Using MySQL command line
mysql -u root -p rotki_demo < migrations/002_add_rpc_nodes_table.sql

# Option 2: Using GORM auto-migration (uncomment in main.go)
# Edit cmd/server/main.go and uncomment the AutoMigrate section
```

### 2. Start the Backend

```bash
# Build and run
go run cmd/server/main.go

# Or use make
make run
```

The API will be available at `http://localhost:8080`

### 3. Start the Frontend

```bash
cd frontend
npm install  # if not already done
npm run dev
```

The UI will be available at `http://localhost:3000`

### 4. Access RPC Nodes Settings

1. Open http://localhost:3000
2. Click on the sidebar menu
3. Expand **Settings**
4. Click **RPC Nodes**

## Configuration Examples

### Ethereum Mainnet

#### Free Public Nodes

```json
{
  "chain_id": "eth",
  "name": "LlamaNodes",
  "url": "https://eth.llamarpc.com",
  "weight": 40,
  "priority": 1,
  "is_enabled": true
}
```

```json
{
  "chain_id": "eth",
  "name": "PublicNode",
  "url": "https://ethereum-rpc.publicnode.com",
  "weight": 30,
  "priority": 0,
  "is_enabled": true
}
```

```json
{
  "chain_id": "eth",
  "name": "1RPC",
  "url": "https://1rpc.io/eth",
  "weight": 30,
  "priority": 0,
  "is_enabled": true
}
```

#### Paid Nodes (Higher Priority)

```json
{
  "chain_id": "eth",
  "name": "Alchemy",
  "url": "https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY",
  "weight": 100,
  "priority": 10,
  "is_enabled": true
}
```

```json
{
  "chain_id": "eth",
  "name": "Infura",
  "url": "https://mainnet.infura.io/v3/YOUR_PROJECT_ID",
  "weight": 100,
  "priority": 10,
  "is_enabled": true
}
```

### BSC (Binance Smart Chain)

```json
{
  "chain_id": "bsc",
  "name": "BSC RPC",
  "url": "https://bsc-dataseed1.binance.org",
  "weight": 50,
  "priority": 1,
  "is_enabled": true
}
```

```json
{
  "chain_id": "bsc",
  "name": "PublicNode BSC",
  "url": "https://bsc-rpc.publicnode.com",
  "weight": 50,
  "priority": 0,
  "is_enabled": true
}
```

### Polygon

```json
{
  "chain_id": "matic",
  "name": "Polygon RPC",
  "url": "https://polygon-rpc.com",
  "weight": 100,
  "priority": 1,
  "is_enabled": true
}
```

### Arbitrum

```json
{
  "chain_id": "arb",
  "name": "Arbitrum One",
  "url": "https://arb1.arbitrum.io/rpc",
  "weight": 100,
  "priority": 1,
  "is_enabled": true
}
```

### Optimism

```json
{
  "chain_id": "op",
  "name": "Optimism",
  "url": "https://mainnet.optimism.io",
  "weight": 100,
  "priority": 1,
  "is_enabled": true
}
```

## Testing the API

### Using the Test Script

```bash
# Make sure backend is running
./scripts/test_rpc_nodes.sh
```

### Manual Testing with curl

Create a node:
```bash
curl -X POST http://localhost:8080/api/v1/rpc-nodes \
  -H "Content-Type: application/json" \
  -d '{
    "chain_id": "eth",
    "name": "Test Node",
    "url": "https://eth.llamarpc.com",
    "weight": 100,
    "priority": 1,
    "timeout": 30,
    "is_enabled": true
  }'
```

List all nodes:
```bash
curl http://localhost:8080/api/v1/rpc-nodes | jq '.'
```

Check connection:
```bash
curl -X POST http://localhost:8080/api/v1/rpc-nodes/1/check | jq '.'
```

## Load Balancing Strategy

### Priority-First

Nodes are selected based on:
1. **Priority** (higher first)
2. **Weight** (percentage within same priority)
3. **Connection status** (only connected nodes)

### Example Configuration

```
Priority 10 (Premium Nodes):
- Alchemy: 100% weight
- Infura: 100% weight
→ 50% of requests go to Alchemy, 50% to Infura

Priority 1 (Free Tier 1):
- LlamaNodes: 50% weight
- 1RPC: 50% weight
→ Used only if priority 10 nodes fail

Priority 0 (Fallback):
- PublicNode: 100% weight
→ Used only if all higher priority nodes fail
```

## Monitoring & Maintenance

### Check Node Status

Via API:
```bash
curl http://localhost:8080/api/v1/rpc-nodes/grouped | jq '.'
```

Via UI:
1. Go to Settings → RPC Nodes
2. Look at "Connectivity" column
3. Green badge = CONNECTED
4. Red badge = DISCONNECTED

### Manual Connection Check

Click the "Check Connection" button in the UI, or:

```bash
curl -X POST http://localhost:8080/api/v1/rpc-nodes/1/check
```

### Disable Problematic Nodes

In the UI:
1. Toggle the switch to disable
2. Node will be skipped in routing

Via API:
```bash
curl -X PUT http://localhost:8080/api/v1/rpc-nodes/1 \
  -H "Content-Type: application/json" \
  -d '{"is_enabled": false, ...other fields...}'
```

## Troubleshooting

### "Node shows as disconnected"

1. Verify URL is correct
2. Test manually:
   ```bash
   curl -X POST YOUR_RPC_URL \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
   ```
3. Check if API key is needed
4. Increase timeout value

### "Cannot create node"

1. Ensure chains table has the chain_id
2. Check foreign key constraints
3. Verify chain exists:
   ```bash
   curl http://localhost:8080/api/v1/chains
   ```

### "Frontend not showing nodes"

1. Check backend is running
2. Open browser console for errors
3. Verify API endpoints:
   ```bash
   curl http://localhost:8080/api/v1/rpc-nodes/grouped
   ```
4. Check CORS settings

## Best Practices

### 1. Use Multiple Nodes Per Chain

- Minimum 2 nodes for redundancy
- Mix free and paid for cost optimization
- Different providers for diversity

### 2. Set Appropriate Timeouts

- Public nodes: 30-60 seconds
- Paid nodes: 10-30 seconds
- High-traffic: 5-10 seconds

### 3. Regular Health Checks

- Run check-all daily
- Monitor connection status
- Replace failing nodes

### 4. Weight Distribution

- Premium nodes: 100% weight, high priority
- Free reliable: 50-100% weight, medium priority
- Experimental: 20-50% weight, low priority

### 5. Gradual Migration from DeBank

Week 1: Configure nodes, test with 10% traffic
Week 2: Increase to 50% traffic
Week 3: Monitor stability, increase to 80%
Week 4: Full migration, keep DeBank as fallback

## Next Steps

After setting up RPC nodes:

1. **Integrate with Data Provider**
   - Modify provider to use RPC nodes
   - Implement weighted selection
   - Add failover logic

2. **Add Background Health Checks**
   - Periodic connection testing
   - Auto-disable failing nodes
   - Alert on issues

3. **Implement Request Routing**
   - Select node based on weight
   - Track request counts
   - Balance load automatically

4. **Monitor Performance**
   - Track response times
   - Log errors per node
   - Optimize weight distribution

## Support

For issues or questions:
1. Check logs: `logs/app.log`
2. Review documentation: `docs/RPC_NODES.md`
3. Test with script: `scripts/test_rpc_nodes.sh`
4. Verify database: Check `rpc_nodes` table

## Summary

✅ Database table created
✅ Backend API ready
✅ Frontend UI implemented
✅ Connection testing working
✅ CRUD operations complete

You can now:
- Add RPC nodes for any chain
- Configure load balancing
- Monitor connection status
- Manage nodes through UI
- Prepare for DeBank migration
