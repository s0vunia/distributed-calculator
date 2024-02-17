package main

import (
	"context"
	"fmt"
	"log"
	"myproject/internal/config"
	config2 "myproject/project/internal/config"
	"myproject/project/internal/repositories/agent"
	"myproject/project/internal/repositories/expression"
	"myproject/project/internal/repositories/queue"
	"myproject/project/internal/repositories/subExpression"
	orchestrator2 "myproject/project/internal/services/orchestrator"
	"net/http"
	"time"
)

// Start инициализирует и запускает оркестратор
func Start() {
	dataSourceName := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		"postgres", "5432", "testttdb", "testttuser", "testttpass")
	expressionRepo, err := expression.NewPostgresRepository(dataSourceName)
	if err != nil {
		log.Printf("Failed to connect postgres: %v", err)
		return
	}
	subExpressionRepo, err := subExpression.NewPostgresRepository(dataSourceName)
	if err != nil {
		log.Printf("Failed to connect postgres: %v", err)
		return
	}
	agentRepo, err := agent.NewPostgresRepository(dataSourceName)
	if err != nil {
		log.Printf("Failed to connect agent postgres: %v", err)
		return
	}

	// Попытки реконнекта к rabbitmq (так как он не сразу отвечает)
	ticker := time.NewTicker(time.Second / 2)
	var queueRepo *queue.RabbitMQRepository
	for _ = range ticker.C {
		queueRepo, err = queue.NewRabbitMQRepository(config2.UrlRabbit, config.NameQueueWithTasks)
		if err != nil {
			log.Printf("Failed to start queue: %v", err)
			continue
		}
		break
	}

	ctx := context.Background()
	newOrchestrator := orchestrator2.NewOrchestrator(ctx, expressionRepo, subExpressionRepo, queueRepo, agentRepo)

	// Регистрация хендлеров
	http.HandleFunc("/expression/", orchestrator2.EndpointExpression(newOrchestrator))
	http.HandleFunc("/expression", orchestrator2.EndpointExpression(newOrchestrator))
	http.HandleFunc("/expressions", orchestrator2.GetExpressions(newOrchestrator))
	http.HandleFunc("/sub_expressions", orchestrator2.GetSubExpressions(newOrchestrator))
	http.HandleFunc("/agents", orchestrator2.GetAgents(newOrchestrator))

	port := ":8080"
	log.Printf("Starting server on %s", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func main() {
	Start()
}
