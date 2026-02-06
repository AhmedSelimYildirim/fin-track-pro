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
	}

	for _, model := range modelsToCreate {
		_, err := DB.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}

	fmt.Println("PostgreSQL baglantisi basarili")
}
