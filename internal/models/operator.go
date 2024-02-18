package models

import "time"

type Operator struct {
	Op      string        `json:"op"`
	Timeout time.Duration `json:"timeout"`
}
