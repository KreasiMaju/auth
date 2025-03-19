package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kreasimaju/auth/utils"
	"github.com/labstack/echo/v4"
)

// EchoAuthMiddleware adalah middleware autentikasi untuk Echo
func EchoAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Mendapatkan token dari header
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header is required",
				})
			}

			// Format harus berupa "Bearer TOKEN"
			parts := strings.Split(auth, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header format must be Bearer TOKEN",
				})
			}

			tokenString := parts[1]

			// Validasi token JWT
			token, err := utils.ValidateJWT(tokenString)
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
			}

			// Ambil klaim dari token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Could not parse token claims",
				})
			}

			// Tetapkan klaim user ke konteks
			c.Set("user", claims)

			return next(c)
		}
	}
}

// EchoRoleMiddleware adalah middleware untuk memeriksa peran pengguna pada Echo
func EchoRoleMiddleware(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Mendapatkan user dari konteks (yang diatur oleh middleware auth)
			user, ok := c.Get("user").(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "User not authenticated",
				})
			}

			// Memeriksa peran pengguna
			userRole, ok := user["role"].(string)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "User has no role assigned",
				})
			}

			// Periksa apakah peran pengguna ada dalam daftar peran yang diizinkan
			for _, role := range roles {
				if role == userRole {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "User does not have the required role",
			})
		}
	}
}
