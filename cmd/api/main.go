package main

import (
	"log"
	"time"

	"github.com/MasterKaif/RiverSide/Internal/handlers"
	"github.com/MasterKaif/RiverSide/Internal/midldewares"
	"github.com/MasterKaif/RiverSide/Internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	utils.InitDB()
	r := gin.Default() 

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api/v1")

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

	r.Run(":8000")
	log.Println("Server started on port 8000")

}
