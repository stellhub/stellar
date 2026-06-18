package internal

type createOrderRequest struct {
	CustomerID string  `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Channel    string  `json:"channel"`
}

type createOrderResponse struct {
	OrderID string  `json:"order_id"`
	Status  string  `json:"status"`
	Amount  float64 `json:"amount"`
	Channel string  `json:"channel"`
}
