package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBPort            string
	DBURL             string
	JWTSecret         string
	ExchangeAPIKey    string
	CurrencyBeaconKey string
	RedisURL          string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		DBHost:            os.Getenv("DB_HOST"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		DBPort:            os.Getenv("DB_PORT"),
		DBURL:             os.Getenv("DB_URL"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		ExchangeAPIKey:    os.Getenv("EXCHANGE_RATE_API_KEY"),
		CurrencyBeaconKey: os.Getenv("CURRENCY_BEACON_KEY"),
		RedisURL:          os.Getenv("REDIS_URL"),
	}
}
