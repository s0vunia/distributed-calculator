package subExpression

import (
	"context"
	"github.com/google/uuid"
	"myproject/internal/models"
)

type Repository interface {
	// CreateSubExpression создает subexpression
	CreateSubExpression(ctx context.Context, subExpression *models.SubExpression) (*models.SubExpression, error)
	// GetSubExpressions вовзращает канал c subexpressions
	GetSubExpressions() chan *models.SubExpression
	// UpdateSubExpressions обновляет expression
	UpdateSubExpressions(ctx context.Context, expression *models.SubExpression) error
	// GetSubExpressionsList возвращает список subexpressions
	GetSubExpressionsList(ctx context.Context) ([]*models.SubExpression, error)
	// DeleteSubExpressionsByExpressionId удаляет subexpression по его ID
	DeleteSubExpressionsByExpressionId(ctx context.Context, expressionId uuid.UUID) error
	// UpdateSubExpressionAgent обновляет agent_id у subexpression
	UpdateSubExpressionAgent(ctx context.Context, idSubExpression, agentId uuid.UUID) error
	// DeleteSubExpressionById удаляет subexpression по его id
	DeleteSubExpressionById(ctx context.Context, id uuid.UUID) error
	// GetNotCalculatedSubExpressionsByAgentId удаляет неподсчитанные subexpression по agent_id
	GetNotCalculatedSubExpressionsByAgentId(ctx context.Context, agentId uuid.UUID) ([]*models.SubExpression, error)
	// ReplaceExpressionsIds меняет sub_expression1 и sub_expression2 с oldId на newId
	ReplaceExpressionsIds(ctx context.Context, oldId uuid.UUID, newId uuid.UUID) error
}
