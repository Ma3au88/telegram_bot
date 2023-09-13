package service

import (
	"TestTelegramBot/internal/repository"
	"TestTelegramBot/internal/repository/postgres"
	"context"
)

type SubscribersService struct {
	repo repository.Subscribers
}

func NewSubscribersService(repo repository.Subscribers) *SubscribersService {
	return &SubscribersService{repo: repo}
}

func (ds *SubscribersService) StartAuto(ctx context.Context, delivery postgres.Subscribers) error {
	return ds.repo.StartAuto(ctx, delivery)
}

func (ds *SubscribersService) StopAuto(ctx context.Context, chatID int64) error {
	return ds.repo.StopAuto(ctx, chatID)
}

func (ds *SubscribersService) AutoUpdateUser(ctx context.Context) ([]postgres.Subscribers, error) {
	return ds.repo.AutoUpdateUser(ctx)
}

func (ds *SubscribersService) UpdateLastMessageTime(ctx context.Context, chatID int64) error {
	return ds.repo.UpdateLastMessageTime(ctx, chatID)
}
