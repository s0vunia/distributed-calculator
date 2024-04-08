package orchestrator

import (
	"go/parser"
	"log"
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
