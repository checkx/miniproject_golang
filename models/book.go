package models

import "gorm.io/gorm"

type Books struct {
	gorm.Model
	Tittle string `json:"tittle" form:"tittle" validate:"required"`
	Author string `json:"author" form:"author" validate:"required"`
	Year   int    `json:"year" form:"year" validate:"required"`
}

type BooksResponse struct {
	ID     uint   `json:"id"`
	Tittle string `json:"tittle" form:"tittle" validate:"required"`
	Author string `json:"author" form:"author" validate:"required"`
	Year   int    `json:"year" form:"year" validate:"required"`
}
