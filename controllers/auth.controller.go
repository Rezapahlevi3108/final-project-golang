package controllers

import (
	"net/http"
	"strings"
	"time"

	"final-project-golang/initializers"
	"final-project-golang/models"
	"final-project-golang/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if !utils.IsValidEmail(payload.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email format"})
		return
	}

	var existingUser models.User
	if ac.DB.Where("email = ?", strings.ToLower(payload.Email)).First(&existingUser).RowsAffected != 0 {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	}

	if payload.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Username is required"})
		return
	}
	if ac.DB.Where("username = ?", payload.Username).First(&existingUser).RowsAffected != 0 {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Username already exists"})
		return
	}

	if len(payload.Password) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Password must be at least 6 characters long"})
		return
	}

	if payload.Age < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Age must be at least 8"})
		return
	}

	if payload.ProfileImageURL != "" && !utils.IsValidURL(payload.ProfileImageURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid profile image URL format"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User{
		Username:        payload.Username,
		Email:           strings.ToLower(payload.Email),
		Password:        hashedPassword,
		Age:             payload.Age,
		ProfileImageURL: payload.ProfileImageURL,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	userResponse := gin.H{
		"id":                newUser.ID,
		"email":             newUser.Email,
		"username":          newUser.Username,
		"age":               newUser.Age,
		"profile_image_url": newUser.ProfileImageURL,
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": userResponse})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Invalid email or password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Invalid email or password"})
		return
	}

	config, _ := initializers.LoadConfig(".")

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}
