package models

import "github.com/google/uuid"

type SubExpression struct {
	Id               uuid.UUID     `json:"id" pg:"type:uuid"`
	ExpressionId     uuid.UUID     `json:"expressionId" pg:"type:uuid"`
	Val1             float64       `json:"val1"`
	Val2             float64       `json:"val2"`
	SubExpressionId1 uuid.NullUUID `json:"subExpressionId1" pg:"type:uuid"`
	SubExpressionId2 uuid.NullUUID `json:"subExpressionId2" pg:"type:uuid"`
	Action           string        `json:"action"`
	Result           float64       `json:"result"`
	IsLast           bool          `json:"isLast"`
	Error            bool          `json:"error"`
	AgentId          uuid.NullUUID `json:"agentId"`
}
