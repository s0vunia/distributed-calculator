package gRPCorchestrator

import (
	"context"
	"fmt"
	protos "github.com/s0vunia/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func TestMain(m *testing.M) {
	agents := 3
	cmd := exec.Command("make", "up", fmt.Sprintf("AGENTS=%d", agents))
	cmd.Dir = "../../"
	cmd.Run()

	conn, _ := grpc.Dial("localhost:44044", grpc.WithInsecure())

	// test connection
	authClient := protos.NewAuthClient(conn)
	_, err := authClient.Register(context.Background(), &protos.RegisterRequest{
		Login:    "",
		Password: "",
	})
	st, ok := status.FromError(err)
	delay := ReconnectDelay
	for ok && st.Code() == codes.Unavailable {
		fmt.Printf("Reconnect...\n")
		time.Sleep(delay)
		delay *= ReconnectBackoff
		if delay > MaxReconnectDelay {
			delay = MaxReconnectDelay
		}
		_, err := authClient.Register(context.Background(), &protos.RegisterRequest{
			Login:    "",
			Password: "",
		})
		st, ok = status.FromError(err)
	}

	os.Exit(m.Run())
}
