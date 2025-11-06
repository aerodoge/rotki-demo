-- Add RPC nodes table
CREATE TABLE IF NOT EXISTS `rpc_nodes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `chain_id` varchar(50) NOT NULL,
  `name` varchar(255) NOT NULL,
  `url` varchar(500) NOT NULL,
  `weight` int NOT NULL DEFAULT 100,
  `is_enabled` tinyint(1) NOT NULL DEFAULT 1,
  `is_connected` tinyint(1) NOT NULL DEFAULT 0,
  `last_checked` datetime(3) DEFAULT NULL,
  `priority` int NOT NULL DEFAULT 0,
  `timeout` int NOT NULL DEFAULT 30,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_rpc_nodes_chain_id` (`chain_id`),
  CONSTRAINT `fk_rpc_nodes_chain` FOREIGN KEY (`chain_id`) REFERENCES `chains` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Add indexes for better query performance
CREATE INDEX `idx_rpc_nodes_enabled` ON `rpc_nodes` (`is_enabled`);
CREATE INDEX `idx_rpc_nodes_chain_enabled` ON `rpc_nodes` (`chain_id`, `is_enabled`);
