package main

import (
	"go-jwt-api/database"
	"go-jwt-api/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDatabase()
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	r := gin.New()
	r.Use(gin.Logger())

	routes.AuthorRoutes(r)
	routes.ArticleRoutes(r)

	r.Run(":" + port)
}
