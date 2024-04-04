package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"log/slog"
	"myproject/internal/app"
	"myproject/internal/config"
	"myproject/internal/repositories/agent"
	appRepo "myproject/internal/repositories/app"
	"myproject/internal/repositories/expression"
	"myproject/internal/repositories/queue"
	"myproject/internal/repositories/subExpression"
	"myproject/internal/repositories/user"
	"myproject/internal/services/auth"
	"myproject/internal/services/orchestrator"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Add this line for logging filename and line number!
	log.SetReportCaller(true)

	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

// Start инициализирует и запускает оркестратор
func Start() {
	dataSourceName := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		"postgres", "5432", "testttdb", "testttuser", "testttpass")
	expressionRepo, err := expression.NewPostgresRepository(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect postgres: %v", err)
		return
	}
	subExpressionRepo, err := subExpression.NewPostgresRepository(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect postgres: %v", err)
		return
	}
	agentRepo, err := agent.NewPostgresRepository(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect agent postgres: %v", err)
		return
	}

	expressionsQueueRepo, err := queue.NewRabbitMQRepository(config.UrlRabbit, config.NameQueueWithTasks)
	if err != nil {
		log.Fatalf("Failed to start queue: %v", err)
	}
	calculationsQueueRepository, err := queue.NewRabbitMQRepository(config.UrlRabbit, config.NameQueueWithFinishedTasks)
	if err != nil {
		log.Fatalf("Failed to start queue: %v", err)
	}
	heartbeatsQueueRepository, err := queue.NewRabbitMQRepository(config.UrlRabbit, config.NameQueueWithHeartbeats)
	if err != nil {
		log.Fatalf("Failed to start queue: %v", err)
	}
	rpcQueueRepository, err := queue.NewRabbitMQRepository(config.UrlRabbit, config.NameQueueWithRPC)
	if err != nil {
		log.Fatalf("Failed to start queue: %v", err)
	}
	userRepository, err := user.NewPostgresRepository(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to start queue: %v", err)
	}
	appRepository, err := appRepo.NewPostgresRepository(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to start queue: %v", err)
	}

	ctx := context.Background()
	logSlog := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	newOrchestrator := orchestrator.NewOrchestrator(ctx, expressionRepo, subExpressionRepo, expressionsQueueRepo,
		calculationsQueueRepository, heartbeatsQueueRepository, rpcQueueRepository, agentRepo)
	newAuth := auth.New(logSlog, userRepository, userRepository, appRepository, time.Hour*10)

	// Регистрация хендлеров
	application := app.New(logSlog, newOrchestrator, newAuth, 8080, 44044, time.Hour*10)
	go func() {
		application.ServerHTTP.MustRun()
	}()
	go func() {
		application.GRPCServer.MustRun()
	}()
	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.ServerHTTP.Stop()
	log.Info("Gracefully stopped")

}

func main() {
	Start()
}
