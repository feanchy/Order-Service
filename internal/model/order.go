package model

type Order struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}
