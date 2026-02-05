package repository

import (
	"context"
	"database/sql"
	"fin-track-pro/internal/config"
	"fin-track-pro/internal/core/models"
	"fmt"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

func ConnectDB() {
	cfg := config.LoadConfig()

	if cfg.DBUser == "" {
		log.Fatal("❌ Hata: .env dosyasi okunamadi veya DB_USER bos! Calisma dizinini (Working Directory) kontrol et .")
	}

	pgconn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort)),
		pgdriver.WithUser(cfg.DBUser),
		pgdriver.WithPassword(cfg.DBPassword),
		pgdriver.WithDatabase(cfg.DBName),
		pgdriver.WithInsecure(true),
	)

	sqldb := sql.OpenDB(pgconn)
	DB = bun.NewDB(sqldb, pgdialect.New())

	ctx := context.Background()

	_, err := DB.NewCreateTable().Model((*models.User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatal("❌ User tablosu olusturulamadi: ", err)
	}

	_, err = DB.NewCreateTable().Model((*models.Asset)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatal("❌ Asset tablosu olusturulamadi: ", err)
	}

	fmt.Println("✅ Bun motoru tikir tikir calisiyor ! 🚀")
}
