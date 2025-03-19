package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kreasimaju/auth/utils"
)

// FiberAuthMiddleware adalah middleware autentikasi untuk Fiber
func FiberAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Mendapatkan token dari header
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Format harus berupa "Bearer TOKEN"
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header format must be Bearer TOKEN",
			})
		}

		tokenString := parts[1]

		// Validasi token JWT
		token, err := utils.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Ambil klaim dari token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not parse token claims",
			})
		}

		// Tetapkan klaim user ke konteks
		c.Locals("user", claims)

		return c.Next()
	}
}

// FiberRoleMiddleware adalah middleware untuk memeriksa peran pengguna pada Fiber
func FiberRoleMiddleware(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Mendapatkan user dari konteks (yang diatur oleh middleware auth)
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		claims, ok := user.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user claims",
			})
		}

		// Memeriksa peran pengguna
		userRole, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "User has no role assigned",
			})
		}

		// Periksa apakah peran pengguna ada dalam daftar peran yang diizinkan
		for _, role := range roles {
			if role == userRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User does not have the required role",
		})
	}
}
