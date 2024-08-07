package database

import (
	"chat-app-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.DBConnectionString,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = db

	// Create the uuid-ossp extension if it doesn't exist
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, err
	}

	// Create the users table if it doesn't exist
	if err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id UUID DEFAULT uuid_generate_v4(),
            username VARCHAR(255) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL,
            role VARCHAR(50) NOT NULL DEFAULT 'user',
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP WITH TIME ZONE,
            PRIMARY KEY (id)
        )
    `).Error; err != nil {
		return nil, err
	}

	// Create the messages table if it doesn't exist
	if err := DB.Exec(`
    CREATE TABLE IF NOT EXISTS messages (
        id UUID DEFAULT uuid_generate_v4(),
        sender_id UUID NOT NULL,
        receiver_id UUID NOT NULL,
        content TEXT NOT NULL,
        reply_to_id UUID,
        timestamp BIGINT DEFAULT EXTRACT(EPOCH FROM now()),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP WITH TIME ZONE,
        PRIMARY KEY (id),
        CONSTRAINT fk_sender
            FOREIGN KEY(sender_id) 
            REFERENCES users(id),
        CONSTRAINT fk_receiver
            FOREIGN KEY(receiver_id) 
            REFERENCES users(id)
    )
`).Error; err != nil {
		return nil, err
	}

	return DB, nil
}
