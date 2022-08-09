package middleware

import (
	"fmt"
	"rest-echo-gorm/constants"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(id, username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["username"] = username
	claims["expired"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	fmt.Println("Token adalah : ", token)
	fmt.Println("Claims adalah : ", claims)

	return token.SignedString([]byte(constants.SecretKey))
}
