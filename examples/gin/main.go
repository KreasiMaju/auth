package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	// Setup Gin
	r := gin.Default()

	// Route publik
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to API!")
	})

	// Auth routes
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/login", loginPage)
		authGroup.POST("/login", login)
		authGroup.GET("/register", registerPage)
		authGroup.POST("/register", register)
		authGroup.GET("/google", googleLogin)
		authGroup.GET("/google/callback", googleCallback)
		authGroup.GET("/logout", logout)
	}

	// Rute yang dilindungi
	api := r.Group("/api")
	api.Use(auth.GinMiddleware()) // Tambahkan middleware auth
	{
		api.GET("/profile", getProfile)

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(auth.GinRoleMiddleware("admin")) // Tambahkan middleware peran
		{
			admin.GET("/users", listUsers)
		}
	}

	// Mulai server
	r.Run(":8080")
}

// Handler untuk login page
func loginPage(c *gin.Context) {
	c.String(http.StatusOK, "Login Page")
}

// Handler untuk login
func login(c *gin.Context) {
	// Implementasi login
	c.String(http.StatusOK, "Login Success")
}

// Handler untuk register page
func registerPage(c *gin.Context) {
	c.String(http.StatusOK, "Register Page")
}

// Handler untuk register
func register(c *gin.Context) {
	// Implementasi register
	c.String(http.StatusOK, "Register Success")
}

// Handler untuk Google login
func googleLogin(c *gin.Context) {
	// Implementasi Google login
	c.String(http.StatusOK, "Google Login")
}

// Handler untuk Google callback
func googleCallback(c *gin.Context) {
	// Implementasi Google callback
	c.String(http.StatusOK, "Google Callback")
}

// Handler untuk logout
func logout(c *gin.Context) {
	// Implementasi logout
	c.String(http.StatusOK, "Logout Success")
}

// Handler untuk profile
func getProfile(c *gin.Context) {
	// Implementasi get profile
	c.String(http.StatusOK, "Profile")
}

// Handler untuk list users
func listUsers(c *gin.Context) {
	// Implementasi list users
	c.String(http.StatusOK, "Users")
}
