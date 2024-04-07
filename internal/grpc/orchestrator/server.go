package orchestratorgrpc

import (
	"context"
	"errors"
	orchv1 "github.com/s0vunia/protos/gen/go/orchestrator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"myproject/internal/repositories"
	"myproject/internal/services/orchestrator"
)

type serverAPI struct {
	orchv1.UnimplementedOrchestratorServer
	orchestrator orchestrator.IOrchestrator
}

func Register(gRPCServer *grpc.Server, orchestrator orchestrator.IOrchestrator) {
	orchv1.RegisterOrchestratorServer(gRPCServer, &serverAPI{orchestrator: orchestrator})
}

func (s *serverAPI) CreateExpression(
	ctx context.Context,
	in *orchv1.CreateExpressionRequest,
) (*orchv1.CreateExpressionResponse, error) {
	if in.Expression == "" {

		return nil, status.Error(codes.InvalidArgument, "expression is required")
	}
	if in.IdempotencyKey == "" {

		return nil, status.Error(codes.InvalidArgument, "idempotencyKey is required")
	}

	var expressionId string
	expressionByKey, err := s.orchestrator.GetExpressionByKey(ctx, in.IdempotencyKey)
	if expressionByKey != nil {
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to create expression")
		}
		expressionId = expressionByKey.Id
	} else {
		err, expressionId = s.orchestrator.CreateExpression(ctx, in.Expression, in.IdempotencyKey)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to create expression")
		}
	}
	return &orchv1.CreateExpressionResponse{ExpressionId: expressionId}, nil
}

func (s *serverAPI) GetExpression(
	ctx context.Context,
	in *orchv1.GetExpressionRequest,
) (*orchv1.GetExpressionResponse, error) {
	if in.ExpressionId == "" {

		return nil, status.Error(codes.InvalidArgument, "expressionId is required")
	}

	expression, err := s.orchestrator.GetExpression(ctx, in.ExpressionId)
	if err != nil {
		if errors.Is(err, repositories.ErrExpressionNotFound) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to get expression")
	}

	return &orchv1.GetExpressionResponse{
		Result:         float32(expression.Result),
		ExpressionId:   expression.Id,
		IdempotencyKey: expression.IdempotencyKey,
		Value:          expression.Value,
		State:          string(expression.State),
	}, nil
}
