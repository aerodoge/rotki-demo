# Database Migrations

This directory contains SQL migration files for the database schema.

## Files

- `001_initial_schema.sql` - Initial database schema with all tables

## Usage

### With Docker
Migrations are automatically applied when the MySQL container starts for the first time.

### Manual Application
```bash
# Using MySQL command line
mysql -u root -p rotki_demo < migrations/001_initial_schema.sql

# Using Docker
docker exec -i rotki-mysql mysql -uroot -protki123 rotki_demo < migrations/001_initial_schema.sql
```

## Schema Overview

### Core Tables
- **wallets** - Wallet metadata
- **addresses** - Blockchain addresses linked to wallets
- **tokens** - Token balances for each address
- **chains** - Blockchain information
- **asset_snapshots** - Historical balance snapshots
- **sync_jobs** - Background sync job tracking

### Relationships
```
wallets (1:N) -> addresses (1:N) -> tokens
                            |
                            +-----> asset_snapshots
```

## Adding New Migrations

When adding new migrations:

1. Create a new file: `002_description.sql`
2. Use proper naming: `NNN_description.sql`
3. Include both UP and DOWN migrations
4. Test thoroughly before committing

Example:
```sql
-- 002_add_nft_table.sql
CREATE TABLE nfts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    address_id BIGINT NOT NULL,
    ...
);
```

## Notes

- Migrations are run in alphabetical order
- The application also uses GORM AutoMigrate as a fallback
- For production, consider using a proper migration tool like golang-migrate
