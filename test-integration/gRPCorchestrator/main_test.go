package gRPCorchestrator

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"os"
	"os/exec"
	"testing"
	"time"
)

type exprRes struct {
	timeout time.Duration
	res     float32
}

// ReconnectDelay - начальная задержка перед попыткой переподключения.
const ReconnectDelay = 10 * time.Second

// MaxReconnectDelay - максимальная задержка перед попыткой переподключения.
const MaxReconnectDelay = 30 * time.Second

// ReconnectBackoff - коэффициент увеличения задержки.
const ReconnectBackoff = 2

// CountTryReconnect - количество попыток переподключения.
const CountTryReconnect = 7

func TestMain(m *testing.M) {
	agents := 3
	cmd := exec.Command("make", "up-for-test-integration", fmt.Sprintf("AGENTS=%d", agents))
	cmd.Dir = "../../"
	cmd.Run()

	conn, _ := grpc.Dial("localhost:44044", grpc.WithInsecure())

	healthClient := grpc_health_v1.NewHealthClient(conn)

	// Проверяем здоровье сервера
	delay := ReconnectDelay
	countTries := 0
	for {
		countTries++
		if countTries == CountTryReconnect {
			fmt.Println("gRPC-Server is unreachable ")
			break
		}
		resp, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
		if err == nil && resp.Status == grpc_health_v1.HealthCheckResponse_SERVING {
			fmt.Println("gRPC-Server is healthy")
			break
		}

		fmt.Printf("Try to connect grpc-Server...\n")
		time.Sleep(delay)
		delay *= ReconnectBackoff
		if delay > MaxReconnectDelay {
			delay = MaxReconnectDelay
		}
	}

	code := m.Run()

	cmd = exec.Command("make", "down", fmt.Sprintf("AGENTS=%d", agents))
	cmd.Dir = "../../"
	cmd.Run()

	os.Exit(code)
}
