package routes

import (
	"final-project-golang/controllers"
	"final-project-golang/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("users")
	router.PUT("", middleware.DeserializeUser(), uc.userController.UpdateMe)
	router.DELETE("", middleware.DeserializeUser(), uc.userController.DeleteMe)
}
