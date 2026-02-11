package service

import (
	"context"
	"encoding/json"
	"fin-track-pro/internal/infrastructure/config"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type MarketService struct {
	cfg *config.Config
	rdb *redis.Client
}

func NewMarketService(cfg *config.Config, rdb *redis.Client) *MarketService {
	return &MarketService{cfg: cfg, rdb: rdb}
}

func (s *MarketService) GetCurrencyRates() (map[string]float64, error) {
	ctx := context.Background()
	cached, err := s.rdb.Get(ctx, "rates:currency").Result()
	if err == nil {
		var rates map[string]float64
		if err := json.Unmarshal([]byte(cached), &rates); err == nil {
			return rates, nil
		}
	}

	url := fmt.Sprintf("https://api.currencybeacon.com/v1/latest?api_key=%s&base=USD", s.cfg.CurrencyBeaconKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	tryRate := data.Rates["TRY"]
	if tryRate <= 0 {
		return nil, fmt.Errorf("gecersiz kur verisi")
	}

	finalRates := make(map[string]float64)
	finalRates["USD"] = tryRate

	if eurRate, ok := data.Rates["EUR"]; ok && eurRate != 0 {
		finalRates["EUR"] = tryRate / eurRate
	} else {
		finalRates["EUR"] = tryRate * 1.08 // Fallback kabaca parite
	}

	cacheData, _ := json.Marshal(finalRates)
	s.rdb.Set(ctx, "rates:currency", cacheData, 6*time.Hour) // 24 saat cok uzun, 6 saat daha ideal
	return finalRates, nil
}

func (s *MarketService) GetMetalPrice(metalCode string) (float64, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("price:metal:%s", metalCode)
	cached, err := s.rdb.Get(ctx, cacheKey).Float64()
	if err == nil {
		return cached, nil
	}

	rates, err := s.GetCurrencyRates()
	if err != nil {
		return 0, err
	}

	symbol := "XAU"
	if metalCode == "SILVER" {
		symbol = "XAG"
	}

	url := fmt.Sprintf("https://api.currencybeacon.com/v1/latest?api_key=%s&base=USD&symbols=%s", s.cfg.CurrencyBeaconKey, symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	// Ounce to Gram conversion
	priceInUSD := 1 / data.Rates[symbol]
	gramPriceTRY := (priceInUSD * rates["USD"]) / 31.1035

	s.rdb.Set(ctx, cacheKey, gramPriceTRY, 1*time.Hour) // Altin fiyati saatlik guncellenmeli
	return gramPriceTRY, nil
}

func (s *MarketService) GetCryptoPrice(coinID string) (float64, error) {
	ctx := context.Background()
	cacheKey := "price:crypto:BTC"
	cached, err := s.rdb.Get(ctx, cacheKey).Float64()
	if err == nil {
		return cached, nil
	}

	url := "https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT"
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	priceUSD, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0, err
	}

	rates, err := s.GetCurrencyRates()
	if err != nil {
		return 0, err
	}

	priceTRY := priceUSD * rates["USD"]
	s.rdb.Set(ctx, cacheKey, priceTRY, 5*time.Minute) // Kripto cok oynak, 5 dk ideal
	return priceTRY, nil
}
