package config

// Config adalah struktur konfigurasi utama untuk auth
type Config struct {
	Database  Database  `json:"database"`
	Providers Providers `json:"providers"`
	JWT       JWT       `json:"jwt"`
	Session   Session   `json:"session"`
	OTP       OTP       `json:"otp"`
}

// Database adalah konfigurasi untuk koneksi database
type Database struct {
	Type        string `json:"type"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Database    string `json:"database"`
	AutoMigrate bool   `json:"auto_migrate"`
}

// Providers berisi konfigurasi untuk berbagai provider OAuth
type Providers struct {
	Google   OAuth `json:"google"`
	Twitter  OAuth `json:"twitter"`
	GitHub   OAuth `json:"github"`
	Facebook OAuth `json:"facebook"`
	Local    bool  `json:"local"`    // Aktifkan autentikasi lokal (email/password)
	OTPAuth  bool  `json:"otp_auth"` // Aktifkan autentikasi OTP
}

// OAuth berisi konfigurasi untuk provider OAuth
type OAuth struct {
	Enabled      bool     `json:"enabled"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	CallbackURL  string   `json:"callback_url"`
	Scopes       []string `json:"scopes"`
}

// JWT berisi konfigurasi untuk token JWT
type JWT struct {
	Secret    string `json:"secret"`
	ExpiresIn int64  `json:"expires_in"` // dalam detik
}

// Session berisi konfigurasi untuk pengelolaan sesi
type Session struct {
	Secret    string `json:"secret"`
	ExpiresIn int64  `json:"expires_in"` // dalam detik
}

// OTP berisi konfigurasi untuk One-Time Password
type OTP struct {
	Enabled     bool        `json:"enabled"`
	DefaultType string      `json:"default_type"` // "sms", "email", "whatsapp"
	Length      int         `json:"length"`       // jumlah digit OTP
	ExpiresIn   int64       `json:"expires_in"`   // dalam detik
	SMS         OTPProvider `json:"sms"`
	Email       OTPProvider `json:"email"`
	WhatsApp    OTPProvider `json:"whatsapp"`
}

// OTPProvider berisi konfigurasi untuk provider OTP
type OTPProvider struct {
	Enabled  bool   `json:"enabled"`
	APIKey   string `json:"api_key"`
	Sender   string `json:"sender"`
	Template string `json:"template"`
}
