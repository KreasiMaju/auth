package main

import (
	"net/http"

	"github.com/kreasimaju/auth"
	"github.com/kreasimaju/auth/config"
	"github.com/kreasimaju/auth/middleware"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
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

	// Setup Echo
	e := echo.New()

	// Middleware
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.CORS())

	// Route publik
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to API!")
	})

	// Auth routes
	auth := e.Group("/auth")
	{
		auth.GET("/login", loginPage)
		auth.POST("/login", login)
		auth.GET("/register", registerPage)
		auth.POST("/register", register)
		auth.GET("/google", googleLogin)
		auth.GET("/google/callback", googleCallback)
		auth.GET("/logout", logout)
	}

	// Rute yang dilindungi
	api := e.Group("/api")
	api.Use(middleware.EchoAuthMiddleware()) // Tambahkan middleware auth
	{
		api.GET("/profile", getProfile)

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.EchoRoleMiddleware("admin")) // Tambahkan middleware peran
		{
			admin.GET("/users", listUsers)
		}
	}

	// Mulai server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler untuk login page
func loginPage(c echo.Context) error {
	return c.String(http.StatusOK, "Login Page")
}

// Handler untuk login
func login(c echo.Context) error {
	// Implementasi login
	return c.String(http.StatusOK, "Login Success")
}

// Handler untuk register page
func registerPage(c echo.Context) error {
	return c.String(http.StatusOK, "Register Page")
}

// Handler untuk register
func register(c echo.Context) error {
	// Implementasi register
	return c.String(http.StatusOK, "Register Success")
}

// Handler untuk Google login
func googleLogin(c echo.Context) error {
	// Implementasi Google login
	return c.String(http.StatusOK, "Google Login")
}

// Handler untuk Google callback
func googleCallback(c echo.Context) error {
	// Implementasi Google callback
	return c.String(http.StatusOK, "Google Callback")
}

// Handler untuk logout
func logout(c echo.Context) error {
	// Implementasi logout
	return c.String(http.StatusOK, "Logout Success")
}

// Handler untuk profile
func getProfile(c echo.Context) error {
	// Implementasi get profile
	return c.String(http.StatusOK, "Profile")
}

// Handler untuk list users
func listUsers(c echo.Context) error {
	// Implementasi list users
	return c.String(http.StatusOK, "Users")
}
