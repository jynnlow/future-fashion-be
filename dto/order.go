package dto

type OrderRequest struct {
	Total     float32      `json:"total"`
	Status    string       `json:"status"`
	Snapshots []*CartModel `json:"snapshots"`
	UserID    uint         `json:"user_id"`
}

type OrderResponse struct {
	ID        uint         `json:"id,omitempty"`
	Total     float32      `json:"total"`
	Status    string       `json:"status"`
	Snapshots []*CartModel `json:"snapshots,omitempty"`
}

type CartModel struct {
	Id       string         `json:"id"`
	Item     string         `json:"item"`
	Price    float32        `json:"price"`
	Sizing   string         `json:"sizing"`
	Quantity int            `json:"quantity"`
	Product  ProductRequest `json:"product"`
}

type ListOrdersResponse struct {
	Orders []*OrderResponse `json:"orders"`
}
