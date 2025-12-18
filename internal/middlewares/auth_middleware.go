package middlewares

import (
	"net/http"
	"strings"

	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/services"
	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "user_id"

func AuthRequired(jwtSvc *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing bearer token"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")

		claims, err := jwtSvc.ParseAccess(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		userIDf, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token claims"})
			return
		}

		c.Set(CtxUserIDKey, uint(userIDf))
		c.Next()
	}
}
