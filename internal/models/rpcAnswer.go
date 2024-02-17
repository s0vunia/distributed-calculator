package models

import "github.com/google/uuid"

type RPCAnswer struct {
	IdSubExpression uuid.UUID `json:"idSubExpression"`
	IdAgent         uuid.UUID `json:"idAgent"`
}
