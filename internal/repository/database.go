package repository

import (
	"fmt"
	"log"

	"github.com/AhmedSelimYildirim/fin-track-pro/internal/config"      // Config paketini buraya dahil et
	"github.com/AhmedSelimYildirim/fin-track-pro/internal/core/models" // Modellerini dahil et
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// 1. Önce yapılandırma (config) bilgilerini yüklüyoruz
	cfg := config.LoadConfig()

	// 2. Bağlantı dizesini (DSN) config'den gelen bilgilerle oluşturuyoruz
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	// 3. Veritabanı bağlantısını açıyoruz
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Veritabanına bağlanırken hata oluştu: ", err)
	}

	// 4. Modelleri (User ve Asset) veritabanına tablo olarak işliyoruz (Auto-Migration)
	err = db.AutoMigrate(&models.User{}, &models.Asset{})
	if err != nil {
		log.Fatal("❌ Tablolar oluşturulurken hata: ", err)
	}

	fmt.Println("✅ Veritabanı bağlantısı sağlandı ve tüm tablolar senkronize edildi!")
	DB = db
}
