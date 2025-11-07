# Rotki Demo API 文档

本项目使用 Swagger 自动生成 API 文档。

## 访问 Swagger UI

启动服务器后，访问以下地址查看交互式 API 文档：

```
http://localhost:8080/swagger/index.html
```

## 生成文档

如果修改了 API 注释，需要重新生成文档：

```bash
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

## 文档位置

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON 文档**: http://localhost:8080/swagger/doc.json
- **YAML 文档**: http://localhost:8080/swagger/swagger.yaml

## API 端点概览

### 钱包管理 (Wallets)
- `POST /api/v1/wallets` - 创建钱包
- `GET /api/v1/wallets` - 获取钱包列表
- `GET /api/v1/wallets/{id}` - 获取钱包详情
- `PUT /api/v1/wallets/{id}` - 更新钱包
- `DELETE /api/v1/wallets/{id}` - 删除钱包
- `POST /api/v1/wallets/{id}/refresh` - 刷新钱包数据

### 地址管理 (Addresses)
- `POST /api/v1/addresses` - 创建地址
- `GET /api/v1/addresses` - 获取地址列表
- `GET /api/v1/addresses/{id}` - 获取地址详情
- `PUT /api/v1/addresses/{id}` - 更新地址
- `DELETE /api/v1/addresses/{id}` - 删除地址
- `POST /api/v1/addresses/{id}/refresh` - 刷新地址资产

### 链信息 (Chains)
- `GET /api/v1/chains` - 获取所有支持的区块链列表

### RPC 节点管理 (RPC Nodes)
- `POST /api/v1/rpc-nodes` - 创建 RPC 节点
- `GET /api/v1/rpc-nodes` - 获取 RPC 节点列表
- `GET /api/v1/rpc-nodes/grouped` - 按链分组获取 RPC 节点
- `GET /api/v1/rpc-nodes/{id}` - 获取 RPC 节点详情
- `PUT /api/v1/rpc-nodes/{id}` - 更新 RPC 节点
- `DELETE /api/v1/rpc-nodes/{id}` - 删除 RPC 节点
- `POST /api/v1/rpc-nodes/{id}/check` - 检查单个 RPC 节点连接
- `POST /api/v1/rpc-nodes/check-all` - 检查所有 RPC 节点连接

## 添加 API 注释

在 handler 函数上方添加 Swagger 注释，例如：

```go
// CreateWallet 创建一个新的钱包
// @Summary      创建钱包
// @Description  创建一个新的钱包
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Param        wallet  body      CreateWalletRequest  true  "钱包信息"
// @Success      201     {object}  github_com_rotki-demo_internal_models.Wallet
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /wallets [post]
func (h *WalletHandler) CreateWallet(c *gin.Context) {
    // ...
}
```

## 工具依赖

- [swaggo/swag](https://github.com/swaggo/swag) - Swagger 文档生成工具
- [gin-swagger](https://github.com/swaggo/gin-swagger) - Gin 框架 Swagger 中间件

## 相关链接

- [Swagger 官方文档](https://swagger.io/)
- [Swaggo 文档](https://github.com/swaggo/swag#getting-started)
- [OpenAPI 规范](https://spec.openapis.org/oas/v3.0.0)
