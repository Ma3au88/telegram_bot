package repository

import (
	"TestTelegramBot/internal/repository/postgres"
	"TestTelegramBot/pkg/cryptocompare"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
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

type Repository struct {
	Currency
	User
	Subscribers
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Currency:    postgres.NewRepository(db),
		User:        postgres.NewRepository(db),
		Subscribers: postgres.NewRepository(db),
	}
}
