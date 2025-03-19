package utils

import (
	"fmt"
	"log"

	"github.com/kreasimaju/auth/config"
	"github.com/kreasimaju/auth/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB menginisialisasi koneksi database berdasarkan konfigurasi
func InitDB(config config.Database) (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	switch config.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username, config.Password, config.Host, config.Port, config.Database)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
			config.Host, config.Username, config.Password, config.Database, config.Port)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Database), &gorm.Config{})
	default:
		return nil, fmt.Errorf("database type not supported: %s", config.Type)
	}

	if err != nil {
		return nil, err
	}

	DB = db

	// Auto migrate database jika diminta
	if config.AutoMigrate {
		err = MigrateDB(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// MigrateDB melakukan migrasi tabel-tabel yang diperlukan
func MigrateDB(db *gorm.DB) error {
	log.Println("Running automatic database migrations...")

	// Daftar model untuk migrasi
	err := db.AutoMigrate(
		&models.User{},
		&models.UserProvider{},
		&models.Session{},
		&models.Token{},
		&models.OTPCode{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}
