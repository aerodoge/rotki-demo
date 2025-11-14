package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Get database connection from environment or use default
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/rotki_demo?parseTime=true"
	}

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Insert Linea chain
	query := `
		INSERT INTO chains (id, name, chain_type, logo_url, is_active)
		VALUES ('linea', 'Linea', 'EVM', '/images/chains/linea.png', TRUE)
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			chain_type = VALUES(chain_type),
			logo_url = VALUES(logo_url),
			is_active = VALUES(is_active),
			updated_at = CURRENT_TIMESTAMP
	`

	result, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to insert Linea chain:", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Successfully added/updated Linea chain. Rows affected: %d\n", rowsAffected)
}
