package postgres

import (
	"TestTelegramBot/pkg/cryptocompare"
	"context"
	"fmt"
	"time"
)

func (r *Repository) SaveCurrency(ctx context.Context, coins []cryptocompare.Price) error {
	var id int64
	query := `INSERT INTO crypto_prices (coin, value, date) VALUES ($1, $2, $3) RETURNING id`

	for _, coin := range coins {
		row := r.db.QueryRow(ctx, query, coin.Coin, coin.Value, time.Now())
		if err := row.Scan(&id); err != nil {
			return fmt.Errorf("database scan: %w", err)
		}
	}

	return nil
}

func (r *Repository) GetCurrency(ctx context.Context, coins []string) ([]CryptoPrices, error) {
	query := `SELECT * FROM crypto_prices WHERE coin = $1 ORDER BY id DESC LIMIT 1;`

	var price CryptoPrices
	prices := make([]CryptoPrices, 0, len(coins))

	for _, coin := range coins {
		row, err := r.db.Query(ctx, query, coin)
		if err != nil {
			return nil, fmt.Errorf("sending query to database: %w", err)
		}

		for row.Next() {
			err := row.Scan(&price.Id, &price.Coin, &price.Value, &price.Date)
			if err != nil {
				return nil, fmt.Errorf("database scan: %w", err)
			}

			prices = append(prices, price)
		}

		if err := row.Err(); err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
	}

	return prices, nil
}
