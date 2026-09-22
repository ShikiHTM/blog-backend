package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shikihtm/blog-backend/internal/model"
)

func RequireAuthCookie(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "TOKEN_NOT_FOUND",
			})
		}

		token, err := jwt.ParseWithClaims(tokenString, &model.JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "TOKEN_INVALID",
			})
		}

		if claims, ok := token.Claims.(*model.JwtClaims); ok {
			c.Set("username", claims.Username)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "INVALID_PAYLOAD",
			})
		}
	}
}
