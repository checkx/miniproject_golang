package models

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Name     string `json:"name" form:"name" validate:"required,omitempty"`
	Email    string `json:"email" form:"email"  validate:"required"`
	Password string `json:"password" form:"password"  validate:"required"`
}

type UsersResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}
