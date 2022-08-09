package routes

import (
	"rest-echo-gorm/constants"
	"rest-echo-gorm/controllers"
	"rest-echo-gorm/helpers"
	m "rest-echo-gorm/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	// Initialaze validator
	e.Validator = &helpers.CustomValidator{Validator: validator.New()}

	// Middleware
	m.LogMiddleware(e)

	// Login user
	e.POST("/auth/login", controllers.LoginUserController)

	// Register user
	e.POST("/users", controllers.CreateUserController)

	// Book
	e.GET("/books", controllers.GetBooksController)
	e.GET("/books/:id", controllers.GetBookController)

	// JWT Grouping
	jwtAuth := e.Group("")
	jwtAuth.Use(middleware.JWT([]byte(constants.SecretKey)))

	jwtAuth.GET("/users", controllers.GetUsersController)
	jwtAuth.GET("/users/:id", controllers.GetUserController)
	jwtAuth.PUT("/users/:id", controllers.UpdateUserController)
	jwtAuth.DELETE("/users/:id", controllers.DeleteUserController)

	// Books routes
	jwtAuth.POST("/books", controllers.CreateBookController)
	jwtAuth.PUT("/books/:id", controllers.UpdateBookController)
	jwtAuth.DELETE("/books/:id", controllers.DeleteBookController)

	return e
}
