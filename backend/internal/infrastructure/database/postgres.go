package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"fin-track-pro/internal/infrastructure/config"
	"fin-track-pro/internal/model"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

func ConnectDB() {
	cfg := config.LoadConfig()

	if cfg.DBUser == "" {
		log.Fatal("DB_USER bos")
	}

	isInsecure := true
	if os.Getenv("RENDER") == "true" || os.Getenv("DB_SSLMODE") == "require" {
		isInsecure = false
	}

	pgconn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort)),
		pgdriver.WithUser(cfg.DBUser),
		pgdriver.WithPassword(cfg.DBPassword),
		pgdriver.WithDatabase(cfg.DBName),
		pgdriver.WithInsecure(isInsecure),
	)

	sqldb := sql.OpenDB(pgconn)
	DB = bun.NewDB(sqldb, pgdialect.New())

	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	modelsToCreate := []interface{}{
		(*model.User)(nil),
		(*model.Asset)(nil),
		(*model.Transaction)(nil),
		(*model.Reminder)(nil),
		(*model.MarketHistory)(nil),
	}

	for _, m := range modelsToCreate {
		_, err := DB.NewCreateTable().Model(m).IfNotExists().Exec(ctx)
		if err != nil {
			log.Printf("Tablo kontrol hatasi: %v", err)
		}
	}

	_, _ = DB.ExecContext(ctx, "ALTER TABLE assets ALTER COLUMN ayar SET DEFAULT 0")
	_, _ = DB.ExecContext(ctx, "UPDATE assets SET ayar = 0 WHERE ayar IS NULL")

	fmt.Println("✅ PostgreSQL baglantisi ve iliskisel tablolar hazir")
}
