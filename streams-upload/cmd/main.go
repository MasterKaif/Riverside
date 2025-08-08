package main


package main

import (
	"log"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting server...")
	// utils.InitDB()
	// utils.InitKafka()

	// Initialize S3
	// if err := utils.InitS3(); err != nil {
	// 	log.Fatalf("Failed to initialize S3: %v", err)
	// }
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
	a

	// Auth routes
	


	log.Println("About to start Gin server on port 8080")
	err := r.Run("0.0.0.0:8000")
	if err != nil {
		log.Fatalf("Gin server failed to start: %v", err)
	}
	log.Println("Server started on port 8080")
}
