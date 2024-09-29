package middleware

import (
	"fmt"
	"net/http"
	"server/models"
	"server/types"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Authorization header is missing"})
			c.Abort()
			return
		}

		userData, userDataErr := models.VerifyJWTToken(tokenString)
		if userDataErr != nil {
			fmt.Println(userDataErr)
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid token"})
			c.Abort()
			return
		}

		user := &types.User{
			ID:    userData.ID,
			Name:  userData.Name,
			Email: userData.Email,
		}
		c.Set("user", user)

		c.Next()
	}
}
