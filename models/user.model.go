package models

import (
	"time"
)

type User struct {
	ID              int64  `gorm:"primaryKey"`
	Username        string `gorm:"size:50;not null"`
	Email           string `gorm:"size:150;not null"`
	Password        string `gorm:"type:text;not null"`
	Age             int    `gorm:"not null"`
	ProfileImageURL string `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Photos          []Photo       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Comments        []Comment     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	SocialMedias    []SocialMedia `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type SignUpInput struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	Age             int    `json:"age" binding:"required"`
	ProfileImageURL string `json:"profile_image_url" validate:"omitempty"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UpdateCurrentUserRequest struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Age             int    `json:"age" binding:"required"`
	ProfileImageURL string `json:"profile_image_url,omitempty" validate:"omitempty"`
}
