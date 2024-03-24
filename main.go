package main

import (
	"log"

	"final-project-golang/controllers"
	"final-project-golang/initializers"
	"final-project-golang/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	PhotoController      controllers.PhotoController
	PhotoRouteController routes.PhotoRouteController

	CommentController      controllers.CommentController
	CommentRouteController routes.CommentRouteController

	SocialMediaController      controllers.SocialMediaController
	SocialMediaRouteController routes.SocialMediaRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController, UserController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	PhotoController = controllers.NewPhotoController(initializers.DB)
	PhotoRouteController = routes.NewRoutePhotoController(PhotoController)

	CommentController = controllers.NewCommentController(initializers.DB)
	CommentRouteController = routes.NewRouteCommentController(CommentController)

	SocialMediaController = controllers.NewSocialMediaController(initializers.DB)
	SocialMediaRouteController = routes.NewRouteSocialMediaController(SocialMediaController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	AuthRouteController.AuthRoute(&server.RouterGroup)
	UserRouteController.UserRoute(&server.RouterGroup)
	PhotoRouteController.PhotoRoute(&server.RouterGroup)
	CommentRouteController.CommentRoute(&server.RouterGroup)
	SocialMediaRouteController.SocialMediaRoute(&server.RouterGroup)
	log.Fatal(server.Run(":" + config.ServerPort))
}
