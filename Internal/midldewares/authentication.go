package middlewares

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MasterKaif/RiverSide/Internal/models"
	"github.com/MasterKaif/RiverSide/Internal/usecases"
	"github.com/MasterKaif/RiverSide/Internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}
		tokenString = tokenParts[1]
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		claims, err := usecases.ParseJWT(tokenString)
		if err != nil {
			log.Println("Error validating token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		var user models.Users
		if err := utils.DB.Where("id = ?", claims.Sub).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		c.Set("userID", claims.Sub)
		// c.Set("iat", claims.IssuedAt)
		// c.Set("exp", claims.ExpiresAt)

		c.Next()
	}
}
