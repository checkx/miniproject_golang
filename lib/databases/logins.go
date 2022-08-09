package databases

import (
	"rest-echo-gorm/models"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(user models.Users) (*models.Users, error) {

	email := user.Email
	password := user.Password
	var err error

	user = models.Users{}

	err = DB.Debug().Model(models.Users{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, err
	}
	return &user, nil
}
