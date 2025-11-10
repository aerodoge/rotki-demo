package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:rotki123@tcp(127.0.0.1:3306)/rotki_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 删除旧约束
	result := db.Exec("ALTER TABLE tokens DROP INDEX uk_address_chain_token")
	if result.Error != nil {
		log.Fatal("Failed to drop old index:", result.Error)
	}
	fmt.Println("Dropped old index successfully")

	// 添加新约束
	result = db.Exec("ALTER TABLE tokens ADD UNIQUE KEY uk_address_chain_token_protocol (address_id, chain_id, token_id, protocol_id)")
	if result.Error != nil {
		log.Fatal("Failed to add new index:", result.Error)
	}

	fmt.Println("Migration completed successfully!")
	os.Exit(0)
}
