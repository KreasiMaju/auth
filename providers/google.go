package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kreasimaju/auth/config"
	"github.com/kreasimaju/auth/models"
	"github.com/kreasimaju/auth/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOAuthConfig *oauth2.Config
)

// InitGoogle menginisialisasi provider Google
func InitGoogle(cfg config.OAuth) {
	googleOAuthConfig = &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.CallbackURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Tambahkan scopes tambahan jika disediakan
	if len(cfg.Scopes) > 0 {
		googleOAuthConfig.Scopes = append(googleOAuthConfig.Scopes, cfg.Scopes...)
	}
}

// GoogleLoginURL mengembalikan URL untuk login melalui Google
func GoogleLoginURL(state string) string {
	return googleOAuthConfig.AuthCodeURL(state)
}

// GoogleUser mewakili respons dari Google API
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// HandleGoogleCallback menangani callback dari Google OAuth
func HandleGoogleCallback(code string) (*models.User, error) {
	// Exchange code untuk token
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	// Ambil data pengguna dari Google
	googleUser, err := getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %s", err.Error())
	}

	// Cari pengguna di database atau buat baru
	user, err := findOrCreateGoogleUser(googleUser, token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// getGoogleUserInfo mengambil informasi pengguna dari Google API
func getGoogleUserInfo(accessToken string) (*GoogleUser, error) {
	// Buat request ke Google API
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Periksa apakah respons berhasil
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status code %d", resp.StatusCode)
	}

	// Baca dan parse respons
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleUser GoogleUser
	if err := json.Unmarshal(data, &googleUser); err != nil {
		return nil, err
	}

	return &googleUser, nil
}

// findOrCreateGoogleUser mencari atau membuat pengguna berdasarkan data Google
func findOrCreateGoogleUser(googleUser *GoogleUser, token *oauth2.Token) (*models.User, error) {
	db := utils.DB
	if db == nil {
		return nil, errors.New("database connection not initialized")
	}

	// Cari provider user
	var provider models.UserProvider
	err := db.Where("provider_name = ? AND provider_id = ?", "google", googleUser.ID).First(&provider).Error

	// Jika provider ditemukan, ambil user dari database
	if err == nil {
		var user models.User
		if err := db.First(&user, provider.UserID).Error; err != nil {
			return nil, err
		}

		// Update token
		provider.AccessToken = token.AccessToken
		provider.RefreshToken = token.RefreshToken
		expiresAt := time.Now().Add(time.Duration(token.Expiry.Unix()-time.Now().Unix()) * time.Second)
		provider.ExpiresAt = &expiresAt

		db.Save(&provider)

		// Update last login
		now := time.Now()
		user.LastLogin = &now
		db.Save(&user)

		return &user, nil
	}

	// Jika provider tidak ditemukan, periksa apakah email sudah terdaftar
	var existingUser models.User
	if err := db.Where("email = ?", googleUser.Email).First(&existingUser).Error; err == nil {
		// Pengguna sudah ada, tambahkan provider baru
		newProvider := models.UserProvider{
			UserID:       existingUser.ID,
			ProviderName: "google",
			ProviderID:   googleUser.ID,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}

		if token.Expiry.Unix() > 0 {
			expiresAt := time.Now().Add(time.Duration(token.Expiry.Unix()-time.Now().Unix()) * time.Second)
			newProvider.ExpiresAt = &expiresAt
		}

		if err := db.Create(&newProvider).Error; err != nil {
			return nil, err
		}

		// Update last login
		now := time.Now()
		existingUser.LastLogin = &now
		db.Save(&existingUser)

		return &existingUser, nil
	}

	// Buat pengguna baru
	newUser := models.User{
		Email:      googleUser.Email,
		FirstName:  googleUser.GivenName,
		LastName:   googleUser.FamilyName,
		IsVerified: googleUser.VerifiedEmail,
		Role:       "user", // Default role
	}

	// Mulai transaction
	tx := db.Begin()

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Tambahkan provider baru
	newProvider := models.UserProvider{
		UserID:       newUser.ID,
		ProviderName: "google",
		ProviderID:   googleUser.ID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	if token.Expiry.Unix() > 0 {
		expiresAt := time.Now().Add(time.Duration(token.Expiry.Unix()-time.Now().Unix()) * time.Second)
		newProvider.ExpiresAt = &expiresAt
	}

	// Simpan data tambahan sebagai JSON
	providerData, _ := json.Marshal(googleUser)
	newProvider.Data = string(providerData)

	if err := tx.Create(&newProvider).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	tx.Commit()

	// Update last login
	now := time.Now()
	newUser.LastLogin = &now
	db.Save(&newUser)

	return &newUser, nil
}
