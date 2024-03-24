package models

import (
	"time"
)

type Photo struct {
	ID        int64  `gorm:"primaryKey"`
	Title     string `gorm:"size:100;not null"`
	Caption   string `gorm:"size:200"`
	PhotoURL  string `gorm:"type:text;not null"`
	UserID    int64  `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Comments  []Comment
}

type CreatePhotoRequest struct {
	Title    string `json:"title" binding:"required"`
	Caption  string `json:"caption,omitempty"`
	PhotoURL string `json:"photo_url" binding:"required"`
}

type UpdatePhoto struct {
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption,omitempty"`
	PhotoURL string `json:"photo_url" validate:"required"`
}
