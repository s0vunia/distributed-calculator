package agent

import (
	"myproject/project/internal/models"
)

type Repository interface {
	// Create создает агента с id
	Create(id string) error
	// IsExists проверяет, существует ли агент с id
	IsExists(id string) (bool, error)
	// CreateIfNotExistsAndUpdateHeartbeat создает агента, если не создан, в противном случае - обновляет heartbeat
	CreateIfNotExistsAndUpdateHeartbeat(id string) error
	// GetAgents возвращает список всех агентов
	GetAgents() ([]*models.Agent, error)
}
