-- Add status column to wallets table
ALTER TABLE wallets ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'Enabled';
