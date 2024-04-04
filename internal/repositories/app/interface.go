package app

import (
	"context"
	"myproject/internal/models"
)

type Repository interface {
	App(ctx context.Context, appID int) (models.App, error)
}
