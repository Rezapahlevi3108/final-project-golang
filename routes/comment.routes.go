package routes

import (
	"final-project-golang/controllers"
	"final-project-golang/middleware"

	"github.com/gin-gonic/gin"
)

type CommentRouteController struct {
	commentController controllers.CommentController
}

func NewRouteCommentController(commentController controllers.CommentController) CommentRouteController {
	return CommentRouteController{commentController}
}

func (cc *CommentRouteController) CommentRoute(rg *gin.RouterGroup) {
	router := rg.Group("comments")
	router.Use(middleware.DeserializeUser())
	router.POST("", cc.commentController.CreateComment)
	router.GET("", cc.commentController.GetComments)
	router.PUT("/:commentId", cc.commentController.UpdateComment)
	router.GET("/:commentId", cc.commentController.GetCommentByID)
	router.DELETE("/:commentId", cc.commentController.DeleteComment)
}
