package services

import (
	"github.com/heyjoakim/devops-21/models"
)

// GetUserID returns user ID for username
func GetUserID(username string) (uint, error) {
	var user models.User
	getUserIDErr := GetDBInstance().db.First(&user, "username = ?", username).Error
	if getUserIDErr != nil {
		LogError(models.Log{
			Message: getUserIDErr.Error(),
			Data:    username,
		})
	}
	return user.UserID, getUserIDErr
}

func GetUserFromUsername(username string) (models.User, error) {
	var user models.User
	err := GetDBInstance().db.Where("username = ?", username).First(&user).Error
	if err != nil {
		LogError(models.Log{
			Message: err.Error(),
			Data:    username,
		})
	}
	return user, err
}

func GetUser(userID uint) models.User {
	var user models.User
	getUserErr := GetDBInstance().db.First(&user, "user_id = ?", userID).Error
	if getUserErr != nil {
		LogError(models.Log{
			Message: getUserErr.Error(),
			Data:    userID,
		})
	}
	return user
}

// CreateUser creates a new user in the database
func CreateUser(user models.User) error {
	createUserErr := GetDBInstance().db.Create(&user).Error
	if createUserErr != nil {
		LogError(models.Log{
			Message: createUserErr.Error(),
			Data:    user,
		})
	}
	return createUserErr
}

// GetUserCount returns the number of users reigstered in the system
func GetUserCount() int64 {
	var count int64
	countErr := GetDBInstance().db.Find(&models.User{}).Count(&count).Error
	if countErr != nil {
		LogError(models.Log{
			Message: countErr.Error(),
		})
	}
	return count
}
