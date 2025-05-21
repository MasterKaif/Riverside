package main

import (
	"log"
	"time"

	"github.com/MasterKaif/RiverSide/Internal/handlers"
	middlewares "github.com/MasterKaif/RiverSide/Internal/midldewares"
	"github.com/MasterKaif/RiverSide/Internal/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting server...")
	utils.InitDB()
	r := gin.Default()
	log.Println("Initializing Gin server...")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api/v1")

	// WebSocket route
	api.GET("/ws", handlers.WebSocketHandler)

	// health check route
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Auth routes
	api.POST("/auth/login", handlers.LoginHandler)
	api.POST("/auth/signup", handlers.SignupHandler)
	api.GET("/auth/me", middlewares.AuthMiddleware(), handlers.ValidateTokenHandler)

	protected := api.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/studio/create", handlers.StudioCreateHandler)
		protected.POST("/studio/join", handlers.StudioJoinHandler)
	}

	log.Println("About to start Gin server on port 8080")
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Gin server failed to start: %v", err)
	}
	log.Println("Server started on port 8080")
}
