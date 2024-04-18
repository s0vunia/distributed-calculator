package main

import (
	log "github.com/sirupsen/logrus"
	"myproject/internal/config"
	"myproject/internal/repositories/queue"
	"myproject/internal/services/agent"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Add this line for logging filename and line number!
	log.SetReportCaller(true)

	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

// Start инициализирует и запускате агента
func Start() {
	cfg := config.MustLoad()
	expressionsQueueRepo, err := queue.NewRabbitMQRepository(cfg.UrlRabbit, cfg.Queue.NameQueueWithTasks)
	if err != nil {
		log.Fatalf("Failed to start queue: %v", err)
		return
	}

	calculationQueueRepo, err := queue.NewRabbitMQRepository(cfg.UrlRabbit, cfg.Queue.NameQueueWithFinishedTasks)
	if err != nil {
		log.Fatalf("Failed to start queue: %v", err)
		return
	}
	heartbeatQueueRepo, err := queue.NewRabbitMQRepository(cfg.UrlRabbit, cfg.Queue.NameQueueWithHeartbeats)
	if err != nil {
		log.Fatalf("Failed to start queue: %v", err)
		return
	}
	rpcQueueRepo, err := queue.NewRabbitMQRepository(cfg.UrlRabbit, cfg.Queue.NameQueueWithRPC)
	if err != nil {
		log.Fatalf("Failed to start queue: %v", err)
		return
	}
	a := agent.NewAgent(expressionsQueueRepo, calculationQueueRepo, heartbeatQueueRepo, rpcQueueRepo, cfg.CalculationTimeouts)
	a.Start()
}

func main() {
	Start()
}
