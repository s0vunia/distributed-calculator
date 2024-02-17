package models

type Agent struct {
	Id        string `json:"id"`
	Heartbeat int64  `json:"heartbeat"`
}
