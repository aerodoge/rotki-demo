# Rotki Demo - DeFi Asset Management System

A DeFi asset tracking application similar to Rotki, built with Go backend and Vue.js frontend, using DeBank API for blockchain data.

## Features

- **Wallet Management**: Create and manage multiple wallets
- **Address Tracking**: Add EVM addresses to wallets and track their assets
- **Real-time Data**: Sync token balances and DeFi positions from DeBank API
- **Auto-refresh**: Automatic periodic syncing of all addresses
- **Manual Refresh**: On-demand refresh for individual addresses or entire wallets
- **Asset Display**: View tokens, protocols, and total values across all chains
- **Extensible Architecture**: Provider interface allows easy switching from DeBank to custom data sources

## Architecture

### Backend (Go)
- **Web Framework**: Gin
- **Configuration**: Viper
- **Logging**: Zap
- **Database**: MySQL with GORM
- **Caching**: Redis (optional)
- **Rate Limiting**: Token bucket algorithm for API calls

### Frontend (Vue.js)
- **Framework**: Vue 3 with Composition API
- **State Management**: Pinia
- **HTTP Client**: Axios
- **Build Tool**: Vite

### Key Design Patterns

1. **Provider Interface Pattern**: Abstraction layer for data sources
   - Current: DeBank API
   - Future: Self-queried blockchain data, other APIs

2. **Repository Pattern**: Database access layer separation

3. **Service Layer**: Business logic encapsulation

4. **Rate Limiting**: Built-in protection against API rate limits

## Project Structure

```
rotki-demo/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── handler/             # HTTP request handlers
│   │   └── router/              # Route definitions
│   ├── config/                  # Configuration management
│   ├── database/                # Database initialization
│   ├── logger/                  # Logging setup
│   ├── models/                  # Database models
│   ├── provider/                # Data provider interface
│   │   └── debank/              # DeBank API implementation
│   ├── repository/              # Data access layer
│   └── service/                 # Business logic
├── frontend/
│   ├── src/
│   │   ├── api/                 # API client
│   │   ├── components/          # Vue components
│   │   ├── stores/              # Pinia stores
│   │   └── views/               # Page components
│   └── vite.config.js
├── docs/
│   └── database_schema.sql      # Database schema
├── config.yaml                  # Application configuration
└── go.mod                       # Go dependencies
```

## Getting Started

### Prerequisites

- Docker & Docker Compose (recommended)
- OR: Go 1.21+ & Node.js 18+ & MySQL 8.0+
- DeBank API Key (get from https://docs.cloud.debank.com)

### 🐳 Quick Start with Docker (Recommended)

The fastest way to get started:

```bash
# 1. Run the quick start script
./scripts/quick-start.sh

# 2. Start the frontend (in a new terminal)
cd frontend && npm run dev

# 3. Open http://localhost:3000
```

That's it! The script will:
- Start MySQL and Redis containers
- Start the backend service
- Install frontend dependencies
- Wait for all services to be ready

For more Docker options, see [Docker Guide](docs/DOCKER.md).

### 📦 Manual Setup (Without Docker)

If you prefer to run MySQL locally:

### Backend Setup

1. Install dependencies:
```bash
cd /Users/miles/go/src/rotki-demo
go mod download
```

2. Create database:
```bash
mysql -u root -p
CREATE DATABASE rotki_demo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

3. Import schema:
```bash
mysql -u root -p rotki_demo < migrations/001_initial_schema.sql
```

4. Configure application:
```bash
cp config.yaml config.yaml
# Edit config.yaml and add your DeBank API key
```

5. Run the server:
```bash
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080`

### Frontend Setup

1. Install dependencies:
```bash
cd frontend
npm install
```

2. Configure environment:
```bash
cp .env.example .env
# Edit .env if needed
```

3. Run development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

## 🐳 Docker Commands

```bash
# Start only databases (for local development)
make docker-db
# or: docker-compose up -d mysql redis

# Start all services (including backend)
make docker-up
# or: docker-compose up -d

# View logs
make docker-logs
# or: docker-compose logs -f

# Stop all services
make docker-down
# or: docker-compose down

# Clean everything (including data)
make docker-clean
# or: docker-compose down -v
```

For complete Docker documentation, see [docs/DOCKER.md](docs/DOCKER.md).

## 🔄 Starting and Stopping Services

### Stop Services

**Stop Backend (Docker):**
```bash
# Stop all Docker services
docker-compose down

# Or use Makefile
make docker-down
```

**Stop Frontend:**
```bash
# Press Ctrl+C in the terminal running the frontend
```

### Restart Services

**Restart Backend:**

```bash
# Option 1: Full restart (including database)
docker-compose down
docker-compose up -d

# Option 2: Restart backend only (keep database running)
docker-compose restart backend

# Option 3: Rebuild backend (after code changes)
docker-compose build backend
docker-compose up -d backend
```

**Restart Frontend:**
```bash
cd frontend
npm run dev
```

### Complete Restart (Clean All Data)

```bash
# Stop and remove all containers and volumes
docker-compose down -v

# Start fresh
docker-compose up -d

# Start frontend
cd frontend && npm run dev
```

### Quick Restart (Keep Data)

```bash
# Restart all Docker services
docker-compose restart

# Frontend will auto-reload (Vite hot reload)
```

### Check Service Status

```bash
# View all services
docker-compose ps

# View backend logs
docker-compose logs backend

# Follow logs in real-time
docker-compose logs -f backend

# Test backend health
curl http://localhost:8080/health

# Check frontend
# Open http://localhost:3000 in browser
```

### Development Workflow

**Daily Development (Run locally with Docker database):**
```bash
# Terminal 1: Start database only
make docker-db

# Terminal 2: Start backend
go run cmd/server/main.go

# Terminal 3: Start frontend
cd frontend && npm run dev
```

**Using Docker for Backend:**
```bash
# Terminal 1: Start all services
docker-compose up -d
docker-compose logs -f backend

# Terminal 2: Start frontend
cd frontend && npm run dev
```

**Stop Everything:**
```bash
# Stop Docker services
docker-compose down

# Stop frontend: Press Ctrl+C in terminal
```

### Makefile Commands

```bash
# Docker operations
make docker-up          # Start all services
make docker-down        # Stop all services
make docker-restart     # Restart all services
make docker-build       # Rebuild images
make docker-logs        # View logs
make docker-ps          # Check status
make docker-clean       # Clean all data (removes volumes)
make docker-db          # Start only MySQL + Redis

# Frontend operations
cd frontend
npm run dev             # Start dev server
npm run build           # Build for production

# Quick start script
./scripts/quick-start.sh  # One-command startup
./scripts/stop.sh         # Stop all services
```

## API Endpoints

### Wallets
- `GET /api/v1/wallets` - List all wallets
- `POST /api/v1/wallets` - Create a wallet
- `GET /api/v1/wallets/:id` - Get wallet details
- `PUT /api/v1/wallets/:id` - Update wallet
- `DELETE /api/v1/wallets/:id` - Delete wallet
- `POST /api/v1/wallets/:id/refresh` - Refresh all addresses in wallet

### Addresses
- `GET /api/v1/addresses` - List all addresses
- `GET /api/v1/addresses?wallet_id=:id` - List addresses by wallet
- `POST /api/v1/addresses` - Add an address
- `GET /api/v1/addresses/:id` - Get address details
- `DELETE /api/v1/addresses/:id` - Delete address
- `POST /api/v1/addresses/:id/refresh` - Refresh address data

## Configuration

### Database Configuration
```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: ""
  database: rotki_demo
  max_idle_conns: 10
  max_open_conns: 100
```

### DeBank API Configuration
```yaml
debank:
  api_key: "YOUR_API_KEY"
  base_url: "https://pro-openapi.debank.com"
  rate_limit:
    requests_per_second: 5
    burst: 10
  cache_ttl: 60
  timeout: 30
```

### Sync Configuration
```yaml
sync:
  enabled: true
  interval: 300        # Sync every 5 minutes
  batch_size: 10       # Process 10 addresses concurrently
```

## DeBank API Integration

### Rate Limiting Strategy
- Token bucket algorithm with configurable rate
- Default: 5 requests/second with burst of 10
- Automatic backoff on rate limit errors

### Cost Optimization
1. **Caching**: 60-second TTL on API responses
2. **Batch Requests**: Use `all_token_list` endpoint to get all chains at once
3. **Periodic Sync**: Configurable interval to avoid unnecessary calls
4. **On-demand Refresh**: Manual refresh only when needed

### API Endpoints Used
- `/v1/user/total_balance` - Get total value across all chains
- `/v1/user/all_token_list` - Get all tokens for an address
- `/v1/user/used_chain_list` - Get chains used by address
- `/v1/user/all_complex_protocol_list` - Get DeFi protocol positions

## Switching Data Providers

The application uses a provider interface pattern for easy switching:

### Current: DeBank Provider
```go
provider := debank.NewDeBankProvider(&cfg.DeBank)
```

### Future: Custom Provider
```go
// Implement the DataProvider interface
type CustomProvider struct {
    // Your implementation
}

func (p *CustomProvider) GetTokenList(ctx context.Context, address string, chainIDs []string) ([]TokenInfo, error) {
    // Query blockchain directly or use another API
}

// Use it
provider := NewCustomProvider(config)
```

The interface ensures all providers implement the same methods:
- `GetTotalBalance()`
- `GetTokenList()`
- `GetUsedChainList()`
- `GetProtocolList()`

## Monitoring and Logging

Logs are output using structured logging (Zap):
- Debug: Detailed request/response logs
- Info: Important events (server start, sync completed)
- Error: Errors and failures

Configure log level in `config.yaml`:
```yaml
log:
  level: debug  # debug, info, warn, error
  output: stdout
```

## Performance Considerations

1. **Database Indexing**: All foreign keys and frequently queried fields are indexed
2. **Batch Operations**: Token upserts use batch operations
3. **Concurrent Sync**: Multiple addresses synced in parallel
4. **Connection Pooling**: Database connection pool configured
5. **Rate Limiting**: Prevents API throttling

## Security Considerations

1. **Input Validation**: All user inputs validated
2. **SQL Injection**: GORM parameterized queries
3. **CORS**: Configurable CORS policies
4. **API Keys**: Stored in config, never committed to git
5. **Error Handling**: No sensitive data in error responses

## Future Enhancements

- [ ] Support for Bitcoin addresses
- [ ] Support for Solana addresses
- [ ] Historical balance tracking
- [ ] Chart visualizations
- [ ] Export to CSV/PDF
- [ ] Transaction history
- [ ] NFT tracking
- [ ] DeFi protocol details
- [ ] Custom tags and labels
- [ ] Multi-user support with authentication
- [ ] Webhook notifications

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
