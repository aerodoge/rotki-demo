# Architecture Documentation

## System Overview

Rotki Demo is a full-stack DeFi asset management application that tracks cryptocurrency holdings across multiple blockchain addresses. It uses DeBank API as the data source but is designed with an abstraction layer to easily switch to other data providers.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────┐
│                     Frontend (Vue.js)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   Sidebar    │  │ EVM Accounts │  │   Modals     │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│  ┌─────────────────────────────────────────────────┐  │
│  │           Pinia Store (State Management)         │  │
│  └─────────────────────────────────────────────────┘  │
│  ┌─────────────────────────────────────────────────┐  │
│  │              Axios (HTTP Client)                 │  │
│  └─────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                           │
                           │ HTTP/JSON
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    Backend (Go/Gin)                      │
│  ┌─────────────────────────────────────────────────┐  │
│  │              HTTP Router (Gin)                   │  │
│  │    /api/v1/wallets  /api/v1/addresses           │  │
│  └─────────────────────────────────────────────────┘  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   Handlers   │  │   Services   │  │ Repositories │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│  ┌─────────────────────────────────────────────────┐  │
│  │         Data Provider Interface                  │  │
│  │  ┌────────────────┐   ┌─────────────────────┐  │  │
│  │  │ DeBank Client  │   │  Future: Custom RPC │  │  │
│  │  └────────────────┘   └─────────────────────┘  │  │
│  └─────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
           │                           │
           │                           │
           ▼                           ▼
    ┌────────────┐           ┌──────────────────┐
    │   MySQL    │           │   DeBank API     │
    │  Database  │           │ (Rate Limited)   │
    └────────────┘           └──────────────────┘
```

## Backend Architecture

### Layer Breakdown

#### 1. API Layer (`internal/api`)
- **Handler**: Request/response handling, validation
- **Router**: Route definitions and middleware setup
- **Responsibilities**: HTTP concerns only

#### 2. Service Layer (`internal/service`)
- **SyncService**: Background synchronization logic
- **Business Logic**: Orchestrates operations across multiple repositories
- **Transaction Management**: Ensures data consistency

#### 3. Repository Layer (`internal/repository`)
- **WalletRepository**: Wallet CRUD operations
- **AddressRepository**: Address CRUD operations
- **TokenRepository**: Token data management
- **Pure Data Access**: No business logic

#### 4. Provider Layer (`internal/provider`)
- **Interface Definition**: DataProvider interface
- **DeBank Implementation**: Current data source
- **Extensibility**: Easy to add new providers

### Key Design Patterns

#### Provider Pattern
```go
type DataProvider interface {
    GetTotalBalance(ctx context.Context, address string) (*TotalBalanceResponse, error)
    GetTokenList(ctx context.Context, address string, chainIDs []string) ([]TokenInfo, error)
    // ... other methods
}
```

**Benefits:**
- Decouples data source from business logic
- Easy to switch between DeBank, custom RPC, or other APIs
- Testable with mock providers

#### Repository Pattern
```go
type WalletRepository struct {
    db *gorm.DB
}

func (r *WalletRepository) GetByID(id uint) (*models.Wallet, error) {
    // Database access only
}
```

**Benefits:**
- Separates data access from business logic
- Easier to test
- Can switch database implementations

#### Service Pattern
```go
type SyncService struct {
    dataProvider provider.DataProvider
    addressRepo  *repository.AddressRepository
    tokenRepo    *repository.TokenRepository
}
```

**Benefits:**
- Encapsulates complex business logic
- Coordinates between multiple repositories and providers
- Manages transactions and error handling

## Frontend Architecture

### Component Structure

```
App.vue
├── Sidebar.vue (Navigation)
└── EVMAccounts.vue (Main view)
    ├── Wallet rows (expandable)
    └── Address rows (nested)
        └── Token displays
```

### State Management (Pinia)

```javascript
walletStore
├── state
│   ├── wallets[]
│   ├── addresses[]
│   ├── selectedWallet
│   └── loading
├── getters
│   ├── getAddressesByWallet()
│   ├── getTotalValueByWallet()
│   └── getTotalValue()
└── actions
    ├── fetchWallets()
    ├── createWallet()
    ├── refreshWallet()
    └── ...
```

## Data Flow

### Adding a New Address

1. **User Action**: Clicks "Add Address" button
2. **Frontend**: Opens modal, user fills form
3. **API Call**: POST /api/v1/addresses
4. **Backend Handler**: Validates request
5. **Repository**: Inserts address into database
6. **Background Trigger**: Initiates sync for new address
7. **Provider**: Calls DeBank API
8. **Data Storage**: Stores tokens in database
9. **Response**: Returns created address to frontend
10. **UI Update**: Displays new address with tokens

### Automatic Sync Process

```
SyncService.Start()
    │
    ├─> Ticker (every 5 minutes)
    │
    ├─> GetAllNeedingSync()
    │     │
    │     └─> Query addresses where last_synced_at > interval
    │
    ├─> ProcessInBatches()
    │     │
    │     └─> Process 10 addresses concurrently
    │
    └─> For each address:
          │
          ├─> GetTokenList() from Provider
          │
          ├─> UpsertBatch() to database
          │
          └─> UpdateLastSynced()
```

## Database Schema

### Core Tables

**wallets**
- Stores wallet metadata
- One-to-many with addresses

**addresses**
- Stores blockchain addresses
- Links to wallet
- Tracks last sync time

**tokens**
- Stores token balances
- Links to address and chain
- Updated on each sync

**asset_snapshots**
- Historical snapshots
- For tracking balance changes over time

### Relationships

```
wallets (1) ──< (N) addresses (1) ──< (N) tokens
                                  │
                                  └──< (N) asset_snapshots
```

## Rate Limiting Strategy

### Token Bucket Algorithm

```go
limiter := rate.NewLimiter(
    rate.Limit(5),  // 5 requests per second
    10,             // Burst of 10
)

// Before each request
limiter.Wait(ctx)
```

### Benefits
- Prevents API throttling
- Smooth traffic distribution
- Configurable per environment

### Cost Optimization

1. **Caching**: 60s TTL on responses
2. **Batch Endpoints**: Use `all_token_list` instead of per-chain calls
3. **Periodic Sync**: Configurable interval (default 5 minutes)
4. **On-demand Refresh**: Only when user requests

## Scalability Considerations

### Current Limitations
- Single server instance
- In-memory rate limiting
- No distributed caching

### Future Improvements

1. **Horizontal Scaling**
   - Add Redis for distributed rate limiting
   - Add Redis for caching
   - Load balancer for multiple instances

2. **Database Optimization**
   - Read replicas for queries
   - Partitioning for large tables
   - Materialized views for aggregations

3. **Async Processing**
   - Message queue (RabbitMQ/Kafka) for sync jobs
   - Worker pools for parallel processing
   - Job status tracking

## Security Considerations

### Current Implementation
- Input validation on all endpoints
- GORM parameterized queries (SQL injection protection)
- CORS middleware with configurable origins
- API keys stored in config (not in code)

### Production Requirements
- [ ] Add authentication (JWT)
- [ ] Rate limiting per user
- [ ] Request signing for API calls
- [ ] Audit logging
- [ ] HTTPS/TLS enforcement
- [ ] Secrets management (Vault)

## Testing Strategy

### Recommended Tests

**Unit Tests**
```go
// Repository tests
func TestWalletRepository_Create(t *testing.T) {
    // Test database operations
}

// Service tests with mocks
func TestSyncService_SyncAddress(t *testing.T) {
    mockProvider := &MockProvider{}
    // Test business logic
}
```

**Integration Tests**
```go
// API tests
func TestWalletAPI_CreateWallet(t *testing.T) {
    // Test full HTTP flow
}
```

**Frontend Tests**
```javascript
// Component tests
describe('EVMAccounts', () => {
  it('displays wallets correctly', () => {
    // Test component rendering
  })
})
```

## Deployment

### Development
```bash
# Backend
go run cmd/server/main.go

# Frontend
cd frontend && npm run dev
```

### Production
```bash
# Build backend
make build

# Build frontend
cd frontend && npm run build

# Run
./bin/rotki-demo
```

### Docker (Recommended)
```bash
docker-compose up -d
```

## Monitoring

### Metrics to Track
- API response times
- DeBank API call count
- Sync success/failure rate
- Database query performance
- Token balance discrepancies

### Logging
- Structured logging with Zap
- Log levels: DEBUG, INFO, WARN, ERROR
- Request/response logging
- Error stack traces

## Future Enhancements

### Phase 2
- [ ] Support for Bitcoin/Solana addresses
- [ ] Historical balance charts
- [ ] Transaction history
- [ ] NFT tracking

### Phase 3
- [ ] Multi-user support
- [ ] Authentication/authorization
- [ ] Custom alerts
- [ ] Export functionality

### Phase 4
- [ ] Self-hosted blockchain nodes
- [ ] Custom data provider
- [ ] Remove DeBank dependency
- [ ] Advanced analytics
