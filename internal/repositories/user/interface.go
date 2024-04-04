package user

import (
	"context"
	"myproject/internal/models"
)

type Repository interface {
	// Create создает пользователя
	Create(ctx context.Context, login string, passHash []byte) (int64, error)
	// Get возвращает пользователя по id
	Get(ctx context.Context, login string) (models.User, error)
}
