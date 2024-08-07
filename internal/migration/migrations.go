package migrations

import (
    "log"
    "gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
    
    if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
        log.Fatalf("Failed to create UUID extension: %v", err)
    }

    
    if err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id UUID DEFAULT uuid_generate_v4(),
            username VARCHAR(255) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL,
            role VARCHAR(50) NOT NULL DEFAULT 'user',
            PRIMARY KEY (id)
        )
    `).Error; err != nil {
        log.Fatalf("Failed to create users table: %v", err)
    }
}
