package controllers

import (
	"go-jwt-api/database"
	"go-jwt-api/models"
	"go-jwt-api/response"

	"github.com/gin-gonic/gin"
)

func CreateArticle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var article models.Article
		if err := ctx.BindJSON(&article); err != nil {
			response.SendErrorResponse(ctx, err.Error(), nil)
			return
		}
		// Get AuthorID from middleware
		authorID, exists := ctx.Get("author_id")
		if !exists {
			response.SendErrorResponse(ctx, "Unauthorized", nil)
			return
		}

		article.AuthorID = authorID.(uint)
		if err := database.DB.Create(&article).Error; err != nil {
			response.SendErrorResponse(ctx, "Failed to create article", nil)
			return
		}

		// Preload the Author relationship
		if err := database.DB.Preload("Author").First(&article, article.ID).Error; err != nil {
			response.SendErrorResponse(ctx, "Failed to load author", nil)
			return
		}

		response.SendSuccessResponse(ctx, "Article created successfully", article)
	}
}

// GetArticles for an authenticated author
func GetArticles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var articles []models.Article

		authorID, exists := ctx.Get("author_id")
		if !exists {
			response.SendErrorResponse(ctx, "Unauthorized", nil)
			return
		}

		if err := database.DB.Preload("Author").Where("author_id = ?", authorID).Find(&articles).Error; err != nil {
			response.SendErrorResponse(ctx, "Failed to retrieve articles", nil)
			return
		}

		response.SendSuccessResponse(ctx, "Article created successfully", articles)
	}
}

func UpdateArticle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		articleID := ctx.Param("id")
		var article models.Article

		// Find the article by ID
		if err := database.DB.First(&article, articleID).Error; err != nil {
			response.SendErrorResponse(ctx, "Article not found", nil)
			return
		}

		// Retrieve the authorID from the JWT token in the context
		loggedInAuthorIDStr, exists := ctx.Get("author_id")
		if !exists {
			response.SendErrorResponse(ctx, "Unauthorized", nil)
			return
		}

		// Ensure the logged-in user is the author of the article
		if article.AuthorID != loggedInAuthorIDStr {
			response.SendErrorResponse(ctx, "You can only update your own articles", nil)
			return
		}

		// Bind the update data
		var inputData models.Article
		if err := ctx.ShouldBindJSON(&inputData); err != nil {
			response.SendErrorResponse(ctx, err.Error(), nil)
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
			response.SendErrorResponse(ctx, err.Error(), nil)
			return
		}
		response.SendSuccessResponse(ctx, "Article updated successfully", article)
	}
}

func DeleteArticle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		articleID := ctx.Param("id")
		var article models.Article

		// Find the article by ID
		if err := database.DB.First(&article, articleID).Error; err != nil {
			response.SendErrorResponse(ctx, "Article not found", nil)
			return
		}

		// Retrieve the authorID from the JWT token in the context
		loggedInAuthorIDStr, exists := ctx.Get("author_id")
		if !exists {
			response.SendErrorResponse(ctx, "Unauthorized", nil)
			return
		}

		// Ensure the logged-in user is the author of the article
		if article.AuthorID != loggedInAuthorIDStr {
			response.SendErrorResponse(ctx, "You can only delete your own articles", nil)
			return
		}

		// Delete the article
		if err := database.DB.Delete(&article).Error; err != nil {
			response.SendErrorResponse(ctx, "Failed to delete article", nil)
			return
		}
		response.SendSuccessResponse(ctx, "Article deleted successfully", nil)
	}
}
