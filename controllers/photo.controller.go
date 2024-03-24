package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"final-project-golang/models"
	"final-project-golang/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PhotoController struct {
	DB *gorm.DB
}

func NewPhotoController(DB *gorm.DB) PhotoController {
	return PhotoController{DB}
}

func (pc *PhotoController) CreatePhoto(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreatePhotoRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if payload.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Title is required"})
		return
	}

	if payload.PhotoURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Photo URL is required"})
		return
	}

	if payload.PhotoURL != "" && !utils.IsValidURL(payload.PhotoURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid profile image URL format"})
		return
	}

	now := time.Now()
	newPhoto := models.Photo{
		Title:     payload.Title,
		Caption:   payload.Caption,
		PhotoURL:  payload.PhotoURL,
		UserID:    currentUser.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := pc.DB.Create(&newPhoto)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"message": "Photo with that title already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":        newPhoto.ID,
		"caption":   newPhoto.Caption,
		"title":     newPhoto.Title,
		"photo_url": newPhoto.PhotoURL,
		"user_id":   newPhoto.UserID,
	})
}

func (pc *PhotoController) UpdatePhoto(ctx *gin.Context) {
	photoID := ctx.Param("photoId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.UpdatePhoto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Title is required"})
		return
	}

	if payload.PhotoURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Photo URL is required"})
		return
	}

	if !utils.IsValidURL(payload.PhotoURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid photo URL format"})
		return
	}

	var updatedPhoto models.Photo
	result := pc.DB.First(&updatedPhoto, "id = ?", photoID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No photo with that ID exists"})
		return
	}

	if updatedPhoto.UserID != currentUser.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You are not authorized to update this photo"})
		return
	}

	now := time.Now()
	updatedPhoto.Title = payload.Title
	updatedPhoto.Caption = payload.Caption
	updatedPhoto.PhotoURL = payload.PhotoURL
	updatedPhoto.UpdatedAt = now

	pc.DB.Save(&updatedPhoto)

	responseData := gin.H{
		"id":        updatedPhoto.ID,
		"caption":   updatedPhoto.Caption,
		"title":     updatedPhoto.Title,
		"photo_url": updatedPhoto.PhotoURL,
		"user_id":   updatedPhoto.UserID,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

func (pc *PhotoController) FindPhotoByID(ctx *gin.Context) {
	photoID := ctx.Param("photoId")

	var photo models.Photo
	result := pc.DB.First(&photo, "id = ?", photoID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No photo with that ID exists"})
		return
	}

	user := models.User{}
	pc.DB.First(&user, photo.UserID)

	responseData := gin.H{
		"id":        photo.ID,
		"caption":   photo.Caption,
		"title":     photo.Title,
		"photo_url": photo.PhotoURL,
		"user_id":   photo.UserID,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}
func (pc *PhotoController) FindPhotos(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var photos []models.Photo
	results := pc.DB.Limit(intLimit).Offset(offset).Find(&photos)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	var responseData []gin.H
	for _, photo := range photos {
		user := models.User{}
		pc.DB.First(&user, photo.UserID)

		responseData = append(responseData, gin.H{
			"id":        photo.ID,
			"caption":   photo.Caption,
			"title":     photo.Title,
			"photo_url": photo.PhotoURL,
			"user_id":   photo.UserID,
			"user": gin.H{
				"id":       user.ID,
				"email":    user.Email,
				"username": user.Username,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

func (pc *PhotoController) DeletePhoto(ctx *gin.Context) {
	photoID := ctx.Param("photoId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var photo models.Photo
	result := pc.DB.First(&photo, "id = ?", photoID)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No photo with that ID exists"})
		return
	}

	if photo.UserID != currentUser.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You are not authorized to delete this photo"})
		return
	}

	result = pc.DB.Delete(&photo)

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No photo with that ID exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
