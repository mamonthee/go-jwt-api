package controllers

import (
	"fmt"
	"go-jwt-api/database"
	"go-jwt-api/helpers"
	"go-jwt-api/models"
	"go-jwt-api/redis"
	"go-jwt-api/response"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providePassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providePassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}

func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var author models.Author

		// Bind JSON body to the author struct
		if err := ctx.BindJSON(&author); err != nil {
			response.SendErrorResponse(ctx, err.Error(), nil)
			// ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		// Validate the author
		validationErr := validate.Struct(author)
		if validationErr != nil {
			response.SendErrorResponse(ctx, validationErr.Error(), nil)
			// ctx.JSON(http.StatusBadRequest, gin.H{"err": validationErr.Error()})
			return
		}

		// Check if email or phone already exists
		var existingAuthor models.Author // Check if email already exists

		if err := database.DB.Where("email = ?", author.Email).First(&existingAuthor).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				fmt.Println("Email does not exist, proceeding to create")
			} else {
				log.Printf("Database error: %v\n", err)
				response.SendErrorResponse(ctx, "Database error: "+err.Error(), nil)
				// ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		} else {
			response.SendErrorResponse(ctx, "Email already exists", nil)
			return
		}

		// Hash the password
		password := HashPassword(author.Password)

		author.Password = string(password)
		author.IsActive = true

		// Save the user to the database
		if err := database.DB.Create(&author).Error; err != nil {
			response.SendErrorResponse(ctx, "Failed to create author", nil)
			return
		}

		// Respond with the created user
		response.SendSuccessResponse(ctx, "Author created successfully", author)
		// ctx.JSON(http.StatusOK, gin.H{"message": "Author created successfully", "author": author})
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var author models.Author
		var foundAuthor models.Author

		// Bind JSON body to the author struct
		if err := ctx.BindJSON(&author); err != nil {
			response.SendErrorResponse(ctx, err.Error(), nil)
			return
		}

		// Find the author in the database
		result := database.DB.Where("email = ?", author.Email).First(&foundAuthor)

		if result.Error != nil || result.RowsAffected == 0 {
			response.SendErrorResponse(ctx, "Invalid email!", nil)
			return
		}

		// Verify the password
		passwordIsValid, msg := VerifyPassword(author.Password, foundAuthor.Password)
		// Check if the account is active
		if !foundAuthor.IsActive {
			response.SendErrorResponse(ctx, "Account is deactivated", nil)
			return
		}

		if passwordIsValid != true {
			response.SendErrorResponse(ctx, msg, nil)
			return
		}

		if foundAuthor.Email == "" {
			response.SendErrorResponse(ctx, "email not found", nil)
			return
		}

		// Generate a JWT token
		token, refreshToken, err := helpers.GenerateTokens(foundAuthor.ID)
		if err != nil {
			response.SendErrorResponse(ctx, err.Error(), nil)
			return
		}

		// Return the token
		response.SendSuccessResponse(ctx, "Login successful", gin.H{
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

func UpdateAuthor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the logged-in author's ID from the JWT token
		loggedInAuthorID, exists := ctx.Get("author_id")
		if !exists {
			response.SendErrorResponse(ctx, "Unauthorized", nil)
			return
		}

		// Get the author's details from the database
		var author models.Author
		if err := database.DB.First(&author, loggedInAuthorID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.SendErrorResponse(ctx, "Author not found", nil)
			} else {
				response.SendErrorResponse(ctx, "Failed to fetch author", nil)
			}
			return
		}

		// Check if the account is active
		if !author.IsActive {
			response.SendErrorResponse(ctx, "Account is deactivated", nil)
			return
		}

		var input models.Author
		if err := ctx.BindJSON(&input); err != nil {
			response.SendErrorResponse(ctx, err.Error(), nil)
			return
		}

		// Update the author's details
		if input.UserName != "" {
			author.UserName = input.UserName
		}
		if input.Email != "" && input.Email != author.Email {
			// Make sure the new email doesn't already exist
			var existingAuthor models.Author
			if err := database.DB.Where("email = ?", input.Email).First(&existingAuthor).Error; err == nil {
				response.SendErrorResponse(ctx, "Email is already in use", nil)
				return
			}
			author.Email = input.Email
		}
		if input.Password != "" {
			// Hash the new password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
			if err != nil {
				response.SendErrorResponse(ctx, "Failed to hash password", nil)
				return
			}
			author.Password = string(hashedPassword)
		}

		// Save the updated author details
		if err := database.DB.Save(&author).Error; err != nil {
			response.SendErrorResponse(ctx, "Failed to update author", nil)
			return
		}

		// Return a success response
		response.SendSuccessResponse(ctx, "Author updated successfully", author)
	}
}

func Deactivate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the logged-in author's ID from the JWT token
		loggedInAuthorID, exists := ctx.Get("author_id")
		if !exists {
			response.SendErrorResponse(ctx, "Unauthorized", nil)
			return
		}

		// Get the author's details from the database
		var author models.Author
		if err := database.DB.First(&author, loggedInAuthorID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.SendErrorResponse(ctx, "Author not found", nil)
			} else {
				response.SendErrorResponse(ctx, "Failed to fetch author", nil)
			}
			return
		}

		// Set the author's IsActive field to false to deactivate the account
		author.IsActive = false

		// Save the updated author details
		if err := database.DB.Save(&author).Error; err != nil {
			response.SendErrorResponse(ctx, "Failed to deactivate author", nil)
			return
		}

		// Blacklist token
		token := ctx.GetHeader("Authorization")[7:]

		err := redis.AddTokenToBlacklist(token, time.Hour*24) // Use the token's expiry time
		if err != nil {
			response.SendErrorResponse(ctx, "Could not blacklist token", nil)
			return
		}

		// Return a success response
		response.SendSuccessResponse(ctx, "Author account deactivated successfully", nil)
	}
}
