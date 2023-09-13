package service

import (
	"TestTelegramBot/internal/repository"
	"TestTelegramBot/internal/repository/postgres"
	"TestTelegramBot/pkg/cryptocompare"
	"context"
)

type Currency interface {
	SaveCurrency(ctx context.Context, coins []cryptocompare.Price) error
	GetCurrency(ctx context.Context, coins []string) ([]postgres.CryptoPrices, error)
}

type User interface {
	SaveUser(ctx context.Context, chatID int64) error
	GetUser(ctx context.Context, chatID int64) (int64, error)
}

type Subscribers interface {
	StartAuto(ctx context.Context, delivery postgres.Subscribers) error
	StopAuto(ctx context.Context, chatID int64) error
	AutoUpdateUser(ctx context.Context) ([]postgres.Subscribers, error)
	UpdateLastMessageTime(ctx context.Context, chatID int64) error
}

type GetPrices interface {
	GetPrice(ctx context.Context, coinName []string) ([]cryptocompare.Price, error)
}

type Service struct {
	Currency
	User
	Subscribers
	GetPrices
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Currency:    NewCurrencyService(repository),
		User:        NewUserService(repository),
		Subscribers: NewSubscribersService(repository),
		GetPrices:   cryptocompare.NewCryptocompare(),
	}
}
