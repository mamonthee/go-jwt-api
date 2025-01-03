package models

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserName  string         `json:"user_name" validate:"required,min=2,max=100"`
	Email     string         `json:"email" gorm:"unique" validate:"required,email"`
	Password  string         `json:"password" validate:"required,min=8"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
