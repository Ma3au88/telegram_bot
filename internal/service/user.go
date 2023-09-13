package service

import (
	"TestTelegramBot/internal/repository"
	"context"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (uts *UserService) SaveUser(ctx context.Context, chatID int64) error {
	return uts.repo.SaveUser(ctx, chatID)
}
func (uts *UserService) GetUser(ctx context.Context, chatID int64) (int64, error) {
	return uts.repo.GetUser(ctx, chatID)
}
