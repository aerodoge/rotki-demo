# Rotki Demo 运行指南

## 目录
- [系统架构](#系统架构)
- [启动服务](#启动服务)
- [ngrok 配置](#ngrok-配置)
- [访问地址](#访问地址)
- [停止服务](#停止服务)
- [常见问题](#常见问题)

## 系统架构

本项目由三个主要服务组成：

1. 前端开发服务器 (Vite)
- 进程: node vite
- 启动命令: cd frontend && npm run dev
- 端口: 3000
- 作用: 运行 Vue.js 前端应用

2. 后端 API 服务器 (Go)
- 进程: go run cmd/server/main.go
- 启动命令: go run cmd/server/main.go
- 端口: 8080
- 作用: 运行 Go 后端 API 服务

3. ngrok 隧道 (公网访问)
- 进程: ngrok http 3000
- 启动命令: ngrok http 3000 --log=stdout
- 端口: 转发 3000 端口到公网
- 公网地址: https://sodless-judgmentally-kelsie.ngrok-free.dev
- 作用: 让公司内外的人都能通过公网访问前端

## 启动服务

### 方式一：手动启动（推荐用于开发）

需要打开 3 个终端窗口，分别运行以下命令：

#### 1. 启动后端服务 (终端 1)

```bash
cd /Users/miles/go/src/rotki-demo
go run cmd/server/main.go
```

**输出**：
- 服务监听在 `:8080`
- 显示 API 路由列表
- 数据库连接成功提示

**健康检查**：
```bash
curl http://localhost:8080/health
```

#### 2. 启动前端服务 (终端 2)

```bash
cd /Users/miles/go/src/rotki-demo/frontend
npm run dev
```

**输出**：
```
VITE v5.4.21  ready in 163 ms

➜  Local:   http://localhost:3000/
➜  Network: http://192.168.31.151:3000/
```

**注意**：前端会自动热更新，修改代码后会实时刷新。

#### 3. 启动 ngrok 隧道 (终端 3)

```bash
ngrok http 3000 --log=stdout
```

**输出**：
```
started tunnel obj=tunnels name=command_line
addr=http://localhost:3000
url=https://sodless-judgmentally-kelsie.ngrok-free.dev
```

**注意**：每次重启 ngrok，URL 可能会变化（免费版）。

### 方式二：使用后台进程

如果需要让服务在后台运行：

```bash
# 启动后端（后台）
cd /Users/miles/go/src/rotki-demo
nohup go run cmd/server/main.go > backend.log 2>&1 &

# 启动前端（后台）
cd /Users/miles/go/src/rotki-demo/frontend
nohup npm run dev > frontend.log 2>&1 &

# 启动 ngrok（后台）
nohup ngrok http 3000 --log=stdout > ngrok.log 2>&1 &
```

查看日志：
```bash
tail -f backend.log
tail -f frontend.log
tail -f ngrok.log
```

## ngrok 配置

### 安装 ngrok

```bash
brew install ngrok/ngrok/ngrok
```

### 配置 authtoken

1. 注册 ngrok 账号：https://dashboard.ngrok.com/signup
2. 获取 authtoken：https://dashboard.ngrok.com/get-started/your-authtoken
3. 配置 authtoken：

```bash
ngrok config add-authtoken YOUR_AUTHTOKEN
```

### ngrok 基本用法

#### 暴露本地端口到公网

```bash
# 暴露 3000 端口（前端）
ngrok http 3000

# 暴露 8080 端口（后端）
ngrok http 8080

# 带日志输出
ngrok http 3000 --log=stdout

# 自定义子域名（需要付费版）
ngrok http 3000 --subdomain=myapp
```

#### 查看 ngrok 状态

访问本地管理界面：http://localhost:4040

在这里可以看到：
- 当前的公网 URL
- 所有 HTTP 请求/响应
- 请求详情和日志

#### 停止 ngrok

```bash
# 如果在前台运行，按 Ctrl+C

# 如果在后台运行，找到进程 ID
ps aux | grep ngrok
kill <PID>
```

### ngrok 免费版限制

- 每次启动 URL 会变化
- 只能同时运行 1 个隧道
- 每分钟请求数限制
- 会话超时限制（8 小时）

### 同时暴露前后端（需要付费版）

免费版只能运行一个隧道，如果需要同时暴露前后端，可以：

**方案 1**：只暴露前端，后端通过 Vite proxy 转发
```bash
# frontend/vite.config.js 已配置 proxy
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true
  }
}

# 只需要暴露前端
ngrok http 3000
```

**方案 2**：升级到 ngrok 付费版
- 可以同时运行多个隧道
- 固定的子域名
- 更高的请求限制

**方案 3**：使用 ngrok 配置文件
```yaml
# ~/.config/ngrok/ngrok.yml
tunnels:
  frontend:
    addr: 3000
    proto: http
  backend:
    addr: 8080
    proto: http
```

启动：
```bash
ngrok start --all
```

## 访问地址

### 本地开发访问

- **前端**: http://localhost:3000
- **后端 API**: http://localhost:8080/api/v1
- **后端健康检查**: http://localhost:8080/health
- **Swagger 文档**: http://localhost:8080/swagger/index.html
- **ngrok 管理**: http://localhost:4040

### 局域网访问（公司内部）

- **前端**: http://192.168.31.151:3000
- **后端 API**: http://192.168.31.151:8080/api/v1

### 公网访问（任何地方）

- **前端**: https://sodless-judgmentally-kelsie.ngrok-free.dev
- **注意**: 免费版 URL 每次重启会变化

## 停止服务

### 停止所有服务

```bash
# 方法 1: 找到进程并停止
ps aux | grep -E "(ngrok|vite|go run)" | grep -v grep
kill <PID>

# 方法 2: 使用 pkill
pkill ngrok
pkill -f "go run cmd/server/main.go"
pkill -f "vite"
```

### 停止单个服务

如果服务在前台运行，在对应终端按 `Ctrl+C`。

如果服务在后台运行：
```bash
# 停止后端
pkill -f "go run cmd/server/main.go"

# 停止前端
pkill -f "vite"

# 停止 ngrok
pkill ngrok
```

## 常见问题

### 1. 端口已被占用

**错误**：`bind: address already in use`

**解决**：
```bash
# 查看占用端口的进程
lsof -i :3000
lsof -i :8080

# 停止占用的进程
kill <PID>
```

### 2. ngrok 显示 "ERR_NGROK_8012"

**原因**：ngrok 无法连接到本地端口

**解决**：
1. 确保前端服务正在运行（localhost:3000）
2. 检查端口是否正确
3. 重启前端服务

### 3. 前端访问后端 CORS 错误

**原因**：跨域配置问题

**解决**：后端已配置 CORS，允许所有来源：
```go
// internal/api/router/router.go
config := cors.DefaultConfig()
config.AllowOrigins = []string{"*"}
```

### 4. ngrok "Blocked request" 错误

**原因**：Vite 默认只允许 localhost 访问

**解决**：已在 `vite.config.js` 中配置：
```javascript
server: {
  host: '0.0.0.0',
  allowedHosts: [
    '.ngrok-free.dev',
    '.ngrok.io',
    'localhost'
  ]
}
```

### 5. 前端 API 请求失败

**检查**：
1. 后端是否正在运行：`curl http://localhost:8080/health`
2. 前端 .env 配置：`VITE_API_URL=/api/v1`
3. Vite proxy 配置：`frontend/vite.config.js`

### 6. 数据库连接失败

**错误**：`Failed to initialize database`

**解决**：
1. 检查 MySQL 是否运行：`brew services list | grep mysql`
2. 检查配置：`config.yaml`
3. 确认数据库已创建：
```bash
mysql -u root -protki123 -e "SHOW DATABASES;"
```

### 7. ngrok URL 变化导致前端无法访问

**原因**：免费版 ngrok 每次重启 URL 会变化

**解决**：
- 升级到 ngrok 付费版（固定域名）
- 或使用局域网 IP 访问（公司内部）
- 或每次重启后通知用户新的 URL

## 开发建议

### 1. 修改代码后的刷新

- **前端**：Vite 会自动热更新，无需重启
- **后端**：需要重启服务（Ctrl+C 后重新运行）
- **ngrok**：无需重启

### 2. 查看日志

```bash
# 前端日志（终端 2）
# 直接在运行 npm run dev 的终端查看

# 后端日志（终端 1）
# 直接在运行 go run 的终端查看

# ngrok 日志
# 访问 http://localhost:4040
# 或查看运行 ngrok 的终端
```

### 3. 调试建议

- 使用浏览器开发者工具（F12）查看网络请求
- 查看 ngrok 管理界面 (http://localhost:4040) 的请求详情
- 检查后端日志中的 API 请求和响应

## 环境配置

### 前端环境变量

文件：`frontend/.env`

```env
# API 配置 - 使用相对路径，利用 Vite proxy 代理到本地后端
VITE_API_URL=/api/v1
```

### 后端配置

文件：`config.yaml`

```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 3306
  username: root
  password: "rotki123"
  database: rotki_demo
```

## 项目结构

```
rotki-demo/
├── cmd/
│   └── server/
│       └── main.go           # 后端入口
├── frontend/
│   ├── src/
│   ├── vite.config.js        # Vite 配置
│   └── .env                  # 前端环境变量
├── config.yaml               # 后端配置
├── RUNNING_GUIDE.md          # 本文档
└── README.md                 # 项目说明
```

## 相关链接

- **ngrok 官网**: https://ngrok.com/
- **ngrok 文档**: https://ngrok.com/docs
- **ngrok Dashboard**: https://dashboard.ngrok.com/
- **Vite 文档**: https://vitejs.dev/
- **Vue.js 文档**: https://vuejs.org/
