package config

import "time"

// Все конфигурационные данные
var (
	UrlRabbit                  = "amqp://guest:guest@rabbitmq/"
	NameQueueWithTasks         = "tasks"
	NameQueueWithRPC           = "tasks_rpc"
	NameQueueWithFinishedTasks = "finished_tasks"
	NameQueueWithHeartbeats    = "heartbeats"
	CountOfAgents              = 1
	CountOfGorutinsInAgent     = 10
	TimeCalculatePlus          = time.Second * 10
	TimeCalculateMinus         = time.Second * 10
	TimeCalculateMult          = time.Second * 10
	TimeCalculateDivide        = time.Second * 10
	RetrySubExpressionTimout   = time.Second * 40
)
