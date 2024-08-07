package handlers

import (
	"chat-app-backend/internal/models"
	"chat-app-backend/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func InitWebSocketRoutes(router *gin.Engine) {
	router.GET("/ws", func(c *gin.Context) {
		services.WebSocketServer.HandleConnections(c.Writer, c.Request)
	})
}

func CreateMessageHandler(c *gin.Context) {
	var input struct {
		ReceiverID uuid.UUID `json:"receiver_id" binding:"required"`
		Content    string    `json:"content" binding:"required"`
		ReplyToID  *uuid.UUID `json:"reply_to_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, _ := c.Get("claims")
	userClaims := claims.(*services.Claims)

	message := models.Message{
		SenderID:   userClaims.ID,
		ReceiverID: input.ReceiverID,
		Content:    input.Content,
		ReplyToID:  input.ReplyToID,
	}

	fmt.Println(message)
	if err := services.SaveMessage(&message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message created successfully"})
}

func GetAllConversationsHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*services.Claims)

	messages, err := services.GetAllConversations(userClaims.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func GetConversationWithUserHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*services.Claims)

	otherUserID := c.Param("user_id")
	user2UUID, err := uuid.Parse(otherUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	messages, err := services.GetConversationWithUser(userClaims.ID, user2UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func DeleteMessageHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*services.Claims)

	messageID := c.Param("id")
	messageUUid, err := uuid.Parse(messageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid messageId ID"})
		return
	}

	if err := services.DeleteMessage(userClaims.ID, messageUUid); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

func UpdateMessageHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*services.Claims)

	messageID := c.Param("id")
	messageUUid, err := uuid.Parse(messageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid messageId ID"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateMessage(userClaims.ID, messageUUid, input.Content); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully"})
}

func ReplyToMessageHandler(c *gin.Context) {
	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, _ := c.Get("claims")
	userClaims := claims.(*services.Claims)

	originalMessageID := c.Param("id")
	originalMessageUUID, err := uuid.Parse(originalMessageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	originalMessage, err := services.GetMessageByID(originalMessageUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Original message not found"})
		return
	}

	replyMessage := models.Message{
		SenderID:   userClaims.ID,
		ReceiverID: originalMessage.SenderID, 
		Content:    input.Content,
		ReplyToID:  &originalMessageUUID,
	}
	fmt.Println("original reply", replyMessage)

	if err := services.SaveMessage(&replyMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reply sent successfully"})
}
