package models

import (
	"gorm.io/gorm"
)

// OAuth model untuk autentikasi OAuth
type OAuth struct {
	gorm.Model
	UserID       uint   `gorm:"index" json:"user_id"`
	Provider     string `gorm:"type:varchar(20)" json:"provider"` // "google", "github", "facebook", dll
	ProviderID   string `gorm:"type:varchar(100)" json:"provider_id"`
	AccessToken  string `gorm:"type:varchar(255)" json:"access_token,omitempty"`
	RefreshToken string `gorm:"type:varchar(255)" json:"refresh_token,omitempty"`
	Email        string `gorm:"type:varchar(100)" json:"email"`
	Name         string `gorm:"type:varchar(100)" json:"name"`
	AvatarURL    string `gorm:"type:varchar(255)" json:"avatar_url"`
	Data         string `gorm:"type:text" json:"data"` // Data tambahan dalam format JSON
}

// FindUserOAuth mencari data OAuth berdasarkan provider dan provider ID
func FindUserOAuth(db *gorm.DB, provider, providerID string) (*OAuth, error) {
	var oauth OAuth
	err := db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&oauth).Error
	if err != nil {
		return nil, err
	}
	return &oauth, nil
}
