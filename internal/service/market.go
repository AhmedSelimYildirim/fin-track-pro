package service

import (
	"encoding/json"
	"fin-track-pro/internal/config"
	"fin-track-pro/internal/repository"
	"fmt"
	"net/http"
)

type MarketService struct {
	cfg        *config.Config
	marketRepo *repository.MarketRepository
}

func NewMarketService(cfg *config.Config, repo *repository.MarketRepository) *MarketService {
	return &MarketService{
		cfg:        cfg,
		marketRepo: repo,
	}
}

func (s *MarketService) GetCurrencyRates() (map[string]float64, error) {
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/TRY", s.cfg.ExchangeAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		ConversionRates map[string]float64 `json:"conversion_rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.ConversionRates, nil
}

func (s *MarketService) GetCryptoPrice(coinID string) (float64, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=try", coinID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-cg-demo-api-key", s.cfg.CoinGeckoKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data[coinID]["try"], nil
}

func (s *MarketService) GetMetalPrice(metalCode string) (float64, error) {
	url := fmt.Sprintf("https://www.goldapi.io/api/%s/TRY", metalCode)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-access-token", s.cfg.MetalsAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Price float64 `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Price, nil
}
