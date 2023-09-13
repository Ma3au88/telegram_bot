package service

import (
	"TestTelegramBot/internal/repository"
	"TestTelegramBot/internal/repository/postgres"
	"TestTelegramBot/pkg/cryptocompare"
	"context"
)

type CurrencyService struct {
	repo repository.Currency
}

func NewCurrencyService(repo repository.Currency) *CurrencyService {
	return &CurrencyService{repo: repo}
}

func (ct *CurrencyService) SaveCurrency(ctx context.Context, coins []cryptocompare.Price) error {
	return ct.repo.SaveCurrency(ctx, coins)
}

func (ct *CurrencyService) GetCurrency(ctx context.Context, coins []string) ([]postgres.CryptoPrices, error) {
	return ct.repo.GetCurrency(ctx, coins)
}
