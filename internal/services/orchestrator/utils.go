package orchestrator

import (
	"go/parser"
	"log"
	"myproject/internal/config"
	"myproject/internal/models"
	"time"
)

// validateExpression валидация выражения с помощью сторонней библиотеки
func validateExpression(expression string) bool {
	_, err := parser.ParseExpr(expression)
	if err != nil {
		log.Printf("parse error")
		return false
	}
	return true
}

// getOperators возвращает список операция
func getOperators() []*models.Operator {
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
