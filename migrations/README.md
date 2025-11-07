# 数据库迁移

此目录包含数据库模式的 SQL 迁移文件。

## 文件

- `001_initial_schema.sql` - 包含所有表的初始数据库模式

## 使用方法

### 使用 Docker
当 MySQL 容器首次启动时，迁移会自动应用。

### 手动应用
```bash
# 使用 MySQL 命令行
mysql -u root -p rotki_demo < migrations/001_initial_schema.sql

# 使用 Docker
docker exec -i rotki-mysql mysql -uroot -protki123 rotki_demo < migrations/001_initial_schema.sql
```

## 模式概览

### 核心表
- **wallets** - 钱包元数据
- **addresses** - 链接到钱包的区块链地址
- **tokens** - 每个地址的代币余额
- **chains** - 区块链信息
- **asset_snapshots** - 历史余额快照
- **sync_jobs** - 后台同步作业跟踪

### 关系
```
wallets (1:N) -> addresses (1:N) -> tokens
                            |
                            +-----> asset_snapshots
```

## 添加新迁移

添加新迁移时：

1. 创建新文件：`002_description.sql`
2. 使用正确的命名：`NNN_description.sql`
3. 包括 UP 和 DOWN 迁移
4. 提交前彻底测试

示例：
```sql
-- 002_add_nft_table.sql
CREATE TABLE nfts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    address_id BIGINT NOT NULL,
    ...
);
```

## 注意事项

- 迁移按字母顺序运行
- 应用程序还使用 GORM AutoMigrate 作为后备
- 对于生产环境，考虑使用适当的迁移工具，如 golang-migrate
