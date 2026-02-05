package repository

import (
	"fmt"
	"log"

	"fin-track-pro/internal/config"
	"fin-track-pro/internal/core/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Veritabanına bağlanırken hata oluştu: ", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Asset{})
	if err != nil {
		log.Fatal("❌ Tablolar oluşturulurken hata: ", err)
	}

	fmt.Println("✅ Veritabanı bağlantısı sağlandı ve tüm tablolar senkronize edildi!")
	DB = db
}
