package model

type Order struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type OrderCreated struct {
	OrderID int `json:"order_id"`
}
