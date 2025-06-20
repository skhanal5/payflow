package model

type OrderResponse struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
}

type Item struct {
	Product string  `json:"product"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int              `json:"total"`
}