package tests

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"myproject/internal/models"
	"myproject/internal/repositories/agent"
	"myproject/internal/repositories/expression"
	expressionMock "myproject/internal/repositories/expression/mocks"
	"myproject/internal/repositories/queue"
	"myproject/internal/repositories/subExpression"
	subExpressionMock "myproject/internal/repositories/subExpression/mocks"
	"myproject/internal/services/orchestrator"
	"reflect"
	"testing"
)

func TestOrchestrator_CreateExpression(t *testing.T) {
	type fields struct {
		expressionRepository        expression.Repository
		subExpressionRepository     subExpression.Repository
		agentRepository             agent.Repository
		expressionsQueueRepository  queue.Repository
		calculationsQueueRepository queue.Repository
		heartbeatsQueueRepository   queue.Repository
		rpcQueueRepository          queue.Repository
	}
	type args struct {
		ctx            context.Context
		expression     string
		idempotencyKey string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   error
		want1  string
	}{
		{
			name: "ok",
			args: args{
				ctx:            context.Background(),
				expression:     "2+2*2",
				idempotencyKey: "abc",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expressionRepo := expressionMock.NewRepository(t)
			subexpressionRepo := subExpressionMock.NewRepository(t)

			retExpression := &models.Expression{
				Id:             "1234567",
				IdempotencyKey: tt.args.idempotencyKey,
				Value:          tt.args.expression,
				State:          models.ExpressionInProgress,
			}
			tt.want1 = "1234567"
			expressionRepo.
				On("CreateExpression", mock.Anything, tt.args.expression, tt.args.idempotencyKey).
				Return(retExpression, nil)
			subexpressionRepo.
				On("CreateSubExpression", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, subExpression *models.SubExpression) (*models.SubExpression, error) {
					subExpression.Id = uuid.New()
					return subExpression, nil
				})

			o := &orchestrator.Orchestrator{
				expressionRepository:    expressionRepo,
				subExpressionRepository: subexpressionRepo,
			}
			got, got1 := o.CreateExpression(tt.args.ctx, tt.args.expression, tt.args.idempotencyKey)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateExpression() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CreateExpression() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
func TestOrchestrator_CreateExpression_With_Error_Split_to_Tasks(t *testing.T) {
	type fields struct {
		expressionRepository        expression.Repository
		subExpressionRepository     subExpression.Repository
		agentRepository             agent.Repository
		expressionsQueueRepository  queue.Repository
		calculationsQueueRepository queue.Repository
		heartbeatsQueueRepository   queue.Repository
		rpcQueueRepository          queue.Repository
	}
	type args struct {
		ctx            context.Context
		expression     string
		idempotencyKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want1   string
	}{
		{
			name: "ok",
			args: args{
				ctx:            context.Background(),
				expression:     "2+2*2",
				idempotencyKey: "abc",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expressionRepo := expressionMock.NewRepository(t)
			subexpressionRepo := subExpressionMock.NewRepository(t)

			retExpression := &models.Expression{
				Id:             "1234567",
				IdempotencyKey: tt.args.idempotencyKey,
				Value:          tt.args.expression,
				State:          models.ExpressionInProgress,
			}
			tt.want1 = "1234567"
			expressionRepo.
				On("CreateExpression", mock.Anything, tt.args.expression, tt.args.idempotencyKey).
				Return(retExpression, nil)
			subexpressionRepo.
				On("CreateSubExpression", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, subExpression *models.SubExpression) (*models.SubExpression, error) {
					return nil, errors.New("err split sub expression")
				})
			exprId, _ := uuid.Parse(retExpression.Id)
			subexpressionRepo.
				On("DeleteSubExpressionsByExpressionId", mock.Anything, exprId).
				Return(nil)
			expressionRepo.
				On("DeleteExpressionById", mock.Anything, exprId).
				Return(nil)

			o := &orchestrator.Orchestrator{
				expressionRepository:    expressionRepo,
				subExpressionRepository: subexpressionRepo,
			}
			got, _ := o.CreateExpression(tt.args.ctx, tt.args.expression, tt.args.idempotencyKey)
			if got != nil && tt.wantErr == false {
				t.Errorf("CreateExpression() got = %v, want %v", got, tt.wantErr)
			}
		})
	}
}
