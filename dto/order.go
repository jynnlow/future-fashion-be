package dto

import "time"

type OrderRequest struct {
	Total     float32      `json:"total"`
	Status    string       `json:"status"`
	Snapshots []*CartModel `json:"snapshots"`
	UserID    uint         `json:"user_id"`
}

type EditOrderRequest struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
	UserID uint   `json:"user_id"`
}

type OrderResponse struct {
	ID        uint         `json:"id,omitempty"`
	Total     float32      `json:"total"`
	Status    string       `json:"status"`
	Snapshots []*CartModel `json:"snapshots,omitempty"`
	UserID    uint         `json:"userID"`
	CreatedAt time.Time    `json:"createdAt"`
}

type CartModel struct {
	Id       string          `json:"id"`
	Item     string          `json:"item"`
	Price    float32         `json:"price"`
	Sizing   string          `json:"sizing"`
	Quantity int             `json:"quantity"`
	Product  ProductResponse `json:"product"`
}

type ListOrdersResponse struct {
	Orders []*OrderResponse `json:"orders"`
}
