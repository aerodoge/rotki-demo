-- 添加 protocols 表

CREATE TABLE IF NOT EXISTS `protocols` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `address_id` bigint NOT NULL,
  `protocol_id` varchar(255) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `site_url` varchar(500) DEFAULT NULL,
  `logo_url` varchar(500) DEFAULT NULL,
  `chain_id` varchar(50) NOT NULL,
  `net_usd_value` decimal(30,6) DEFAULT NULL,
  `asset_usd_value` decimal(30,6) DEFAULT NULL,
  `debt_usd_value` decimal(30,6) DEFAULT NULL,
  `position_type` varchar(50) DEFAULT NULL,
  `raw_data` json DEFAULT NULL,
  `last_updated` timestamp DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_address_protocol` (`address_id`,`protocol_id`),
  KEY `idx_protocols_address_id` (`address_id`),
  KEY `idx_protocols_chain_id` (`chain_id`),
  CONSTRAINT `fk_protocols_address` FOREIGN KEY (`address_id`) REFERENCES `addresses` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
