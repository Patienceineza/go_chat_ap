package handlers

import (
	"chat-app-backend/internal/database"
	"chat-app-backend/internal/models"
	"chat-app-backend/internal/services"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Role == "" {
		input.Role = "user"
	}

	if input.Role == "admin" {
		claims, _ := c.Get("claims")
		userClaims := claims.(*services.Claims)
		if userClaims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can create another admin"})
			return
		}
	}

	user := models.User{
		ID:       uuid.New(),
		Username: input.Username,
		Password: input.Password,
		Role:     input.Role,
	}

	if err := services.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "userId": user.ID})
}

func LoginHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.LoginUser(input.Username, input.Password)
	fmt.Println(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := services.GenerateJWT(user.ID.String(), user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully", "token": token, "userId": user.ID})
}

func GetAllUsersHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*services.Claims)

	if userClaims.Role != "admin" {
		user, err := services.GetUserByUsername(userClaims.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user, "userId": userClaims.ID})
		return
	}

	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users, "userId": userClaims.ID})
}

func UpdateUserHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*services.Claims)

	id := c.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input struct {
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Role != "" && userClaims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can update roles"})
		return
	}

	err = services.UpdateUser(userClaims.Username, userId, input.Password, input.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "userId": userId})
}

func DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := services.DeleteUser(userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully", "userId": userId})
}

func GetUserFromTokenHandler(c *gin.Context) {
	
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims := claims.(*services.Claims)

	var user models.User
	if err := database.DB.First(&user, "id = ?", userClaims.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}