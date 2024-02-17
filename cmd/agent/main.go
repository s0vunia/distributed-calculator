package main

import (
	"fmt"
	"myproject/internal/config"
	config2 "myproject/project/internal/config"
	"myproject/project/internal/repositories/queue"
	"myproject/project/internal/services/agent"
)

// Start инициализирует и запускате агента
func Start() {
	queueRepo, err := queue.NewRabbitMQRepository(config2.UrlRabbit, config.NameQueueWithTasks)
	if err != nil {
		fmt.Printf("Failed to start queue: %v", err)
		return
	}
	demonQueueRepo, err := queue.NewRabbitMQRepository(config2.UrlRabbit, config.NameQueueWithTasks)
	if err != nil {
		fmt.Printf("Failed to start queue: %v", err)
		return
	}
	a := agent.NewAgent(queueRepo, demonQueueRepo)
	a.Start()
}

func main() {
	Start()
}
