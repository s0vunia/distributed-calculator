package agent

import (
	"errors"
	"myproject/internal/config"
	"myproject/internal/models"
	"time"
)

// Calculate считает subexpression с паузой из config
func Calculate(expression *models.SubExpression) (ans float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = errors.New("unknown panic")
			}
		}
	}()
	switch expression.Action {
	case "+":
		<-time.After(config.TimeCalculatePlus)
		return expression.Val1 + expression.Val2, nil
	case "-":
		<-time.After(config.TimeCalculateMinus)
		return expression.Val1 - expression.Val2, nil
	case "*":
		<-time.After(config.TimeCalculateMult)
		return expression.Val1 * expression.Val2, nil
	case "/":
		if expression.Val2 == 0 {
			return 0, errors.New("cannot divide by zero")
		}
		<-time.After(config.TimeCalculateDivide)
		return expression.Val1 / expression.Val2, nil
	default:
		err = errors.New("not allowed action")
		return 0, err
	}
}
