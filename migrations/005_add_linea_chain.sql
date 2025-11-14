-- Add Linea chain to the chains table
INSERT INTO chains (id, name, chain_type, logo_url, is_active)
VALUES ('linea', 'Linea', 'EVM', '/images/chains/linea.png', TRUE)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    chain_type = VALUES(chain_type),
    logo_url = VALUES(logo_url),
    is_active = VALUES(is_active),
    updated_at = CURRENT_TIMESTAMP;
