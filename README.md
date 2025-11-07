# Rotki Demo - DeFi 资产管理系统

一个类似于 Rotki 的 DeFi 资产追踪应用，使用 Go 后端和 Vue.js 前端构建，通过 DeBank API 获取区块链数据。

## 功能特性

- **钱包管理**：创建和管理多个钱包
- **地址追踪**：添加 EVM 地址到钱包并追踪其资产
- **实时数据**：从 DeBank API 同步代币余额和 DeFi 持仓
- **自动刷新**：自动定期同步所有地址
- **手动刷新**：按需刷新单个地址或整个钱包
- **资产展示**：查看所有链上的代币、协议和总价值
- **可扩展架构**：提供商接口允许轻松从 DeBank 切换到自定义数据源

## 架构

### 后端 (Go)
- **Web 框架**：Gin
- **配置管理**：Viper
- **日志**：Zap
- **数据库**：MySQL + GORM
- **缓存**：Redis (可选)
- **速率限制**：令牌桶算法用于 API 调用

### 前端 (Vue.js)
- **框架**：Vue 3 + Composition API
- **状态管理**：Pinia
- **HTTP 客户端**：Axios
- **构建工具**：Vite

### 核心设计模式

1. **提供商接口模式**：数据源抽象层
   - 当前：DeBank API
   - 未来：自查询区块链数据、其他 API

2. **仓储模式**：数据库访问层分离

3. **服务层**：业务逻辑封装

4. **速率限制**：内置 API 速率限制保护

## 项目结构

```
rotki-demo/
├── cmd/
│   └── server/
│       └── main.go              # 应用程序入口点
├── internal/
│   ├── api/
│   │   ├── handler/             # HTTP 请求处理器
│   │   └── router/              # 路由定义
│   ├── config/                  # 配置管理
│   ├── database/                # 数据库初始化
│   ├── logger/                  # 日志设置
│   ├── models/                  # 数据库模型
│   ├── provider/                # 数据提供商接口
│   │   └── debank/              # DeBank API 实现
│   ├── repository/              # 数据访问层
│   └── service/                 # 业务逻辑
├── frontend/
│   ├── src/
│   │   ├── api/                 # API 客户端
│   │   ├── components/          # Vue 组件
│   │   ├── stores/              # Pinia 状态存储
│   │   └── views/               # 页面组件
│   └── vite.config.js
├── docs/
│   └── database_schema.sql      # 数据库模式
├── config.yaml                  # 应用程序配置
└── go.mod                       # Go 依赖
```

## 开始使用

### 前置要求

- Docker & Docker Compose（推荐）
- 或者：Go 1.21+ & Node.js 18+ & MySQL 8.0+
- DeBank API 密钥（从 https://docs.cloud.debank.com 获取）

### 🐳 使用 Docker 快速启动（推荐）

最快的启动方式：

```bash
# 1. 运行快速启动脚本
./scripts/quick-start.sh

# 2. 启动前端（在新终端中）
cd frontend && npm run dev

# 3. 打开 http://localhost:3000
```

就这么简单！脚本将会：
- 启动 MySQL 和 Redis 容器
- 启动后端服务
- 安装前端依赖
- 等待所有服务准备就绪

更多 Docker 选项，请查看 [Docker 指南](docs/DOCKER.md)。

### 📦 手动设置（不使用 Docker）

如果你更喜欢在本地运行 MySQL：

### 后端设置

1. 安装依赖：
```bash
cd /Users/miles/go/src/rotki-demo
go mod download
```

2. 创建数据库：
```bash
mysql -u root -p
CREATE DATABASE rotki_demo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

3. 导入模式：
```bash
mysql -u root -p rotki_demo < migrations/001_initial_schema.sql
```

4. 配置应用程序：
```bash
cp config.yaml config.yaml
# 编辑 config.yaml 并添加你的 DeBank API 密钥
```

5. 运行服务器：
```bash
go run cmd/server/main.go
```

API 将在 `http://localhost:8080` 可用

### 前端设置

1. 安装依赖：
```bash
cd frontend
npm install
```

2. 配置环境：
```bash
cp .env.example .env
# 如果需要，编辑 .env
```

3. 运行开发服务器：
```bash
npm run dev
```

前端将在 `http://localhost:3000` 可用

## 🐳 Docker 命令

```bash
# 仅启动数据库（用于本地开发）
make docker-db
# 或：docker-compose up -d mysql redis

# 启动所有服务（包括后端）
make docker-up
# 或：docker-compose up -d

# 查看日志
make docker-logs
# 或：docker-compose logs -f

# 停止所有服务
make docker-down
# 或：docker-compose down

# 清理所有内容（包括数据）
make docker-clean
# 或：docker-compose down -v
```

完整的 Docker 文档，请参见 [docs/DOCKER.md](docs/DOCKER.md)。

## 🔄 启动和停止服务

### 停止服务

**停止后端（Docker）：**
```bash
# 停止所有 Docker 服务
docker-compose down

# 或使用 Makefile
make docker-down
```

**停止前端：**
```bash
# 在运行前端的终端中按 Ctrl+C
```

### 重启服务

**重启后端：**

```bash
# 选项 1：完全重启（包括数据库）
docker-compose down
docker-compose up -d

# 选项 2：仅重启后端（保持数据库运行）
docker-compose restart backend

# 选项 3：重新构建后端（代码更改后）
docker-compose build backend
docker-compose up -d backend
```

**重启前端：**
```bash
cd frontend
npm run dev
```

### 完全重启（清除所有数据）

```bash
# 停止并删除所有容器和卷
docker-compose down -v

# 全新启动
docker-compose up -d

# 启动前端
cd frontend && npm run dev
```

### 快速重启（保留数据）

```bash
# 重启所有 Docker 服务
docker-compose restart

# 前端将自动重载（Vite 热重载）
```

### 检查服务状态

```bash
# 查看所有服务
docker-compose ps

# 查看后端日志
docker-compose logs backend

# 实时跟踪日志
docker-compose logs -f backend

# 测试后端健康状态
curl http://localhost:8080/health

# 检查前端
# 在浏览器中打开 http://localhost:3000
```

### 开发工作流

**日常开发（使用 Docker 数据库本地运行）：**
```bash
# 终端 1：仅启动数据库
make docker-db

# 终端 2：启动后端
go run cmd/server/main.go

# 终端 3：启动前端
cd frontend && npm run dev
```

**使用 Docker 运行后端：**
```bash
# 终端 1：启动所有服务
docker-compose up -d
docker-compose logs -f backend

# 终端 2：启动前端
cd frontend && npm run dev
```

**停止所有服务：**
```bash
# 停止 Docker 服务
docker-compose down

# 停止前端：在终端中按 Ctrl+C
```

### Makefile 命令

```bash
# Docker 操作
make docker-up          # 启动所有服务
make docker-down        # 停止所有服务
make docker-restart     # 重启所有服务
make docker-build       # 重新构建镜像
make docker-logs        # 查看日志
make docker-ps          # 检查状态
make docker-clean       # 清理所有数据（删除卷）
make docker-db          # 仅启动 MySQL + Redis

# 前端操作
cd frontend
npm run dev             # 启动开发服务器
npm run build           # 构建生产版本

# 快速启动脚本
./scripts/quick-start.sh  # 一键启动
./scripts/stop.sh         # 停止所有服务
```

## API 端点

### 钱包
- `GET /api/v1/wallets` - 列出所有钱包
- `POST /api/v1/wallets` - 创建钱包
- `GET /api/v1/wallets/:id` - 获取钱包详情
- `PUT /api/v1/wallets/:id` - 更新钱包
- `DELETE /api/v1/wallets/:id` - 删除钱包
- `POST /api/v1/wallets/:id/refresh` - 刷新钱包中的所有地址

### 地址
- `GET /api/v1/addresses` - 列出所有地址
- `GET /api/v1/addresses?wallet_id=:id` - 按钱包列出地址
- `POST /api/v1/addresses` - 添加地址
- `GET /api/v1/addresses/:id` - 获取地址详情
- `DELETE /api/v1/addresses/:id` - 删除地址
- `POST /api/v1/addresses/:id/refresh` - 刷新地址数据

## 配置

### 数据库配置
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

### DeBank API 配置
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

### 同步配置
```yaml
sync:
  enabled: true
  interval: 300        # 每 5 分钟同步一次
  batch_size: 10       # 并发处理 10 个地址
```

## DeBank API 集成

### 速率限制策略
- 可配置速率的令牌桶算法
- 默认：每秒 5 个请求，突发 10 个
- 速率限制错误时自动退避

### 成本优化
1. **缓存**：API 响应的 60 秒 TTL
2. **批量请求**：使用 `all_token_list` 端点一次获取所有链
3. **定期同步**：可配置间隔以避免不必要的调用
4. **按需刷新**：仅在需要时手动刷新

### 使用的 API 端点
- `/v1/user/total_balance` - 获取所有链的总价值
- `/v1/user/all_token_list` - 获取地址的所有代币
- `/v1/user/used_chain_list` - 获取地址使用的链
- `/v1/user/all_complex_protocol_list` - 获取 DeFi 协议持仓

## 切换数据提供商

应用程序使用提供商接口模式以便于切换：

### 当前：DeBank 提供商
```go
provider := debank.NewDeBankProvider(&cfg.DeBank)
```

### 未来：自定义提供商
```go
// 实现 DataProvider 接口
type CustomProvider struct {
    // 你的实现
}

func (p *CustomProvider) GetTokenList(ctx context.Context, address string, chainIDs []string) ([]TokenInfo, error) {
    // 直接查询区块链或使用其他 API
}

// 使用它
provider := NewCustomProvider(config)
```

接口确保所有提供商实现相同的方法：
- `GetTotalBalance()`
- `GetTokenList()`
- `GetUsedChainList()`
- `GetProtocolList()`

## 监控和日志

使用结构化日志（Zap）输出日志：
- Debug：详细的请求/响应日志
- Info：重要事件（服务器启动、同步完成）
- Error：错误和失败

在 `config.yaml` 中配置日志级别：
```yaml
log:
  level: debug  # debug, info, warn, error
  output: stdout
```

## 性能考虑

1. **数据库索引**：所有外键和频繁查询的字段都已建立索引
2. **批量操作**：代币更新使用批量操作
3. **并发同步**：多个地址并行同步
4. **连接池**：数据库连接池已配置
5. **速率限制**：防止 API 限流

## 安全考虑

1. **输入验证**：所有用户输入都经过验证
2. **SQL 注入**：GORM 参数化查询
3. **CORS**：可配置的 CORS 策略
4. **API 密钥**：存储在配置中，永不提交到 git
5. **错误处理**：错误响应中不包含敏感数据

## 未来增强

- [ ] 支持比特币地址
- [ ] 支持 Solana 地址
- [ ] 历史余额跟踪
- [ ] 图表可视化
- [ ] 导出到 CSV/PDF
- [ ] 交易历史
- [ ] NFT 跟踪
- [ ] DeFi 协议详情
- [ ] 自定义标签
- [ ] 多用户支持与身份验证
- [ ] Webhook 通知

## 许可证

MIT

## 贡献

欢迎贡献！请随时提交 Pull Request。
