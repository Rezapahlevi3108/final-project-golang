package controllers

import (
	"net/http"
	"time"

	"final-project-golang/models"
	"final-project-golang/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SocialMediaController struct {
	DB *gorm.DB
}

func NewSocialMediaController(DB *gorm.DB) SocialMediaController {
	return SocialMediaController{DB}
}

func (smc *SocialMediaController) CreateSocialMedia(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload models.CreateSocialMediaRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if payload.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Name is required"})
		return
	}

	if payload.SocialMediaURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Social media URL is required"})
		return
	}

	if payload.SocialMediaURL != "" && !utils.IsValidURL(payload.SocialMediaURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid profile image URL format"})
		return
	}

	now := time.Now()
	newSocialMedia := models.SocialMedia{
		Name:           payload.Name,
		SocialMediaURL: payload.SocialMediaURL,
		UserID:         currentUser.ID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	result := smc.DB.Create(&newSocialMedia)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":               newSocialMedia.ID,
		"name":             newSocialMedia.Name,
		"social_media_url": newSocialMedia.SocialMediaURL,
		"user_id":          newSocialMedia.UserID,
	})
}
func (smc *SocialMediaController) GetSocialMedias(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	var socialMedias []models.SocialMedia
	result := smc.DB.Preload("User").Where("user_id = ?", currentUser.ID).Find(&socialMedias)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error"})
		return
	}

	var responseData []gin.H
	for _, socialMedia := range socialMedias {
		responseData = append(responseData, gin.H{
			"id":               socialMedia.ID,
			"name":             socialMedia.Name,
			"social_media_url": socialMedia.SocialMediaURL,
			"user_id":          socialMedia.UserID,
			"user": gin.H{
				"id":       socialMedia.User.ID,
				"email":    socialMedia.User.Email,
				"username": socialMedia.User.Username,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

func (smc *SocialMediaController) GetSocialMediaByID(ctx *gin.Context) {
	socialMediaID := ctx.Param("socialMediaId")

	var socialMedia models.SocialMedia
	result := smc.DB.Preload("User").First(&socialMedia, "id = ?", socialMediaID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No social media entry with that ID exists"})
		return
	}
	user := models.User{}
	smc.DB.First(&user, socialMedia.UserID)

	responseData := gin.H{
		"id":               socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user_id":          socialMedia.UserID,
		"user": gin.H{
			"id":       socialMedia.User.ID,
			"email":    socialMedia.User.Email,
			"username": socialMedia.User.Username,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

func (smc *SocialMediaController) UpdateSocialMedia(ctx *gin.Context) {
	socialMediaID := ctx.Param("socialMediaId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload models.UpdateSocialMediaRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if payload.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Name is required"})
		return
	}

	if payload.SocialMediaURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Social media URL is required"})
		return
	}

	if !utils.IsValidURL(payload.SocialMediaURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid profile image URL format"})
		return
	}

	var updatedSocialMedia models.SocialMedia
	result := smc.DB.First(&updatedSocialMedia, "id = ?", socialMediaID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No social media with that ID exists"})
		return
	}

	if updatedSocialMedia.UserID != currentUser.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to update this social media entry"})
		return
	}

	updatedSocialMedia.Name = payload.Name
	updatedSocialMedia.SocialMediaURL = payload.SocialMediaURL
	smc.DB.Save(&updatedSocialMedia)

	responseData := gin.H{
		"id":               updatedSocialMedia.ID,
		"name":             updatedSocialMedia.Name,
		"social_media_url": updatedSocialMedia.SocialMediaURL,
		"user_id":          updatedSocialMedia.UserID,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

func (smc *SocialMediaController) DeleteSocialMedia(ctx *gin.Context) {
	socialMediaID := ctx.Param("socialMediaId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var socialMedia models.SocialMedia
	result := smc.DB.First(&socialMedia, "id = ?", socialMediaID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No social media with that ID exists"})
		return
	}

	if socialMedia.UserID != currentUser.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to delete this social media entry"})
		return
	}

	result = smc.DB.Delete(&socialMedia)
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No social media with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
