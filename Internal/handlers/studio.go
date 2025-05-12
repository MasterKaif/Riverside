package handlers

import (
	"net/http"
	"time"

	"github.com/MasterKaif/RiverSide/Internal/models"
	"github.com/MasterKaif/RiverSide/Internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StudioCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type StudioCreateResponse struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Host        UserCreateResponse `json:"host"` // Assuming you have a UserCreateResponse struct for the host
}

func StudioCreateHandler(c *gin.Context) {
	var req StudioCreateRequest
	if err := c.BindJSON((&req)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}
	var host models.Users
	if err := utils.DB.Where("id = ?", c.MustGet("userID")).First(&host).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}

	session := models.StudioSession{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Host:        host,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := utils.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create studio"})
		return
	}

	sessionDetials := StudioCreateResponse{
		ID:          session.ID,
		Name:        session.Name,
		Description: session.Description,
		Host: UserCreateResponse{
			ID:       host.ID,
			Username: host.Username,
			Email:    host.Email,
		},
	}
	c.JSON(http.StatusOK, gin.H{"message": "Studio created successfully", "session": sessionDetials})
}

func StudioJoinHandler(c *gin.Context) {
	type JoinRequest struct {
		SessionID string `json:"session_id" binding:"required"`
	}

	type JoinResponse struct {
		ID      uuid.UUID            `json:"id"`
		Session StudioCreateResponse `json:"session"`
	}

	var req JoinRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}
	var session models.StudioSession
	var user models.Users

	if err := utils.DB.Where("id = ?", c.MustGet("userID")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := utils.DB.Where("id = ?", req.SessionID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Studio not found"})
		return
	}

	// Check if the user is already a joiner
	var joiner models.SessionJoiners
	if err := utils.DB.Where("session_id = ? AND user_id = ?", session.ID, user.ID).First(&joiner).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already joined the studio"})
		return
	}

	joiner = models.SessionJoiners{
		ID:      uuid.New(),
		Session: session,
		User:    user,
	}

	if err := utils.DB.Create(&joiner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join studio"})
		return
	}

	sessionDetails := JoinResponse{
		ID: joiner.ID,
		Session: StudioCreateResponse{
			ID:          session.ID,
			Name:        session.Name,
			Description: session.Description,
			Host: UserCreateResponse{
				ID:       session.Host.ID,
				Username: session.Host.Username,
				Email:    session.Host.Email,
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined studio successfully", "session": sessionDetails})
}
