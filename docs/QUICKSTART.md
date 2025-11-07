# 🚀 快速启动指南

## 最快启动方式 (Docker)

```bash
# 1. 启动所有服务
./scripts/quick-start.sh

# 2. 编辑配置文件添加 DeBank API key
# 编辑 config.yaml, 修改:
# debank.api_key: "YOUR_DEBANK_API_KEY"

# 3. 启动前端
cd frontend && npm run dev

# 4. 访问应用
open http://localhost:3000
```

## 服务说明

启动后会运行以下服务:

| 服务       | 端口   | 说明                 |
|----------|------|--------------------|
| MySQL    | 3306 | 数据库 (密码: rotki123) |
| Redis    | 6379 | 缓存 (可选)            |
| Backend  | 8080 | Go API 服务          |
| Frontend | 3000 | Vue.js UI          |

## 常用命令

### Docker 操作
```bash
# 查看运行状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 重启服务
docker-compose restart

# 停止所有服务
docker-compose down
```

### 数据库操作
```bash
# 连接 MySQL
docker exec -it rotki-mysql mysql -uroot -protki123 rotki_demo

# 查看 Redis
docker exec -it rotki-redis redis-cli
```

### 开发模式

如果你想本地运行代码 (只用 Docker 作为数据库):

```bash
# 1. 只启动数据库
make docker-db

# 2. 运行后端
go run cmd/server/main.go

# 3. 运行前端
cd frontend && npm run dev
```

## 测试应用

1. 打开 http://localhost:3000
2. 点击 "Add Wallet" 创建钱包
3. 点击 "Add Address" 添加以太坊地址
4. 系统会自动同步数据
5. 点击刷新按钮手动更新

## 示例地址

可以使用这些公开地址测试:

- Vitalik.eth: `0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045`
- Uniswap: `0x1a9C8182C09F50C8318d769245beA52c32BE35BC`

## 故障排查

### 后端无法连接数据库
```bash
# 检查 MySQL 是否启动
docker-compose ps mysql

# 查看 MySQL 日志
docker-compose logs mysql

# 等待更长时间
sleep 30 && docker-compose restart backend
```

### 端口冲突
如果 3306 或 8080 端口被占用:
```bash
# 修改 docker-compose.yml 中的端口映射
# mysql:
#   ports:
#     - "3307:3306"  # 改用 3307
```

### 清理所有数据重新开始
```bash
docker-compose down -v
./scripts/quick-start.sh
```

## 下一步

- 查看 [完整文档](../README.md)
- 了解 [架构设计](ARCHITECTURE.md)
- 阅读 [Docker 详细指南](DOCKER.md)
- 查看 [API 文档](../README.md#api-endpoints)

## 需要帮助?

- GitHub Issues: [报告问题]
- 文档: `docs/` 目录
- 配置示例: `config.yaml.example`
