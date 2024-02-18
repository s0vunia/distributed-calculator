package main

import (
	"context"
	"fmt"
	"log"
	"myproject/internal/config"
	"myproject/internal/repositories/agent"
	"myproject/internal/repositories/expression"
	"myproject/internal/repositories/queue"
	"myproject/internal/repositories/subExpression"
	"myproject/internal/services/orchestrator"
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
		queueRepo, err = queue.NewRabbitMQRepository(config.UrlRabbit, config.NameQueueWithTasks)
		if err != nil {
			log.Printf("Failed to start queue: %v", err)
			continue
		}
		break
	}

	ctx := context.Background()
	newOrchestrator := orchestrator.NewOrchestrator(ctx, expressionRepo, subExpressionRepo, queueRepo, agentRepo)
	// Регистрация хендлеров
	http.HandleFunc("/expression/", orchestrator.EndpointExpression(newOrchestrator))
	http.HandleFunc("/expression", orchestrator.EndpointExpression(newOrchestrator))
	http.HandleFunc("/expressions", orchestrator.GetExpressions(newOrchestrator))
	http.HandleFunc("/sub_expressions", orchestrator.GetSubExpressions(newOrchestrator))
	http.HandleFunc("/agents", orchestrator.GetAgents(newOrchestrator))
	http.HandleFunc("/operators", orchestrator.GetOperators)

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
