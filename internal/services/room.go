package services

import (
	"chat-app-backend/internal/database"
	"chat-app-backend/internal/models"
	"errors"

	"github.com/google/uuid"
)

func CreateRoom(room *models.ChatRoom) error {
	if database.DB == nil {
		return errors.New("database connection is nil")
	}
	if err := database.DB.Create(room).Error; err != nil {
		return err
	}
	return nil
}

func GetRooms() ([]models.ChatRoom, error) {
	var rooms []models.ChatRoom
	if err := database.DB.Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func JoinRoom(userID uuid.UUID, roomID uuid.UUID) error {
	// Implement logic to add user to room
	// You might need a new model to represent the user-room relationship
	// For now, let's assume you have such a table in your database
	var room models.ChatRoom
	if err := database.DB.First(&room, "id = ?", roomID).Error; err != nil {
		return err
	}

	// Add user to room logic
	// Example: database.DB.Create(&models.UserRoom{UserID: userID, RoomID: roomID})

	return nil
}

func LeaveRoom(userID uuid.UUID, roomID uuid.UUID) error {
	// Implement logic to remove user from room
	// Example: database.DB.Where("user_id = ? AND room_id = ?", userID, roomID).Delete(&models.UserRoom{})

	return nil
}
