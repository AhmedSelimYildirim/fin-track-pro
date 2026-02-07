package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBPort            string
	JWTSecret         string
	ExchangeAPIKey    string
	CoinGeckoKey      string
	MetalsAPIKey      string
	CurrencyBeaconKey string
	RedisURL          string
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	return &Config{
		DBHost:            os.Getenv("DB_HOST"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		DBPort:            os.Getenv("DB_PORT"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		ExchangeAPIKey:    os.Getenv("EXCHANGE_RATE_API_KEY"),
		CoinGeckoKey:      os.Getenv("COINGECKO_API_KEY"),
		MetalsAPIKey:      os.Getenv("METALS_API_KEY"),
		CurrencyBeaconKey: os.Getenv("CURRENCY_BEACON_KEY"),
		RedisURL:          os.Getenv("REDIS_URL"),
	}
}
