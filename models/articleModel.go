package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	AuthorID    uint           `json:"author_id" gorm:"not null;index"`
	Author      Author         `json:"author" gorm:"foreignKey:AuthorID;constraint:onUpdate:CASCADE,onDelete:SET NULL"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Hidden in JSON responses
}
