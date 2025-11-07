# 快速设置指南

## 步骤 1：数据库设置

```bash
# 创建数据库
mysql -u root -p
```

```sql
CREATE DATABASE rotki_demo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

## 步骤 2：配置

```bash
# 复制并编辑配置
cp config.yaml.example config.yaml
```

编辑 `config.yaml` 并更新：
- 数据库凭据
- DeBank API 密钥（从 https://docs.cloud.debank.com 获取）

## 步骤 3：后端设置

```bash
# 安装 Go 依赖
go mod download

# 运行服务器（将自动迁移数据库）
go run cmd/server/main.go
```

API 将在 http://localhost:8080 可用

检查健康状态：`curl http://localhost:8080/health`

## 步骤 4：前端设置

```bash
# 安装前端依赖
cd frontend
npm install

# 运行开发服务器
npm run dev
```

前端将在 http://localhost:3000 可用

## 步骤 5：测试应用程序

1. 在浏览器中打开 http://localhost:3000
2. 点击"Add Wallet"创建钱包（例如，"My Wallet"）
3. 点击"Add Address"添加以太坊地址
   - 选择刚刚创建的钱包
   - 输入以太坊地址（例如，`0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb`）
   - 添加标签（可选）
4. 系统将自动从 DeBank 同步地址数据
5. 点击刷新按钮手动更新数据

## 常见问题

### 数据库连接错误
- 检查 MySQL 是否运行：`mysql -u root -p`
- 验证 `config.yaml` 中的凭据
- 确保数据库存在：`SHOW DATABASES;`

### DeBank API 错误
- 验证 `config.yaml` 中的 API 密钥是否正确
- 检查速率限制是否超出
- 确保互联网连接正常

### 前端无法连接到后端
- 验证后端是否在 8080 端口运行
- 检查后端的 CORS 设置
- 验证 `vite.config.js` 中的代理配置

## 使用 curl 进行 API 测试

```bash
# 创建钱包
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{"name":"My Wallet","description":"Main wallet"}'

# 列出钱包
curl http://localhost:8080/api/v1/wallets

# 添加地址
curl -X POST http://localhost:8080/api/v1/addresses \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_id": 1,
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "label": "Main Address",
    "chain_type": "EVM"
  }'

# 刷新地址
curl -X POST http://localhost:8080/api/v1/addresses/1/refresh

# 获取带代币的地址
curl http://localhost:8080/api/v1/addresses/1
```

## 生产部署

### 构建后端
```bash
make build
# 二进制文件将在 bin/rotki-demo
```

### 构建前端
```bash
cd frontend
npm run build
# 静态文件将在 frontend/dist
```

### 在生产环境中运行
```bash
# 更新 config.yaml
# - 将 server.mode 设置为 "release"
# - 将 log.level 设置为 "info"
# - 配置正确的数据库凭据
# - 将 sync.enabled 设置为 true

# 运行后端
./bin/rotki-demo

# 使用 nginx 或类似工具提供前端服务
```

## 开发技巧

### 代码更改时自动重启后端
```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 使用 air 运行
air
```

### 数据库迁移
应用程序在启动时自动运行迁移。要重置：
```bash
mysql -u root -p rotki_demo < docs/database_schema.sql
```

### 查看日志
默认情况下，后端日志输出到标准输出。要写入文件：
```yaml
log:
  output: file
  file_path: logs/app.log
```

### 调整同步间隔
```yaml
sync:
  enabled: true
  interval: 300  # 5 分钟
  batch_size: 10
```
