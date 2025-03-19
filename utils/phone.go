package utils

import (
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// FormatPhoneNumber memformat nomor telepon ke format E.164 (standar internasional)
func FormatPhoneNumber(phoneNumber string, defaultRegion string) (string, error) {
	// Jika default region tidak diisi, gunakan ID (Indonesia)
	if defaultRegion == "" {
		defaultRegion = "ID"
	}

	// Parse nomor telepon
	parsedNumber, err := phonenumbers.Parse(phoneNumber, defaultRegion)
	if err != nil {
		return "", err
	}

	// Validasi nomor telepon
	if !phonenumbers.IsValidNumber(parsedNumber) {
		return "", err
	}

	// Format ke E.164
	formattedNumber := phonenumbers.Format(parsedNumber, phonenumbers.E164)

	return formattedNumber, nil
}

// IsValidPhoneNumber memeriksa apakah nomor telepon valid
func IsValidPhoneNumber(phoneNumber string, defaultRegion string) bool {
	// Jika default region tidak diisi, gunakan ID (Indonesia)
	if defaultRegion == "" {
		defaultRegion = "ID"
	}

	// Parse nomor telepon
	parsedNumber, err := phonenumbers.Parse(phoneNumber, defaultRegion)
	if err != nil {
		return false
	}

	// Validasi nomor telepon
	return phonenumbers.IsValidNumber(parsedNumber)
}

// GetCountryCodeForPhoneNumber mendapatkan kode negara dari nomor telepon
func GetCountryCodeForPhoneNumber(phoneNumber string, defaultRegion string) (string, error) {
	// Jika default region tidak diisi, gunakan ID (Indonesia)
	if defaultRegion == "" {
		defaultRegion = "ID"
	}

	// Parse nomor telepon
	parsedNumber, err := phonenumbers.Parse(phoneNumber, defaultRegion)
	if err != nil {
		return "", err
	}

	// Validasi nomor telepon
	if !phonenumbers.IsValidNumber(parsedNumber) {
		return "", err
	}

	// Dapatkan region dari nomor telepon
	region := phonenumbers.GetRegionCodeForNumber(parsedNumber)

	return region, nil
}

// IsPhoneNumber memeriksa apakah string berisi nomor telepon
func IsPhoneNumber(s string) bool {
	// Hapus semua karakter whitespace
	cleaned := strings.ReplaceAll(s, " ", "")

	// Cek apakah dimulai dengan "+" atau angka
	if !strings.HasPrefix(cleaned, "+") && !strings.HasPrefix(cleaned, "0") {
		return false
	}

	// Cek apakah sisanya angka
	for i := 1; i < len(cleaned); i++ {
		if cleaned[i] < '0' || cleaned[i] > '9' {
			// Izinkan karakter seperti - ( ) untuk format
			if cleaned[i] != '-' && cleaned[i] != '(' && cleaned[i] != ')' {
				return false
			}
		}
	}

	return true
}
