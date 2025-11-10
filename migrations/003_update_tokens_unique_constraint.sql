-- 更新 tokens 表的唯一约束以包含 protocol_id
-- 这样可以允许相同的 token 在钱包和协议中都存在

-- 删除旧的唯一约束
ALTER TABLE tokens DROP INDEX uk_address_chain_token;

-- 添加新的唯一约束，包含 protocol_id
-- MySQL 允许多个 NULL 值，所以钱包代币（protocol_id=NULL）不会冲突
ALTER TABLE tokens ADD UNIQUE KEY uk_address_chain_token_protocol (address_id, chain_id, token_id, protocol_id);
