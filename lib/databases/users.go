package databases

import (
	"log"
	"rest-echo-gorm/config"
	"rest-echo-gorm/helpers"
	"rest-echo-gorm/models"
)

var DB = config.ConnectDB()

func CreateUser(user models.Users) (models.Users, error) {
	err := user.BeforeSave()

	if err != nil {
		log.Fatal(err)
	}
	if err := DB.Create(&user).Error; err != nil {
		return models.Users{}, err
	}

	return user, nil
}

func GetUsers() ([]models.Users, error) {
	var users []models.Users

	if err := DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func GetUser(reqId string) (*models.Users, error) {
	var user models.Users
	if err := DB.Where("Id = ?", reqId).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(reqId string) (models.Users, error) {
	var user models.Users
	if err := DB.Where("Id = ?", reqId).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func DeleteUser(reqId string) (string, error) {
	var user models.Users

	if err := DB.Where("Id = ?", reqId).First(&user).Delete(&user).Error; err != nil {
		return "", err
	}

	userId := helpers.ConvertUintToString(user.ID)

	return userId, nil
}
