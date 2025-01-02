package routes

import (
	"go-jwt-api/controllers"
	"go-jwt-api/middleware"

	"github.com/gin-gonic/gin"
)

func AuthorRoutes(route *gin.Engine) {

	route.POST("/register", controllers.Register())
	route.POST("/login", controllers.Login())

	authorRoutes := route.Group("/author")
	authorRoutes.Use(middleware.AuthenticateJWT())
	authorRoutes.PUT("/update", controllers.UpdateAuthor())
	authorRoutes.PUT("/deactivate", controllers.Deactivate())
}
