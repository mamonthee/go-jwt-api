package controllers

import (
	"go-jwt-api/database"
	"go-jwt-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateArticle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var article models.Article
		if err := ctx.BindJSON(&article); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		// Get AuthorID from middleware
		authorID, exists := ctx.Get("author_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		article.AuthorID = authorID.(uint)
		if err := database.DB.Create(&article).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
			return
		}

		// Preload the Author relationship
		if err := database.DB.Preload("Author").First(&article, article.ID).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load author"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Article created successfully", "article": article})
	}
}

// GetArticles retrieves all articles for an authenticated author
func GetArticles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var articles []models.Article

		authorID, exists := ctx.Get("author_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if err := database.DB.Preload("Author").Where("author_id = ?", authorID).Find(&articles).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"articles": articles})
	}
}

func UpdateArticle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		articleID := ctx.Param("id")
		var article models.Article

		// Find the article by ID
		if err := database.DB.First(&article, articleID).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}

		// Retrieve the authorID from the JWT token in the context
		loggedInAuthorIDStr, exists := ctx.Get("author_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Ensure the logged-in user is the author of the article
		if article.AuthorID != loggedInAuthorIDStr {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own articles"})
			return
		}

		// Bind the update data
		var inputData models.Article
		if err := ctx.ShouldBindJSON(&inputData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the article fields
		article.Title = inputData.Title
		article.Description = inputData.Description

		// Update only if fields are provided
		if inputData.Title != "" {
			article.Title = inputData.Title
		}
		if inputData.Description != "" {
			article.Description = inputData.Description
		}

		// Save to the database
		if err := database.DB.Save(&article).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Article updated successfully", "article": article})
	}
}

func DeleteArticle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		articleID := ctx.Param("id")
		var article models.Article

		// Find the article by ID
		if err := database.DB.First(&article, articleID).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}

		// Retrieve the authorID from the JWT token in the context
		loggedInAuthorIDStr, exists := ctx.Get("author_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Ensure the logged-in user is the author of the article
		if article.AuthorID != loggedInAuthorIDStr {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own articles"})
			return
		}

		// Delete the article
		if err := database.DB.Delete(&article).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
	}
}
