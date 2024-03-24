package routes

import (
	"final-project-golang/controllers"

	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
	userController controllers.UserController
}

func NewAuthRouteController(authController controllers.AuthController, userController controllers.UserController) AuthRouteController {
	return AuthRouteController{authController, userController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("users")

	router.POST("/register", rc.authController.SignUpUser)
	router.POST("/login", rc.authController.SignInUser)
}
