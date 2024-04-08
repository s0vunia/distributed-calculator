package orchestratorutils

import (
	"myproject/internal/config"
	"myproject/internal/models"
	"time"
)

// GetOperators возвращает список операция
func GetOperators() []*models.Operator {
	operatorsMap := map[string]time.Duration{
		"+": config.TimeCalculatePlus / time.Second,
		"-": config.TimeCalculateMinus / time.Second,
		"*": config.TimeCalculateMult / time.Second,
		"/": config.TimeCalculateDivide / time.Second,
	}
	var operators []*models.Operator
	for key, value := range operatorsMap {
		operators = append(operators, &models.Operator{Op: key, Timeout: value})
	}
	return operators
}
