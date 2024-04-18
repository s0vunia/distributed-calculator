package gRPCorchestrator

import (
	"context"
	"fmt"
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
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestGRPCServiceParallelCalculation(t *testing.T) {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)
	cfg := config.MustLoadPath("../../config/local_tests.yaml")

	agents := 2
	cmd := exec.Command("make", "up", fmt.Sprintf("AGENTS=%d", agents))
	cmd.Dir = "../"
	err = cmd.Run()
	//assert.NoError(t, err)

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

	// Время ожидания уменьшено, учитывая, что агенты будут выполнять вычисления параллельно
	expressions := map[string]exprRes{
		"(2*5)/(1+1)": {
			timeout: max(cfg.CalculationTimeouts.TimeCalculateMult, cfg.CalculationTimeouts.TimeCalculatePlus) + cfg.CalculationTimeouts.TimeCalculateDivide + time.Second,
			res:     5,
		},
		"(5+3)*(4+3)": {
			timeout: cfg.CalculationTimeouts.TimeCalculatePlus + cfg.CalculationTimeouts.TimeCalculateMult + time.Second*1,
			res:     56,
		},
		"(456+8462)*(59347-1357)+278": {
			timeout: max(cfg.CalculationTimeouts.TimeCalculatePlus, cfg.CalculationTimeouts.TimeCalculateMinus) + cfg.CalculationTimeouts.TimeCalculateMult + cfg.CalculationTimeouts.TimeCalculatePlus + time.Second*1,
			res:     517155098,
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
		assert.Equal(t, expr.res, getExpressionResponse.Result)
	}
}
