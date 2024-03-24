package models

import (
	"time"
)

type Comment struct {
	ID        int64  `gorm:"primaryKey"`
	UserID    int64  `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID"`
	PhotoID   int64  `gorm:"not null"`
	Photo     Photo  `gorm:"foreignKey:PhotoID"`
	Message   string `gorm:"size:200;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateCommentRequest struct {
	PhotoID int64  `json:"photo_id" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type UpdateCommentRequest struct {
	Message string `json:"message" validate:"required"`
}
