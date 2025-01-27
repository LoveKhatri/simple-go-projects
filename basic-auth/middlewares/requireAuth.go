package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/LoveKhatri/basic-auth/initializers"
	"github.com/LoveKhatri/basic-auth/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")

	if err != nil || tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Token is invalid",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token expired",
			})
			return
		}

		var user models.User

		initializers.DB.First(&user, claims["sub"])

		if uint(claims["sub"].(float64)) != user.ID {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token is invalid",
			})
			return
		}

		c.Set("user", user)

		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Token is invalid",
		})
		return
	}
}
