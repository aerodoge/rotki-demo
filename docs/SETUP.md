# Quick Setup Guide

## Step 1: Database Setup

```bash
# Create database
mysql -u root -p
```

```sql
CREATE DATABASE rotki_demo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

## Step 2: Configuration

```bash
# Copy and edit configuration
cp config.yaml.example config.yaml
```

Edit `config.yaml` and update:
- Database credentials
- DeBank API key (get from https://docs.cloud.debank.com)

## Step 3: Backend Setup

```bash
# Install Go dependencies
go mod download

# Run the server (will auto-migrate database)
go run cmd/server/main.go
```

The API will be available at http://localhost:8080

Check health: `curl http://localhost:8080/health`

## Step 4: Frontend Setup

```bash
# Install frontend dependencies
cd frontend
npm install

# Run development server
npm run dev
```

The frontend will be available at http://localhost:3000

## Step 5: Test the Application

1. Open http://localhost:3000 in your browser
2. Click "Add Wallet" to create a wallet (e.g., "My Wallet")
3. Click "Add Address" to add an Ethereum address
   - Select the wallet you just created
   - Enter an Ethereum address (e.g., `0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb`)
   - Add a label (optional)
4. The system will automatically sync the address data from DeBank
5. Click the refresh button to manually update data

## Common Issues

### Database Connection Error
- Check MySQL is running: `mysql -u root -p`
- Verify credentials in `config.yaml`
- Ensure database exists: `SHOW DATABASES;`

### DeBank API Error
- Verify API key is correct in `config.yaml`
- Check rate limits aren't exceeded
- Ensure internet connection is working

### Frontend Can't Connect to Backend
- Verify backend is running on port 8080
- Check CORS settings in backend
- Verify proxy configuration in `vite.config.js`

## API Testing with curl

```bash
# Create a wallet
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{"name":"My Wallet","description":"Main wallet"}'

# List wallets
curl http://localhost:8080/api/v1/wallets

# Add an address
curl -X POST http://localhost:8080/api/v1/addresses \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_id": 1,
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "label": "Main Address",
    "chain_type": "EVM"
  }'

# Refresh an address
curl -X POST http://localhost:8080/api/v1/addresses/1/refresh

# Get address with tokens
curl http://localhost:8080/api/v1/addresses/1
```

## Production Deployment

### Build Backend
```bash
make build
# Binary will be at bin/rotki-demo
```

### Build Frontend
```bash
cd frontend
npm run build
# Static files will be in frontend/dist
```

### Run in Production
```bash
# Update config.yaml
# - Set server.mode to "release"
# - Set log.level to "info"
# - Configure proper database credentials
# - Set sync.enabled to true

# Run backend
./bin/rotki-demo

# Serve frontend with nginx or similar
```

## Development Tips

### Auto-restart Backend on Changes
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with air
air
```

### Database Migrations
The application automatically runs migrations on startup. To reset:
```bash
mysql -u root -p rotki_demo < docs/database_schema.sql
```

### View Logs
Backend logs go to stdout by default. To write to file:
```yaml
log:
  output: file
  file_path: logs/app.log
```

### Adjust Sync Interval
```yaml
sync:
  enabled: true
  interval: 300  # 5 minutes
  batch_size: 10
```
