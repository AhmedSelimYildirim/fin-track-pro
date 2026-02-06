package service

import (
	"context"
	"encoding/json"
	"fin-track-pro/internal/infrastructure/config"
	"fmt"
	"net/http"
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
		json.Unmarshal([]byte(cached), &rates)
		return rates, nil
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
	json.NewDecoder(resp.Body).Decode(&data)

	tryRate := data.Rates["TRY"]
	finalRates := make(map[string]float64)
	finalRates["USD"] = tryRate
	if eurRate, ok := data.Rates["EUR"]; ok && eurRate != 0 {
		finalRates["EUR"] = tryRate / eurRate
	}

	cacheData, _ := json.Marshal(finalRates)
	s.rdb.Set(ctx, "rates:currency", cacheData, 7*24*time.Hour)
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
	usdToTry := rates["USD"]

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
	json.NewDecoder(resp.Body).Decode(&data)

	priceInUSD := 1 / data.Rates[symbol]
	gramPriceTRY := (priceInUSD * usdToTry) / 31.1035

	s.rdb.Set(ctx, cacheKey, gramPriceTRY, 7*24*time.Hour)
	return gramPriceTRY, nil
}

func (s *MarketService) GetCryptoPrice(coinID string) (float64, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("price:crypto:%s", coinID)
	cached, err := s.rdb.Get(ctx, cacheKey).Float64()
	if err == nil {
		return cached, nil
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=try", coinID)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data map[string]map[string]float64
	json.NewDecoder(resp.Body).Decode(&data)

	price := data[coinID]["try"]
	s.rdb.Set(ctx, cacheKey, price, 7*24*time.Hour)
	return price, nil
}

func (s *MarketService) GetHistoricalRate(date time.Time, base string, target string) (float64, error) {
	ctx := context.Background()
	dateStr := date.Format("2006-01-02")
	cacheKey := fmt.Sprintf("rate:hist:%s:%s:%s", base, target, dateStr)

	cached, err := s.rdb.Get(ctx, cacheKey).Float64()
	if err == nil {
		return cached, nil
	}

	url := fmt.Sprintf("https://api.frankfurter.app/%s?from=%s&to=%s", dateStr, base, target)
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

	rate := data.Rates[target]
	if rate != 0 {
		s.rdb.Set(ctx, cacheKey, rate, 0)
	}

	return rate, nil
}
