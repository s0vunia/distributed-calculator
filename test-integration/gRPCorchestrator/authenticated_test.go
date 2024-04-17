package gRPCorchestrator

import (
	"context"
	"github.com/google/uuid"
	protos "github.com/s0vunia/protos/gen/go/auth"
	"github.com/s0vunia/protos/gen/go/orchestrator"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"myproject/internal/config"
	"testing"
	"time"
)

func TestGRPCServiceAuthenticated(t *testing.T) {
	conn, err := grpc.Dial("localhost:44044", grpc.WithInsecure())
	defer conn.Close()
	assert.NoError(t, err)
	client := orchestrator.NewOrchestratorClient(conn)
	// authenticated
	login, password := "abcdeeshka", "hahahaokey"
	authClient := protos.NewAuthClient(conn)
	registerResponse, err := authClient.Register(context.Background(), &protos.RegisterRequest{
		Login:    login,
		Password: password,
	})
	st, ok := status.FromError(err)
	if ok {
		assert.Equal(t, st.Code(), codes.AlreadyExists)
	}
	log.Printf("%v", registerResponse)

	loginResponse, err := authClient.Login(context.Background(), &protos.LoginRequest{
		Login:    login,
		Password: password,
		AppId:    1,
	})
	assert.NoError(t, err)
	token := loginResponse.Token

	md := metadata.New(map[string]string{
		"authorization": token,
	})

	// Создание контекста с метаданными
	expressions := map[string]exprRes{
		"2+2*2": {
			timeout: config.TimeCalculatePlus + config.TimeCalculateMult + time.Second*2,
			res:     6,
		},
		"(2+2)*2": {
			timeout: config.TimeCalculatePlus + config.TimeCalculateMult + time.Second*2,
			res:     8,
		},
		"6*6*5": {
			timeout: 3*config.TimeCalculateMult + time.Second*2,
			res:     180,
		},
		"(380-54)/2": {
			timeout: config.TimeCalculateDivide + config.TimeCalculateMinus + time.Second*2,
			res:     163,
		},
	}

	for key, expr := range expressions {
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		createExpressionResponse, err := client.CreateExpression(ctx, &orchestrator.CreateExpressionRequest{
			IdempotencyKey: uuid.New().String(),
			Expression:     key,
		})
		assert.NoError(t, err)
		log.Printf("%v", createExpressionResponse)
		tick := time.NewTicker(expr.timeout)
		<-tick.C
		getExpressionResponse, err := client.GetExpression(ctx, &orchestrator.GetExpressionRequest{
			ExpressionId: createExpressionResponse.ExpressionId,
		})
		log.Printf("%v", getExpressionResponse)
		assert.Equal(t, getExpressionResponse.Result, expr.res)
	}
}
