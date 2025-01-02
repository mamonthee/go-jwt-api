package routes

import (
	"go-jwt-api/controllers"
	"go-jwt-api/middleware"

	"github.com/gin-gonic/gin"
)

func ArticleRoutes(route *gin.Engine) {
	articleRoutes := route.Group("/articles")
	articleRoutes.Use(middleware.AuthenticateJWT())
	articleRoutes.POST("/", controllers.CreateArticle())
	articleRoutes.GET("/", controllers.GetArticles())
	articleRoutes.PUT("/:id", controllers.UpdateArticle())
	articleRoutes.DELETE("/:id", controllers.DeleteArticle())
}
