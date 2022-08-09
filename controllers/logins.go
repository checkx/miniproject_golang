package controllers

import (
	"fmt"
	"net/http"
	"rest-echo-gorm/lib/databases"
	"rest-echo-gorm/middleware"
	"rest-echo-gorm/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

func LoginUserController(c echo.Context) error {
	req := models.Users{}
	c.Bind(&req)
	user, err := databases.LoginUser(req)

	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseFormat{
			Status:   http.StatusNotFound,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	token, err := middleware.CreateToken(strconv.Itoa(int(user.ID)), req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ResponseFormat{
			Status:   http.StatusInternalServerError,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	fmt.Println("return user LoginUserController : ", user)
	fmt.Println("return token LoginUserController : ", token)

	return c.JSON(http.StatusOK, &ResponseFormat{
		Status:   http.StatusOK,
		Messages: "Success",
		Data: &models.Logins{
			Name:  user.Name,
			Email: user.Email,
			Token: token,
		},
	})
}
