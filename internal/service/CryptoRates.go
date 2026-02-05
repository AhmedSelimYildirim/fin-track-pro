package service

import (
	"encoding/json"
	"net/http"
)

type CryptoRates map[string]map[string]float64

func GetCryptoRates() (CryptoRates, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum,binancecoin&vs_currencies=usd,try"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data CryptoRates
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
