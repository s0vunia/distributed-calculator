package agent

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"myproject/internal/config"
	"myproject/internal/models"
	"myproject/internal/repositories/queue"
	"time"
)

type IAgent interface {
	// Start запускает агента
	Start()
	// CalculateExpression считает subexpression
	CalculateExpression(task *models.SubExpression)
	// StartHeartbeats отправка heartbeats
	StartHeartbeats()
}

type Agent struct {
	id                   string
	queueRepository      queue.Queue
	demonQueueRepository queue.Queue
}

func NewAgent(queueRepo, demonQueueRepo queue.Queue) *Agent {
	id := uuid.NewString()
	return &Agent{
		id:                   id,
		queueRepository:      queueRepo,
		demonQueueRepository: demonQueueRepo,
	}
}

func (a *Agent) Start() {
	// соединение с очередью subexpressions
	err := a.queueRepository.Connect(config.NameQueueWithTasks)
	if err != nil {
		log.Fatalf("Failed to connect to queue: %v", err)
	}
	defer a.queueRepository.Close()

	tasks, err := a.queueRepository.Consume()
	if err != nil {
		log.Fatalf("Failed to consume tasks from queue: %v", err)
	}

	// начинаем посылать heartbeats
	go a.StartHeartbeats()

	// обработка subexpressions из очереди
	for task := range tasks {
		expressionStruct := &models.SubExpression{}
		_ = json.Unmarshal(task, expressionStruct)

		// формирование ответа оркестратору и том, что взяли subexpression на обработку
		idAgent, _ := uuid.Parse(a.id)
		rpcAnswer := models.RPCAnswer{
			IdSubExpression: expressionStruct.Id,
			IdAgent:         idAgent,
		}
		err := a.queueRepository.Connect(config.NameQueueWithRPC)
		if err != nil {
			log.Printf("error connect to rpc queue")
		}
		rpcJson, err := json.Marshal(rpcAnswer)
		if err != nil {
			log.Printf("error unmarshal rpc")
		}
		err = a.queueRepository.Publish(rpcJson)
		if err != nil {
			log.Printf("error publish rpc")
		}
		a.queueRepository.Close()

		// подсчет subexpression
		a.CalculateExpression(expressionStruct)
	}
}

func (a *Agent) CalculateExpression(task *models.SubExpression) {
	result, err := Calculate(task)
	if err != nil {
		task.Error = true
	}
	task.Result = result

	err = a.queueRepository.Connect(config.NameQueueWithFinishedTasks)
	if err != nil {
		log.Fatalf("Failed to connect to queue: %v", err)
	}
	defer a.queueRepository.Close()

	expressionJson, err := json.Marshal(task)
	err = a.queueRepository.Publish(expressionJson)
	if err != nil {
		log.Printf("Failed to publish finished task to queue: %v", err)
	}
}

func (a *Agent) StartHeartbeats() {
	// Открываем соединение один раз, а не на каждую итерацию
	err := a.demonQueueRepository.Connect(config.NameQueueWithHeartbeats)
	if err != nil {
		log.Fatalf("Failed to connect to queue: %v\n", err)
	}
	defer a.demonQueueRepository.Close() // Закрыть соединение при завершении функции

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop() // Остановить тикер, когда функция завершится

	for range ticker.C {
		agent := &models.Agent{
			Id: a.id,
		}
		agentJson, err := json.Marshal(agent)
		if err != nil {
			log.Printf("Failed to encode agent: %v\n", err)
			continue // Продолжить цикл, если кодирование не удалось
		}

		err = a.demonQueueRepository.Publish(agentJson)
		if err != nil {
			log.Printf("Failed to publish task to queue: %v", err)
		}
	}
}
