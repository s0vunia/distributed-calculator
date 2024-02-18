package main

import (
	"fmt"
	"log"
	"myproject/internal/config"
	"myproject/internal/repositories/queue"
	"myproject/internal/services/agent"
	"time"
)

// Start инициализирует и запускате агента
func Start() {
	ticker := time.NewTicker(time.Second / 2)
	var expressionsQueueRepo *queue.RabbitMQRepository
	var err error
	for _ = range ticker.C {
		expressionsQueueRepo, err = queue.NewRabbitMQRepository(config.UrlRabbit, config.NameQueueWithTasks)
		if err != nil {
			log.Printf("Failed to start queue: %v", err)
			continue
		}
		break
	}
	calculationQueueRepo, err := queue.NewRabbitMQRepository(config.UrlRabbit, config.NameQueueWithFinishedTasks)
	if err != nil {
		fmt.Printf("Failed to start queue: %v", err)
		return
	}
	heartbeatQueueRepo, err := queue.NewRabbitMQRepository(config.UrlRabbit, config.NameQueueWithHeartbeats)
	if err != nil {
		fmt.Printf("Failed to start queue: %v", err)
		return
	}
	rpcQueueRepo, err := queue.NewRabbitMQRepository(config.UrlRabbit, config.NameQueueWithRPC)
	if err != nil {
		fmt.Printf("Failed to start queue: %v", err)
		return
	}
	a := agent.NewAgent(expressionsQueueRepo, calculationQueueRepo, heartbeatQueueRepo, rpcQueueRepo)
	a.Start()
}

func main() {
	Start()
}
