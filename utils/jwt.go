package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kreasimaju/auth/config"
	"github.com/kreasimaju/auth/models"
)

var jwtSecret string

// SetJWTSecret menetapkan secret untuk JWT
func SetJWTSecret(secret string) {
	jwtSecret = secret
}

// GenerateJWT menghasilkan token JWT untuk pengguna
func GenerateJWT(user models.User, cfg config.JWT) (string, error) {
	// Buat token JWT dengan klaim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":     user.ID,
		"email":       user.Email,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"role":        user.Role,
		"is_verified": user.IsVerified,
		"exp":         time.Now().Add(time.Second * time.Duration(cfg.ExpiresIn)).Unix(),
		"iat":         time.Now().Unix(),
	})

	// Tandatangani token dengan secret
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT memvalidasi token JWT
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	if jwtSecret == "" {
		return nil, errors.New("JWT secret not initialized")
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi metode signing
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetUserIDFromToken mengambil ID pengguna dari token JWT
func GetUserIDFromToken(token *jwt.Token) (uint, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("could not parse token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id claim not found or invalid")
	}

	return uint(userID), nil
}

// GetConfig mengambil konfigurasi dari package auth
func GetConfig() config.Config {
	// Ideally, this would import the config from the auth package
	// But to avoid circular dependencies, we might need to redesign this
	// For now, we'll return a mock config for compilation
	return config.Config{}
}
