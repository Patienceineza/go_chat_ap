package handlers

import (
	"chat-app-backend/internal/models"
	"chat-app-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateRoomHandler handles the creation of a new chat room
func CreateRoomHandler(c *gin.Context) {
	var room models.ChatRoom
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room.ID = uuid.New()
	if err := services.CreateRoom(&room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

// GetRoomsHandler handles fetching all chat rooms
func GetRoomsHandler(c *gin.Context) {
	rooms, err := services.GetRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

// JoinRoomHandler handles joining a chat room
func JoinRoomHandler(c *gin.Context) {
	roomID := c.Param("room_id")
	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	claims, _ := c.Get("claims")
	userClaims := claims.(*services.Claims)

	if err := services.JoinRoom(userClaims.ID, roomUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined room successfully"})
}

// LeaveRoomHandler handles leaving a chat room
func LeaveRoomHandler(c *gin.Context) {
	roomID := c.Param("room_id")
	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	claims, _ := c.Get("claims")
	userClaims := claims.(*services.Claims)

	if err := services.LeaveRoom(userClaims.ID, roomUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left room successfully"})
}
