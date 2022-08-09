package controllers

import (
	"net/http"
	"rest-echo-gorm/helpers"
	"rest-echo-gorm/lib/databases"
	"rest-echo-gorm/models"
	"time"

	"github.com/labstack/echo/v4"
)

func CreateBookController(c echo.Context) error {
	req := models.Books{}
	c.Bind(&req)

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, &ResponseFormat{
			Status:   http.StatusBadRequest,
			Messages: "Fail",
			Data:     err.Error(),
		})
	}

	book, _ := databases.CreateBook(&req)

	var bookResponse models.BooksResponse
	bookResponse.ID = book.ID
	bookResponse.Tittle = book.Tittle
	bookResponse.Author = book.Author
	bookResponse.Year = book.Year

	return c.JSON(http.StatusCreated, &ResponseFormat{
		Status:   http.StatusCreated,
		Messages: "Success",
		Data:     &bookResponse,
	})
}

func GetBookController(c echo.Context) error {
	reqId := c.Param("id")

	book, err := databases.GetBook(reqId)

	if err != nil {
		return c.JSON(http.StatusNotFound, &ResponseFormat{
			Status:   http.StatusNotFound,
			Messages: "Fail",
			Data:     err.Error(),
		})
	}

	var bookResponse models.BooksResponse
	bookResponse.ID = book.ID
	bookResponse.Tittle = book.Tittle
	bookResponse.Author = book.Author
	bookResponse.Year = book.Year

	return c.JSON(http.StatusOK, &ResponseFormat{
		Status:   http.StatusOK,
		Messages: "Success",
		Data:     bookResponse,
	})
}

func GetBooksController(c echo.Context) error {
	books, err := databases.GetBooks()

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseFormat{
			Status:   http.StatusBadRequest,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	//var booksResponse []models.BooksResponse
	booksResponse := make([]models.BooksResponse, len(*books))

	for i, book := range *books {
		booksResponse[i].ID = book.ID
		booksResponse[i].Tittle = book.Tittle
		booksResponse[i].Author = book.Author
		booksResponse[i].Year = book.Year
	}

	return c.JSON(http.StatusOK, &ResponseFormat{
		Status:   http.StatusOK,
		Messages: "Success",
		Data:     booksResponse,
	})
}

func UpdateBookController(c echo.Context) error {
	reqId := c.Param("id")
	book, err := databases.GetBook(reqId)

	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseFormat{
			Status:   http.StatusNotFound,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	var bookPayload models.BooksResponse
	bookPayload.ID = helpers.ConvertStringToUint(reqId)
	c.Bind(&bookPayload)

	book.Tittle = bookPayload.Tittle
	book.Author = bookPayload.Author
	book.Year = bookPayload.Year
	book.UpdatedAt = time.Now()

	if err := c.Validate(&bookPayload); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseFormat{
			Status:   http.StatusBadRequest,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	if err := databases.UpdateBook(book); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseFormat{
			Status:   http.StatusBadRequest,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &ResponseFormat{
		Status:   http.StatusOK,
		Messages: "Success",
		Data:     bookPayload,
	})
}

func DeleteBookController(c echo.Context) error {
	reqId := c.Param("id")
	bookId, err := databases.DeleteBook(reqId)

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
		Data:     "Book with id " + bookId + " has been deleted",
	})
}
