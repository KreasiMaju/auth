package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kreasimaju/auth/utils"
)

// GinAuthMiddleware adalah middleware autentikasi untuk Gin
func GinAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mendapatkan token dari header
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Format harus berupa "Bearer TOKEN"
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header format must be Bearer TOKEN",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validasi token JWT
		token, err := utils.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Ambil klaim dari token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not parse token claims",
			})
			c.Abort()
			return
		}

		// Tetapkan klaim user ke konteks
		c.Set("user", claims)

		c.Next()
	}
}

// GinRoleMiddleware adalah middleware untuk memeriksa peran pengguna pada Gin
func GinRoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mendapatkan user dari konteks (yang diatur oleh middleware auth)
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			c.Abort()
			return
		}

		claims, ok := user.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user claims",
			})
			c.Abort()
			return
		}

		// Memeriksa peran pengguna
		userRole, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User has no role assigned",
			})
			c.Abort()
			return
		}

		// Periksa apakah peran pengguna ada dalam daftar peran yang diizinkan
		for _, role := range roles {
			if role == userRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "User does not have the required role",
		})
		c.Abort()
	}
}
