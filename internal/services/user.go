package services

import (
	"chat-app-backend/internal/database"
	"chat-app-backend/internal/models"
	
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

func RegisterUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.ID = uuid.New()
	if err := database.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// LoginUser takes a username and password as strings and returns a pointer to a User struct and an error.
func LoginUser(username, password string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(username string, userId uuid.UUID, password, role string) error {
	var user models.User
	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return err
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	if role != "" {
		user.Role = role
	}

	return database.DB.Save(&user).Error
}

func DeleteUser(userId uuid.UUID) error {
	return database.DB.Delete(&models.User{}, userId).Error
}
