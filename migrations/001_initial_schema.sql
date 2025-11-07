-- Rotki Demo 数据库模式

-- Wallets 表：存储钱包信息
CREATE TABLE wallets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    tags JSON DEFAULT NULL COMMENT '用户自定义标签',
    enabled_chains JSON DEFAULT NULL COMMENT '启用的链 ID 列表，NULL 表示所有链',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Addresses 表：存储与钱包关联的区块链地址
CREATE TABLE addresses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    wallet_id BIGINT NOT NULL,
    address VARCHAR(255) NOT NULL,
    chain_type VARCHAR(50) NOT NULL DEFAULT 'EVM', -- EVM, Bitcoin, Solana 等
    label VARCHAR(255),
    tags JSON DEFAULT NULL COMMENT '用户自定义标签',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_synced_at TIMESTAMP NULL,
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE,
    UNIQUE KEY uk_address_chain (address, chain_type),
    INDEX idx_wallet_id (wallet_id),
    INDEX idx_address (address)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Asset snapshots 表：存储每个地址的定期资产快照
CREATE TABLE asset_snapshots (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    address_id BIGINT NOT NULL,
    snapshot_time TIMESTAMP NOT NULL,
    total_usd_value DECIMAL(30, 6) DEFAULT 0,
    data_source VARCHAR(50) NOT NULL DEFAULT 'debank', -- debank、self-query 等
    raw_data JSON, -- 存储完整响应以供将来参考
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE,
    INDEX idx_address_time (address_id, snapshot_time),
    INDEX idx_snapshot_time (snapshot_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Chains 表：存储支持的区块链信息
CREATE TABLE chains (
    id VARCHAR(50) PRIMARY KEY, -- 例如：'eth'、'bsc'、'polygon'
    name VARCHAR(100) NOT NULL,
    chain_type VARCHAR(50) NOT NULL DEFAULT 'EVM',
    logo_url VARCHAR(512),
    native_token_id VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_chain_type (chain_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tokens 表：存储代币信息以供快速查找
CREATE TABLE tokens (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    address_id BIGINT NOT NULL,
    chain_id VARCHAR(50) NOT NULL,
    token_id VARCHAR(255) NOT NULL, -- 合约地址或原生代币 ID
    symbol VARCHAR(50),
    name VARCHAR(255),
    decimals INT,
    logo_url VARCHAR(512),
    balance VARCHAR(255), -- 存储为字符串以处理来自诈骗代币的超大值
    price DECIMAL(30, 6),
    usd_value DECIMAL(30, 6),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE,
    FOREIGN KEY (chain_id) REFERENCES chains(id),
    UNIQUE KEY uk_address_chain_token (address_id, chain_id, token_id),
    INDEX idx_address_id (address_id),
    INDEX idx_chain_id (chain_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- API rate limiting 表：API 速率限制跟踪
CREATE TABLE api_rate_limits (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    provider VARCHAR(50) NOT NULL, -- 'debank' 等
    endpoint VARCHAR(255) NOT NULL,
    request_count INT DEFAULT 0,
    window_start TIMESTAMP NOT NULL,
    window_end TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_provider_window (provider, window_start, window_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Sync jobs 表：跟踪后台同步操作
CREATE TABLE sync_jobs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    address_id BIGINT,
    wallet_id BIGINT,
    job_type VARCHAR(50) NOT NULL, -- 'full_sync'、'token_sync'、'protocol_sync'
    status VARCHAR(50) NOT NULL, -- 'pending'、'running'、'completed'、'failed'
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE,
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE,
    INDEX idx_status (status),
    INDEX idx_address_id (address_id),
    INDEX idx_wallet_id (wallet_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- RPC nodes 表：存储每条链的 RPC 端点
CREATE TABLE rpc_nodes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    chain_id VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    endpoint VARCHAR(512) NOT NULL,
    weight INT DEFAULT 100 COMMENT '负载均衡权重 (0-100)',
    is_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (chain_id) REFERENCES chains(id) ON DELETE CASCADE,
    INDEX idx_chain_id (chain_id),
    INDEX idx_enabled (is_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
