package models

type ExpressionState string

const (
	ExpressionError      ExpressionState = "error"
	ExpressionInProgress                 = "in_progress"
	ExpressionOk                         = "ok"
)

type Expression struct {
	Result         float64         `json:"result"`
	Id             string          `json:"id"`
	UserId         string          `json:"userId"`
	IdempotencyKey string          `json:"idempotencyKey"`
	Value          string          `json:"value"`
	State          ExpressionState `json:"state"`
}
