package orchestrator

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

type Stack []string

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(str string) {
	*s = append(*s, str)
}

func (s *Stack) Pop() string {
	if s.IsEmpty() {
		return ""
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element
	}
}

func (s *Stack) Peek() string {
	if s.IsEmpty() {
		return ""
	}
	return (*s)[len(*s)-1]
}

func precedence(op string) int {
	if op == "+" || op == "-" {
		return 1
	}
	if op == "*" || op == "/" {
		return 2
	}
	return 0
}

// infixToPostfix преобразование инфиксной записи в постфиксную
func infixToPostfix(expression string) string {
	var result bytes.Buffer
	var stack Stack

	// Создаем регулярное выражение, которое будет использоваться для разделения чисел и операторов
	re := regexp.MustCompile(`\d+|[+/*()-]`)
	tokens := re.FindAllString(expression, -1)

	for _, token := range tokens {
		// Если это число, добавляем его в результат
		if _, err := strconv.Atoi(token); err == nil {
			result.WriteString(token + " ")
		} else {
			// Обработка скобок и операторов
			switch token {
			case "(":
				stack.Push(token)
			case ")":
				for !stack.IsEmpty() && stack.Peek() != "(" {
					result.WriteString(stack.Pop() + " ")
				}
				stack.Pop() // Удаляем открывающую скобку
			default:
				for !stack.IsEmpty() && precedence(token) <= precedence(stack.Peek()) {
					result.WriteString(stack.Pop() + " ")
				}
				stack.Push(token)
			}
		}
	}

	for !stack.IsEmpty() {
		result.WriteString(stack.Pop() + " ")
	}

	return strings.TrimSpace(result.String())
}
