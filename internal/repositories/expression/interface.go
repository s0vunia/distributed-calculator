package expression

import (
	"context"
	"github.com/google/uuid"
	"myproject/internal/models"
)

type Repository interface {
	// CreateExpression создает expression
	CreateExpression(ctx context.Context, value string, idempotentKey string) (*models.Expression, error)
	// GetExpressions возвращает список expression
	GetExpressions(context.Context) ([]*models.Expression, error)
	// GetExpressionById возвращает expression по id
	GetExpressionById(context.Context, string) (*models.Expression, error)
	// GetExpressionByKey возвращает expression по ключу идемпотентности
	GetExpressionByKey(context.Context, string) (*models.Expression, error)
	// UpdateExpression обновляет expression
	UpdateExpression(context.Context, *models.Expression) error
	// UpdateExpressionById обновляет результат expression по ID
	UpdateExpressionById(ctx context.Context, id uuid.UUID, result float64) error
	// DeleteExpressionById удаляет expression по ID
	DeleteExpressionById(ctx context.Context, id uuid.UUID) error
	// UpdateState обновляет статус expression по ID
	UpdateState(ctx context.Context, id string, state models.ExpressionState) error
}
