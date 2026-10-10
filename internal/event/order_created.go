package event

type OrderCreated struct {
	EventID string `json:"event_id"`
	OrderID int    `json:"order_id"`
}
