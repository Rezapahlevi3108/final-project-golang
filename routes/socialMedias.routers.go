package routes

import (
	"final-project-golang/controllers"
	"final-project-golang/middleware"

	"github.com/gin-gonic/gin"
)

type SocialMediaRouteController struct {
	socialMediaController controllers.SocialMediaController
}

func NewRouteSocialMediaController(socialMediaController controllers.SocialMediaController) SocialMediaRouteController {
	return SocialMediaRouteController{socialMediaController}
}

func (smc *SocialMediaRouteController) SocialMediaRoute(rg *gin.RouterGroup) {
	router := rg.Group("socialmedias")
	router.Use(middleware.DeserializeUser())
	router.POST("", smc.socialMediaController.CreateSocialMedia)
	router.GET("", smc.socialMediaController.GetSocialMedias)
	router.GET("/:socialMediaId", smc.socialMediaController.GetSocialMediaByID)
	router.PUT("/:socialMediaId", smc.socialMediaController.UpdateSocialMedia)
	router.DELETE("/:socialMediaId", smc.socialMediaController.DeleteSocialMedia)
}
