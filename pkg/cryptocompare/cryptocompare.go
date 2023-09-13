package cryptocompare

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const URL = "https://min-api.cryptocompare.com/data/pricemulti?fsyms=%s&tsyms=RUR"

type Cryptocompare struct {
	URL    string
	Client *http.Client
}

func NewCryptocompare() *Cryptocompare {
	return &Cryptocompare{URL: URL, Client: &http.Client{Timeout: 5 * time.Second}}
}

type Price struct {
	Coin  string
	Value float64
}

func (c *Cryptocompare) GetPrice(ctx context.Context, coinName []string) ([]Price, error) {
	coins := strings.Join(coinName, ",")
	rec, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(c.URL, coins), nil)
	if err != nil {
		return nil, fmt.Errorf("creating query with context: %w", err)
	}

	resp, err := c.Client.Do(rec)
	if err != nil {
		return nil, fmt.Errorf("http request execution: %w", err)
	}
	defer resp.Body.Close()

	var prices struct {
		BTC map[string]float64 `json:"BTC"`
		ETH map[string]float64 `json:"ETH"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	var priceDuo []Price

	for _, coin := range coinName {
		switch coin {
		case "BTC":
			priceDuo = append(priceDuo, Price{Coin: "BTC", Value: prices.BTC["RUR"]})
		case "ETH":
			priceDuo = append(priceDuo, Price{Coin: "ETH", Value: prices.ETH["RUR"]})
		}
	}

	return priceDuo, nil
}
