package controllers

import (
	"net/http"
	"rest-echo-gorm/helpers"
	"rest-echo-gorm/lib/databases"
	"rest-echo-gorm/models"
	"time"

	"github.com/labstack/echo/v4"
)

type ResponseFormat struct {
	Status   int         `json:"status"`
	Messages string      `json:"messages"`
	Data     interface{} `json:"data"`
}

func CreateUserController(c echo.Context) (err error) {
	req := models.Users{}

	c.Bind(&req)

	//if err = c.Bind(&req); err != nil {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}

	if err = c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, _ := databases.CreateUser(req)

	// return userResponses from database
	var userResponses models.UsersResponse
	userResponses.ID = user.ID
	userResponses.Name = user.Name
	userResponses.Email = user.Email
	userResponses.Password = user.Password

	return c.JSON(http.StatusCreated, ResponseFormat{
		Status:   http.StatusCreated,
		Messages: "Success create user",
		Data:     userResponses,
	})
}

func GetUsersController(c echo.Context) error {
	users, err := databases.GetUsers()

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseFormat{
			Status:   http.StatusBadRequest,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	// make userResponses from database
	var userResponses []models.UsersResponse
	userResponses = make([]models.UsersResponse, len(users))

	for i, user := range users {
		userResponses[i].ID = user.ID
		userResponses[i].Name = user.Name
		userResponses[i].Email = user.Email
		userResponses[i].Password = user.Password
	}

	return c.JSON(http.StatusOK, ResponseFormat{
		Status:   http.StatusOK,
		Messages: "Success",
		Data:     userResponses,
	})
}

func GetUserController(c echo.Context) error {
	user, err := databases.GetUser(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseFormat{
			Status:   http.StatusNotFound,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	var userResponses models.UsersResponse
	userResponses.ID = user.ID
	userResponses.Name = user.Name
	userResponses.Email = user.Email
	userResponses.Password = user.Password

	return c.JSON(http.StatusOK, ResponseFormat{
		Status:   http.StatusOK,
		Messages: "Success",
		Data:     userResponses,
	})
}

func UpdateUserController(c echo.Context) error {
	reqId := c.Param("id")
	user, err := databases.UpdateUser(reqId)

	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseFormat{
			Status:   http.StatusNotFound,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	payload := models.UsersResponse{}
	payload.ID = helpers.ConvertStringToUint(reqId)
	c.Bind(&payload)

	// user.ID = helpers.ConvertStringToUint(reqId)
	user.ID = payload.ID
	user.Name = payload.Name
	user.Email = payload.Email
	user.Password = payload.Password
	user.UpdatedAt = time.Now()

	if err := c.Validate(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseFormat{
			Status:   http.StatusBadRequest,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	if err := databases.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFormat{
			Status:   http.StatusInternalServerError,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	return c.JSON(http.StatusOK, ResponseFormat{
		Status:   http.StatusOK,
		Messages: "Success",
		Data:     payload,
	})
}

func DeleteUserController(c echo.Context) error {
	reqId := c.Param("id")
	user, err := databases.DeleteUser(reqId)

	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseFormat{
			Status:   http.StatusNotFound,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	return c.JSON(http.StatusOK, ResponseFormat{
		Status:   http.StatusOK,
		Messages: "Success",
		Data:     "User with id " + user + " has been deleted",
	})
}
