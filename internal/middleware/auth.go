package middleware

import (
	"library-api/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAuth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		userID, err := authService.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
