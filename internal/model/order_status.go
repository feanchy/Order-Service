package model

const (
	OrderStatusCreated    = "created"
	OrderStatusConfirmed  = "confirmed"
	OrderStatusInDelivery = "in_delivery"
	OrderStatusDelivered  = "delivered"
	OrderStatusCancelled  = "cancelled"
)

func CanTransition(from, to string) bool {
	// chatgpt, у нас тут проверка from тк каждый статус может быть только в 2-х состояниях? то есть либо в след статусе либо отмененном
	//
	switch from {
	case OrderStatusCreated:
		return to == OrderStatusConfirmed || to == OrderStatusCancelled
	case OrderStatusConfirmed:
		return to == OrderStatusInDelivery || to == OrderStatusCancelled
	case OrderStatusInDelivery:
		return to == OrderStatusDelivered
	}

	return false
}
