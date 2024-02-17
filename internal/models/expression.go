package models

type ExpressionState string

const (
	Error      ExpressionState = "error"
	InProgress                 = "in_progress"
	Ok                         = "ok"
)

type Expression struct {
	Result         float64         `json:"result"`
	Id             string          `json:"id"`
	IdempotencyKey string          `json:"idempotencyKey"`
	Value          string          `json:"value"`
	State          ExpressionState `json:"state"`
}
