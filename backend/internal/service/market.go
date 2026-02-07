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
	s.rdb.Set(ctx, "rates:currency", cacheData, 24*time.Hour)
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

	s.rdb.Set(ctx, cacheKey, gramPriceTRY, 12*time.Hour)
	return gramPriceTRY, nil
}

func (s *MarketService) GetCryptoPrice(coinID string) (float64, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("price:crypto:%s", coinID)
	cached, err := s.rdb.Get(ctx, cacheKey).Float64()
	if err == nil {
		return cached, nil
	}

	url := fmt.Sprintf("https://api.coincap.io/v2/assets/%s", coinID)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("crypto api error: %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			PriceUsd string `json:"priceUsd"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	priceUSD, err := strconv.ParseFloat(result.Data.PriceUsd, 64)
	if err != nil {
		return 0, err
	}

	rates, err := s.GetCurrencyRates()
	if err != nil {
		return 0, err
	}

	priceTRY := priceUSD * rates["USD"]
	s.rdb.Set(ctx, cacheKey, priceTRY, 10*time.Minute)
	return priceTRY, nil
}
