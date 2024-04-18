package gRPCorchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

type exprRes struct {
	timeout time.Duration
	res     float32
}

func TestMain(m *testing.M) {
	agents := 3
	cmd := exec.Command("make", "up", fmt.Sprintf("AGENTS=%d", agents))
	cmd.Dir = "../../"
	cmd.Run()
	os.Exit(m.Run())
}
