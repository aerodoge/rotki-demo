# Docker 部署指南

## 快速启动

### 方式一: 只启动数据库 (开发模式)

如果你想本地运行 Go 和前端,只需要 MySQL 和 Redis:

```bash
# 1. 只启动数据库服务
docker-compose up -d mysql redis

# 2. 等待数据库启动完成 (约 10-20 秒)
docker-compose ps

# 3. 配置后端
cp config.yaml.example config.yaml
# 编辑 config.yaml:
# - database.password: rotki123
# - debank.api_key: 你的 API key

# 4. 运行后端
go run cmd/server/main.go

# 5. 运行前端 (新终端)
cd frontend && npm run dev
```

访问: http://localhost:3000

### 方式二: 完整 Docker 部署

启动所有服务 (MySQL + Redis + Backend):

```bash
# 1. 准备配置
cp config.docker.yaml config.yaml
# 编辑 config.yaml 添加你的 DeBank API key

# 2. 启动所有服务
docker-compose up -d

# 3. 查看日志
docker-compose logs -f backend

# 4. 前端本地运行
cd frontend
npm install
npm run dev
```

访问: http://localhost:3000

## Docker Compose 服务

### MySQL
- **端口**: 3306
- **数据库**: rotki_demo
- **用户**: root
- **密码**: rotki123
- **数据持久化**: Docker volume `mysql_data`
- **自动初始化**: 使用 `docs/database_schema.sql`

### Redis
- **端口**: 6379
- **无密码**
- **数据持久化**: Docker volume `redis_data`

### Backend (可选)
- **端口**: 8080
- **自动连接**: MySQL + Redis
- **配置**: 通过 volume 挂载 `config.yaml`
- **日志**: 挂载到 `./logs` 目录

## 常用命令

### 启动服务
```bash
# 启动所有服务
docker-compose up -d

# 只启动数据库
docker-compose up -d mysql redis

# 启动并查看日志
docker-compose up
```

### 查看状态
```bash
# 查看运行状态
docker-compose ps

# 查看日志
docker-compose logs

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f mysql
```

### 停止服务
```bash
# 停止所有服务
docker-compose stop

# 停止并删除容器
docker-compose down

# 停止并删除容器和数据卷 (会删除数据!)
docker-compose down -v
```

### 重启服务
```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart backend
```

### 数据库操作
```bash
# 连接到 MySQL
docker exec -it rotki-mysql mysql -uroot -protki123 rotki_demo

# 导出数据库
docker exec rotki-mysql mysqldump -uroot -protki123 rotki_demo > backup.sql

# 导入数据库
docker exec -i rotki-mysql mysql -uroot -protki123 rotki_demo < backup.sql

# 查看 Redis 数据
docker exec -it rotki-redis redis-cli
```

### 更新代码
```bash
# 重新构建并启动
docker-compose up -d --build backend
```

## 配置说明

### config.yaml (Docker 环境)

```yaml
database:
  host: mysql          # Docker 服务名
  port: 3306
  username: root
  password: rotki123   # 与 docker-compose.yml 一致
  database: rotki_demo

redis:
  host: redis          # Docker 服务名
  port: 6379

debank:
  api_key: "YOUR_KEY"  # 必须配置!
```

### 环境变量 (可选)

你也可以通过环境变量覆盖配置:

```yaml
# docker-compose.yml
services:
  backend:
    environment:
      - DB_HOST=mysql
      - DB_PASSWORD=rotki123
      - DEBANK_API_KEY=your_key
```

## 生产部署建议

### 1. 修改默认密码
```yaml
# docker-compose.yml
services:
  mysql:
    environment:
      MYSQL_ROOT_PASSWORD: "strong_password_here"
```

同时更新 `config.yaml`:
```yaml
database:
  password: "strong_password_here"
```

### 2. 添加 Redis 密码
```yaml
# docker-compose.yml
services:
  redis:
    command: redis-server --requirepass your_redis_password
```

### 3. 使用 Docker Secrets

创建 `docker-compose.prod.yml`:
```yaml
version: '3.8'

services:
  mysql:
    environment:
      MYSQL_ROOT_PASSWORD_FILE: /run/secrets/db_password
    secrets:
      - db_password

secrets:
  db_password:
    file: ./secrets/db_password.txt
```

### 4. 配置反向代理 (Nginx)

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location / {
        root /path/to/frontend/dist;
        try_files $uri $uri/ /index.html;
    }
}
```

### 5. 限制资源使用

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

## 故障排查

### 后端无法连接数据库

**问题**: `Error 2005: Unknown MySQL server host 'mysql'`

**解决**:
```bash
# 检查 MySQL 是否启动
docker-compose ps mysql

# 查看 MySQL 日志
docker-compose logs mysql

# 等待 MySQL 完全启动
docker-compose up -d mysql
sleep 20
docker-compose up -d backend
```

### 数据库连接被拒绝

**问题**: `Error 1045: Access denied`

**解决**: 检查密码是否一致
```bash
# docker-compose.yml 中的密码
MYSQL_ROOT_PASSWORD: rotki123

# config.yaml 中的密码
database:
  password: rotki123  # 必须一致!
```

### 端口冲突

**问题**: `Bind for 0.0.0.0:3306 failed: port is already allocated`

**解决**: 修改端口映射
```yaml
services:
  mysql:
    ports:
      - "3307:3306"  # 使用 3307 映射到容器的 3306
```

然后更新 config.yaml:
```yaml
database:
  port: 3307
```

### 查看容器内部

```bash
# 进入 backend 容器
docker exec -it rotki-backend sh

# 查看配置文件
cat config.yaml

# 测试数据库连接
ping mysql
```

## 数据持久化

### 数据存储位置
- MySQL 数据: Docker volume `mysql_data`
- Redis 数据: Docker volume `redis_data`

### 查看 volume
```bash
# 列出所有 volume
docker volume ls

# 查看 volume 详情
docker volume inspect rotki-demo_mysql_data
```

### 备份数据
```bash
# 备份 MySQL
docker exec rotki-mysql mysqldump -uroot -protki123 rotki_demo > backup_$(date +%Y%m%d).sql

# 备份 Redis
docker exec rotki-redis redis-cli SAVE
docker cp rotki-redis:/data/dump.rdb ./redis_backup_$(date +%Y%m%d).rdb
```

### 恢复数据
```bash
# 恢复 MySQL
docker exec -i rotki-mysql mysql -uroot -protki123 rotki_demo < backup_20240101.sql

# 恢复 Redis
docker cp redis_backup_20240101.rdb rotki-redis:/data/dump.rdb
docker-compose restart redis
```

## 监控和日志

### 日志收集
```bash
# 持续查看所有日志
docker-compose logs -f

# 只查看最近 100 行
docker-compose logs --tail=100

# 保存日志到文件
docker-compose logs > docker_logs.txt
```

### 资源监控
```bash
# 查看资源使用
docker stats

# 只查看特定容器
docker stats rotki-backend rotki-mysql
```

## 开发 vs 生产

### 开发环境
```bash
# 使用 docker-compose.yml
docker-compose up -d mysql redis

# 本地运行代码,支持热重载
go run cmd/server/main.go
cd frontend && npm run dev
```

### 生产环境
```bash
# 构建生产版本
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# 或使用 Makefile
make docker-build
make docker-run
```

## 完整启动流程

```bash
# 1. 克隆项目
git clone <repo>
cd rotki-demo

# 2. 准备配置
cp config.docker.yaml config.yaml
# 编辑 config.yaml 添加 DeBank API key

# 3. 启动数据库
docker-compose up -d mysql redis

# 4. 等待启动完成
sleep 20

# 5. 初始化数据库 (自动完成)
# schema 会自动加载

# 6. 启动后端
docker-compose up -d backend

# 7. 安装前端依赖
cd frontend
npm install

# 8. 启动前端
npm run dev

# 9. 访问应用
open http://localhost:3000
```

搞定! 🎉
