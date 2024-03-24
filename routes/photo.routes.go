package routes

import (
	"final-project-golang/controllers"
	"final-project-golang/middleware"

	"github.com/gin-gonic/gin"
)

type PhotoRouteController struct {
	photoController controllers.PhotoController
}

func NewRoutePhotoController(photoController controllers.PhotoController) PhotoRouteController {
	return PhotoRouteController{photoController}
}

func (pc *PhotoRouteController) PhotoRoute(rg *gin.RouterGroup) {

	router := rg.Group("photos")
	router.Use(middleware.DeserializeUser())
	router.POST("", pc.photoController.CreatePhoto)
	router.GET("", pc.photoController.FindPhotos)
	router.PUT("/:photoId", pc.photoController.UpdatePhoto)
	router.GET("/:photoId", pc.photoController.FindPhotoByID)
	router.DELETE("/:photoId", pc.photoController.DeletePhoto)
}
