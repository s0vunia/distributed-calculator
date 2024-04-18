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
	id                         string
	expressionQueueRepository  queue.Repository
	calculationQueueRepository queue.Repository
	heartbeatQueueRepository   queue.Repository
	rpcQueueRepository         queue.Repository
	calculationTimeouts        config.CalculationTimeoutsConfig
}

func NewAgent(expressionQueueRepo, calculationQueueRepo, heartbeatQueueRepo, rpcQueueRepo queue.Repository,
	timeouts config.CalculationTimeoutsConfig) *Agent {
	id := uuid.NewString()
	return &Agent{
		id:                         id,
		expressionQueueRepository:  expressionQueueRepo,
		calculationQueueRepository: calculationQueueRepo,
		heartbeatQueueRepository:   heartbeatQueueRepo,
		rpcQueueRepository:         rpcQueueRepo,
		calculationTimeouts:        timeouts,
	}
}

func (a *Agent) Start() {
	// соединение с очередью subexpressions
	err := a.expressionQueueRepository.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to queue: %v", err)
	}
	defer a.expressionQueueRepository.Close()

	tasks, err := a.expressionQueueRepository.Consume()
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
		err := a.rpcQueueRepository.Connect()
		if err != nil {
			log.Printf("error connect to rpc queue")
		}
		rpcJson, err := json.Marshal(rpcAnswer)
		if err != nil {
			log.Printf("error unmarshal rpc")
		}
		err = a.rpcQueueRepository.Publish(rpcJson)
		if err != nil {
			log.Printf("error publish rpc")
		}
		a.rpcQueueRepository.Close()

		// подсчет subexpression
		a.CalculateExpression(expressionStruct)
	}
}

func (a *Agent) CalculateExpression(task *models.SubExpression) {
	result, err := Calculate(task, a.calculationTimeouts)
	if err != nil {
		task.Error = true
	}
	task.Result = result

	err = a.calculationQueueRepository.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to queue: %v", err)
	}
	defer a.calculationQueueRepository.Close()

	expressionJson, err := json.Marshal(task)
	err = a.calculationQueueRepository.Publish(expressionJson)
	if err != nil {
		log.Printf("Failed to publish finished task to queue: %v", err)
	}
}

func (a *Agent) StartHeartbeats() {
	// Открываем соединение один раз, а не на каждую итерацию
	err := a.heartbeatQueueRepository.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to queue: %v\n", err)
	}
	defer a.heartbeatQueueRepository.Close() // Закрыть соединение при завершении функции

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

		err = a.heartbeatQueueRepository.Publish(agentJson)
		if err != nil {
			log.Printf("Failed to publish task to queue: %v", err)
		}
	}
}
