package orchestratorgrpc

import (
	"context"
	"errors"
	authv1 "github.com/s0vunia/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"myproject/internal/repositories"
	"myproject/internal/services/orchestrator"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
	orchestrator orchestrator.IOrchestrator
}

func Register(gRPCServer *grpc.Server, orchestrator orchestrator.IOrchestrator) {
	authv1.RegisterAuthServer(gRPCServer, &serverAPI{orchestrator: orchestrator})
}

func (s *serverAPI) CreateExpression(
	ctx context.Context,
	in *authv1.CreateExpressionRequest,
) (*authv1.CreateExpressionResponse, error) {
	if in.Expression == "" {

		return nil, status.Error(codes.InvalidArgument, "expression is required")
	}
	if in.IdempotencyKey == "" {

		return nil, status.Error(codes.InvalidArgument, "idempotencyKey is required")
	}
	err, expressionId := s.orchestrator.CreateExpression(ctx, in.Expression, in.IdempotencyKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create expression")
	}
	return &authv1.CreateExpressionResponse{ExpressionId: expressionId}, nil
}

func (s *serverAPI) GetExpression(
	ctx context.Context,
	in *authv1.GetExpressionRequest,
) (*authv1.GetExpressionResponse, error) {
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

	return &authv1.GetExpressionResponse{
		Result:         float32(expression.Result),
		ExpressionId:   expression.Id,
		IdempotencyKey: expression.IdempotencyKey,
		Value:          expression.Value,
		State:          string(expression.State),
	}, nil
}
