package seed

import (
    "chat-app-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "log"
)

func Seed(db *gorm.DB) {
    // Create dummy users
    users := []models.User{
        {
            ID:       uuid.New(),
            Username: "user1",
            Password: "password1", // Make sure to hash the password in a real application
            Role:     "user",
        },
        {
            ID:       uuid.New(),
            Username: "user2",
            Password: "password2", // Make sure to hash the password in a real application
            Role:     "user",
        },
        {
            ID:       uuid.New(),
            Username: "admin1",
            Password: "adminpass", // Make sure to hash the password in a real application
            Role:     "admin",
        },
    }

    for _, user := range users {
        if err := db.Create(&user).Error; err != nil {
            log.Printf("Failed to create user %s: %v", user.Username, err)
        }
    }
}
