package orchestrator

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	"myproject/internal/models"
	"myproject/internal/repositories/subExpression"
	"regexp"
	"strconv"
	"strings"
)

// splitToSubtasks делать полное арифметическое выражение на подзадачи
func splitToSubtasks(ctx context.Context, expr *models.Expression, subExpressionRepo subExpression.Repository) (tasks []*models.SubExpression, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = errors.New("unknown panic in splitter subtasks")
			}
		}
	}()

	var stack []string
	uuidRegex := regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

	// функция создания subexpression
	getTempVar := func(operand1, operand2, element string, isLast bool) (*models.SubExpression, error) {
		uid, _ := uuid.Parse(expr.Id)
		subExpr := &models.SubExpression{
			ExpressionId: uid,
			IsLast:       isLast,
			Action:       element,
			Error:        false,
		}

		if uuidRegex.MatchString(operand1) {
			operand1Uid, err := uuid.Parse(operand1)
			if err != nil {
				log.Printf("error parse operand1")
				return nil, err
			}
			subExpr.SubExpressionId1 = uuid.NullUUID{
				UUID:  operand1Uid,
				Valid: true,
			}
		} else {
			el, err := strconv.ParseFloat(operand1, 64)
			if err != nil {
				return nil, err
			}
			subExpr.Val1 = el
		}
		if uuidRegex.MatchString(operand2) {
			operand2Uid, err := uuid.Parse(operand2)
			if err != nil {
				log.Printf("error parse operand2")
				return nil, err
			}
			subExpr.SubExpressionId2 = uuid.NullUUID{
				UUID:  operand2Uid,
				Valid: true,
			}
		} else {
			el, err := strconv.ParseFloat(operand2, 64)
			if err != nil {
				return nil, err
			}
			subExpr.Val2 = el
		}
		subExpr, err = subExpressionRepo.CreateSubExpression(ctx, subExpr)
		if err != nil {
			return nil, err
		}
		return subExpr, nil

	}

	// Разделяем входное выражение на элементы.
	elements := strings.Fields(infixToPostfix(expr.Value))

	for i, element := range elements {
		switch element {
		case "+", "-", "*", "/":
			// Всегда должно быть как минимум два элемента в стеке.
			operand2 := stack[len(stack)-1]
			operand1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2] // Удаляем два элемента из стека.

			isLast := false
			if i == len(elements)-1 {
				isLast = true
			}
			tempVar, err := getTempVar(operand1, operand2, element, isLast)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, tempVar)

			// Заносим результат как временную переменную обратно в стек.
			stack = append(stack, tempVar.Id.String())
		default:
			// Добавляем в стек, если это число.
			stack = append(stack, element)
		}
	}

	return tasks, nil
}
