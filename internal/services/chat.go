package services

import (
	"chat-app-backend/internal/database"
	"chat-app-backend/internal/models"
	websockets "chat-app-backend/internal/websocket"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

var WebSocketServer *websockets.WebSocketServer

func InitWebSocketServer() {
	WebSocketServer = websockets.NewWebSocketServer()
	go WebSocketServer.HandleMessages()
}

func SaveMessage(message *models.Message) error {
	if database.DB == nil {
		return errors.New("database connection is nil")
	}

	if message == nil {
		return errors.New("message is nil")
	}
	
	if err := database.DB.Create(message).Error; err != nil {
		return err
	}

	broadcastMessage(message)

	return nil
}

func broadcastMessage(message *models.Message) {
	msgJSON, _ := json.Marshal(message)
	WebSocketServer.Broadcast <- msgJSON
}

func GetAllConversations(userID uuid.UUID) ([]models.Message, error) {
    var messages []models.Message
    if err := database.DB.Preload("Sender").Preload("Receiver").
        Where("sender_id = ? OR receiver_id = ?", userID, userID).
        Find(&messages).Error; err != nil {
        return nil, err
    }
    return messages, nil
}

func GetConversationWithUser(userID uuid.UUID, otherUserID uuid.UUID) ([]models.Message, error) {
    var messages []models.Message
    if err := database.DB.Preload("Sender").Preload("Receiver").
        Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
            userID, otherUserID, otherUserID, userID).
        Find(&messages).Error; err != nil {
        return nil, err
    }
    return messages, nil
}

func DeleteMessage(userID uuid.UUID, messageID uuid.UUID) error {
    var message models.Message
    if err := database.DB.Where("id = ? AND sender_id = ?", messageID, userID).First(&message).Error; err != nil {
        return err
    }
    if err := database.DB.Delete(&message).Error; err != nil {
        return err
    }
    return nil
}

func UpdateMessage(userID uuid.UUID, messageID uuid.UUID, content string) error {
    var message models.Message
    if err := database.DB.Where("id = ? AND sender_id = ?", messageID, userID).First(&message).Error; err != nil {
        return err
    }
    message.Content = content
    if err := database.DB.Save(&message).Error; err != nil {
        return err
    }
    return nil
}

func GetMessageByID(id uuid.UUID) (*models.Message, error) {
	var message models.Message
	if err := database.DB.First(&message, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}