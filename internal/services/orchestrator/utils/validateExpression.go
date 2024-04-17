package orchestratorutils

import (
	"go/parser"
	"log"
)

// ValidateExpression валидация выражения с помощью сторонней библиотеки
func ValidateExpression(expression string) bool {
	_, err := parser.ParseExpr(expression)
	if err != nil {
		log.Printf("parse error")
		return false
	}
	return true
}
