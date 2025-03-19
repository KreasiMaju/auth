package auth

// Version adalah versi saat ini dari package auth.
// Diperbarui saat build menggunakan ldflags.
var Version = "dev"

// GetVersion mengembalikan versi saat ini dari package auth.
func GetVersion() string {
	return Version
}
