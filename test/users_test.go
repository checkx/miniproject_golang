package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"rest-echo-gorm/config"
	"rest-echo-gorm/controllers"
	"rest-echo-gorm/helpers"
	"rest-echo-gorm/models"
	"strconv"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func CleanTable(table []string) {
	for _, tableName := range table {
		config.ConnectDB().Exec("DELETE FROM " + tableName)
	}
}

func insertUserDb(name, email, password string) models.Users {
	user := models.Users{
		Name:     name,
		Email:    email,
		Password: password,
	}

	if err := config.ConnectDB().Create(&user).Error; err != nil {
		panic(err)
	}

	return user
}

func CreateToken() string {
	insertUserDb("user1", "user1@gmail.com", "user123")

	requestBody := strings.NewReader(`{"email":"user1@gmail.com","password":"user123"}`)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", requestBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/auth/login")
	controllers.LoginUserController(c)

	e.Validator = &helpers.CustomValidator{Validator: validator.New()}

	// Decode rec.body and exctract Data.id
	var response struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(rec.Body.String()), &response); err != nil {
		panic(err)
	}

	return response.Data.Token
}

func TestCreateUserController(t *testing.T) {
	t.Run("Test create user with valid payload", func(t *testing.T) {
		config.InitialMigration()

		fmt.Println("get token", CreateToken())

		requestBody := strings.NewReader(`{"name":"user1","email":"user1@gmail.com","password":"user123"}`)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/users", requestBody)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		if assert.NoError(t, controllers.CreateUserController(c)) {
			// Decode rec.body and exctract Data.id
			var response struct {
				Data models.UsersResponse `json:"data"`
			}

			if err := json.Unmarshal([]byte(rec.Body.String()), &response); err != nil {
				t.Fatal(err)
			}

			expectedReturns := &controllers.ResponseFormat{
				Status:   http.StatusCreated,
				Messages: "Success create user",
				Data: models.UsersResponse{
					ID:       response.Data.ID,
					Name:     "user1",
					Email:    "user1@gmail.com",
					Password: "user123",
				},
			}

			var expectedReturnsJson bytes.Buffer
			if err := json.NewEncoder(&expectedReturnsJson).Encode(expectedReturns); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, expectedReturnsJson.String(), rec.Body.String())
		}

		// Clean table
		CleanTable([]string{"users"})
	})

	t.Run("Test create user with bad payload", func(t *testing.T) {
		requestBody := strings.NewReader(`{"name":"user1","email":"user1","password":"user123"}`)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/users", requestBody)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		if assert.NoError(t, controllers.CreateUserController(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}

		// Clean table
		CleanTable([]string{"users"})
	})
}

func TestGetUsersController(t *testing.T) {
	t.Run("Test get users", func(t *testing.T) {
		insertUsers := []models.Users{
			{
				Name:     "user1",
				Email:    "user1@gmail.com",
				Password: "user123",
			},
			{
				Name:     "user2",
				Email:    "user2@gmail.com",
				Password: "user123",
			},
			{
				Name:     "user3",
				Email:    "user3@gmail.com",
				Password: "user123",
			},
		}

		for _, user := range insertUsers {
			insertUserDb(user.Name, user.Email, user.Password)
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, controllers.GetUsersController(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			body := rec.Body.String()

			var response struct {
				Users []models.Users `json:"data"`
			}

			err := json.Unmarshal([]byte(body), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, len(insertUsers), len(response.Users))
		}

		// Clean table
		CleanTable([]string{"users"})
	})
}

func TestGetUserController(t *testing.T) {
	t.Run("Test get user with valid id", func(t *testing.T) {
		user := insertUserDb("user1", "user1@gmail.com", "user123")

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(user.ID)))

		expectedResponse := &controllers.ResponseFormat{
			Status:   http.StatusOK,
			Messages: "Success",
			Data: models.UsersResponse{
				ID:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Password: user.Password,
			},
		}

		var expectedResponseJson bytes.Buffer
		if err := json.NewEncoder(&expectedResponseJson).Encode(expectedResponse); err != nil {
			t.Fatal(err)
		}

		if assert.NoError(t, controllers.GetUserController(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expectedResponseJson.String(), rec.Body.String())
		}

		// Clean table
		CleanTable([]string{"users"})
	})

	t.Run("Test get user with invalid id", func(t *testing.T) {
		insertUserDb("user1", "user1@gmail.com", "user123")

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(1))

		expectedResponse := &controllers.ResponseFormat{
			Status:   http.StatusNotFound,
			Messages: "Failed",
			Data:     "record not found",
		}

		var expectedResponseJson bytes.Buffer
		if err := json.NewEncoder(&expectedResponseJson).Encode(expectedResponse); err != nil {
			t.Fatal(err)
		}

		if assert.NoError(t, controllers.GetUserController(c)) {
			assert.Equal(t, expectedResponse.Status, rec.Code)
			assert.Equal(t, expectedResponseJson.String(), rec.Body.String())
		}

		// Clean table
		CleanTable([]string{"users"})
	})
}

func TestUpdateUserController(t *testing.T) {
	t.Run("Update user with valid payload", func(t *testing.T) {
		// First insert user into DB
		user := insertUserDb("user1", "user1@gmail.com", "user123")

		// Create requestBody for update user
		requestBody := map[string]interface{}{
			"name":     "user2Changed",
			"email":    "user2change@gmail.com",
			"password": "user2",
		}
		// Convert requestBody to json
		requestBodyJson, _ := json.Marshal(requestBody)

		// Initial echo
		e := echo.New()
		// Create request
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(requestBodyJson))
		req.Header.Set("Content-Type", "application/json")
		// Create recorder
		rec := httptest.NewRecorder()
		// Create context
		c := e.NewContext(req, rec)
		// Set user id
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(helpers.ConvertUintToString(user.ID))

		// Validate request
		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		// Create response expected
		expectedReturns := &controllers.ResponseFormat{
			Status:   http.StatusOK,
			Messages: "Success",
			Data: models.UsersResponse{
				ID:       user.ID,
				Name:     "user2Changed",
				Email:    "user2change@gmail.com",
				Password: "user2",
			},
		}

		// Check if no error on controller
		if assert.NoError(t, controllers.UpdateUserController(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			// Convert struct ResponseFormat to json
			var expectedReturnsJson bytes.Buffer
			if err := json.NewEncoder(&expectedReturnsJson).Encode(expectedReturns); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expectedReturnsJson.String(), rec.Body.String())

		}

		// Clean table
		CleanTable([]string{"users"})
	})

	t.Run("Update user with invalid payload", func(t *testing.T) {
		// First insert user into DB
		user := insertUserDb("user1", "user1@gmail.com", "user123")

		// Create requestBody for update user
		requestBody := map[string]interface{}{
			"name":     "user2Changed",
			"email":    "user2changed",
			"password": "user2",
		}
		// Convert requestBody to json
		requestBodyJson, _ := json.Marshal(requestBody)

		// Initial echo
		e := echo.New()
		// Create request
		req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewReader(requestBodyJson))
		req.Header.Set("Content-Type", "application/json")
		// Create recorder
		rec := httptest.NewRecorder()
		// Create context
		c := e.NewContext(req, rec)
		// Set user id
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(helpers.ConvertUintToString(user.ID))

		// Validate request
		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		// Create response expected
		expectedReturns := &controllers.ResponseFormat{
			Status:   http.StatusBadRequest,
			Messages: "Failed",
			Data:     "code=400, message=Key: 'UsersResponse.Email' Error:Field validation for 'Email' failed on the 'email' tag",
		}

		// Check if no error on controller
		if assert.NoError(t, controllers.UpdateUserController(c)) {
			assert.Equal(t, expectedReturns.Status, rec.Code)

			// Convert struct ResponseFormat to json
			var expectedReturnsJson bytes.Buffer
			if err := json.NewEncoder(&expectedReturnsJson).Encode(expectedReturns); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expectedReturnsJson.String(), rec.Body.String())

		}

		// Clean table
		CleanTable([]string{"users"})
	})
}

func TestDeleteUserController(t *testing.T) {
	t.Run("Delete user with valid id", func(t *testing.T) {
		user := insertUserDb("user1", "user1@gmail.com", "user123")

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(user.ID)))

		expectedResponse := &controllers.ResponseFormat{
			Status:   http.StatusOK,
			Messages: "Success",
			Data:     "User with id " + strconv.Itoa(int(user.ID)) + " has been deleted",
		}

		var expectedResponseJson bytes.Buffer
		if err := json.NewEncoder(&expectedResponseJson).Encode(expectedResponse); err != nil {
			t.Fatal(err)
		}

		if assert.NoError(t, controllers.DeleteUserController(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expectedResponseJson.String(), rec.Body.String())
		}
	})

	t.Run("Delete user with invalid id", func(t *testing.T) {
		insertUserDb("user1", "user1@gmail.com", "user123")

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(1))

		expectedResponse := &controllers.ResponseFormat{
			Status:   http.StatusNotFound,
			Messages: "Failed",
			Data:     "record not found",
		}

		var expectedResponseJson bytes.Buffer
		if err := json.NewEncoder(&expectedResponseJson).Encode(expectedResponse); err != nil {
			t.Fatal(err)
		}

		if assert.NoError(t, controllers.DeleteUserController(c)) {
			assert.Equal(t, expectedResponse.Status, rec.Code)
			assert.Equal(t, expectedResponseJson.String(), rec.Body.String())
		}
	})
}
