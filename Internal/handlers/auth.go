package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/MasterKaif/RiverSide/Internal/models"
	"github.com/MasterKaif/RiverSide/Internal/usecases"
	"github.com/MasterKaif/RiverSide/Internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserCreateRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
}

type UserCreateResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func SignupHandler(c *gin.Context) {
	var user UserCreateRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user already exists
	var existingUser models.Users
	if err := utils.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}
	// Here you would typically save the user to a database
	hashedPassword, err := usecases.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	newUser := models.Users{
		ID:        uuid.New(),
		Username:  user.Username,
		Password:  hashedPassword,
		Email:     user.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := utils.DB.Create(&newUser).Error; err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := usecases.GenerateToken(newUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	response := UserCreateResponse{
		ID:       newUser.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": response})
}

func LoginHandler(c *gin.Context) {
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Here you would typically check the user's credentials against a database
	var dbUser models.Users
	if err := utils.DB.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if !usecases.CheckPasswordHash(user.Password, dbUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, error := usecases.GenerateToken(dbUser.ID)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
