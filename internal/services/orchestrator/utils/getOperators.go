package orchestratorutils

import (
	"myproject/internal/config"
	"myproject/internal/models"
	"time"
)

// GetOperators возвращает список операция
func GetOperators(timeouts config.CalculationTimeoutsConfig) []*models.Operator {
	operatorsMap := map[string]time.Duration{
		"+": timeouts.TimeCalculatePlus / time.Second,
		"-": timeouts.TimeCalculateMinus / time.Second,
		"*": timeouts.TimeCalculateMult / time.Second,
		"/": timeouts.TimeCalculateDivide / time.Second,
	}
	var operators []*models.Operator
	for key, value := range operatorsMap {
		operators = append(operators, &models.Operator{Op: key, Timeout: value})
	}
	return operators
}
