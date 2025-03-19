package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kreasimaju/auth"
	"github.com/kreasimaju/auth/config"
)

func main() {
	// Inisialisasi auth dengan konfigurasi
	err := auth.Init(config.Config{
		Database: config.Database{
			Type:        "sqlite",  // SQLite untuk contoh
			Database:    "auth.db", // File database
			AutoMigrate: true,
		},
		Providers: config.Providers{
			Google: config.OAuth{
				Enabled:      true,
				ClientID:     "your-client-id",
				ClientSecret: "your-client-secret",
				CallbackURL:  "http://localhost:8080/auth/google/callback",
			},
			Local: true, // Aktifkan auth lokal
		},
		JWT: config.JWT{
			Secret:    "your-jwt-secret-key",
			ExpiresIn: 86400, // 24 jam
		},
	})

	if err != nil {
		panic("Failed to initialize auth: " + err.Error())
	}

	// Setup Fiber
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Route publik
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to API!")
	})

	// Auth routes
	authGroup := app.Group("/auth")
	{
		authGroup.Get("/login", loginPage)
		authGroup.Post("/login", login)
		authGroup.Get("/register", registerPage)
		authGroup.Post("/register", register)
		authGroup.Get("/google", googleLogin)
		authGroup.Get("/google/callback", googleCallback)
		authGroup.Get("/logout", logout)
	}

	// Rute yang dilindungi
	api := app.Group("/api")
	api.Use(auth.FiberMiddleware()) // Tambahkan middleware auth
	{
		api.Get("/profile", getProfile)

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(auth.FiberRoleMiddleware("admin")) // Tambahkan middleware peran
		{
			admin.Get("/users", listUsers)
		}
	}

	// Mulai server
	app.Listen(":8080")
}

// Handler untuk login page
func loginPage(c *fiber.Ctx) error {
	return c.SendString("Login Page")
}

// Handler untuk login
func login(c *fiber.Ctx) error {
	// Implementasi login
	return c.SendString("Login Success")
}

// Handler untuk register page
func registerPage(c *fiber.Ctx) error {
	return c.SendString("Register Page")
}

// Handler untuk register
func register(c *fiber.Ctx) error {
	// Implementasi register
	return c.SendString("Register Success")
}

// Handler untuk Google login
func googleLogin(c *fiber.Ctx) error {
	// Implementasi Google login
	return c.SendString("Google Login")
}

// Handler untuk Google callback
func googleCallback(c *fiber.Ctx) error {
	// Implementasi Google callback
	return c.SendString("Google Callback")
}

// Handler untuk logout
func logout(c *fiber.Ctx) error {
	// Implementasi logout
	return c.SendString("Logout Success")
}

// Handler untuk profile
func getProfile(c *fiber.Ctx) error {
	// Implementasi get profile
	return c.SendString("Profile")
}

// Handler untuk list users
func listUsers(c *fiber.Ctx) error {
	// Implementasi list users
	return c.SendString("Users")
}
