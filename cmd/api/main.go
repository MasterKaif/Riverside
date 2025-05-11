package main

import (
	"log"

	"github.com/MasterKaif/RiverSide/Internal/handlers"
	"github.com/MasterKaif/RiverSide/Internal/midldewares"
	"github.com/MasterKaif/RiverSide/Internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitDB()
	r := gin.Default() 

	api := r.Group("/api")

	// Auth routes
	api.POST("/auth/login", handlers.LoginHandler)
	api.POST("/auth/signup", handlers.SignupHandler)

	protected := api.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/studio/create", handlers.StudioCreateHandler)
		protected.POST("/studio/join", handlers.StudioJoinHandler)
	}

	r.Run(":8000")
	log.Println("Server started on port 8000")

}
