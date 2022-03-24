package dto

import "errors"

type ProductRequest struct {
	Item     string   `json:"item"`
	Price    float32  `json:"price"`
	Stock    int      `json:"stock"`
	Pictures []string `json:"pictures"`
	XS       *Sizing  `json:"xs"`
	S        *Sizing  `json:"s"`
	M        *Sizing  `json:"m"`
	L        *Sizing  `json:"l"`
	XL       *Sizing  `json:"xl"`
}

type UpdateProductRequest struct {
	ID       uint     `json:"id"`
	Item     string   `json:"item"`
	Price    float32  `json:"price"`
	Stock    int      `json:"stock"`
	Pictures []string `json:"pictures"`
	XS       *Sizing  `json:"xs"`
	S        *Sizing  `json:"s"`
	M        *Sizing  `json:"m"`
	L        *Sizing  `json:"l"`
	XL       *Sizing  `json:"xl"`
}

type ProductResponse struct {
	ID       uint     `json:"id"`
	Item     string   `json:"item"`
	Price    float32  `json:"price"`
	Stock    int      `json:"stock"`
	Pictures []string `json:"pictures"`
	XS       *Sizing  `json:"xs"`
	S        *Sizing  `json:"s"`
	M        *Sizing  `json:"m"`
	L        *Sizing  `json:"l"`
	XL       *Sizing  `json:"xl"`
}

type ListProductsResponse struct {
	Products []*ProductResponse `json:"products"`
}

type Sizing struct {
	Chest float32 `json:"chest"`
	Waist float32 `json:"waist"`
	Hip   float32 `json:"hip"`
}

func (p *ProductRequest) Validate() error {
	if p.Item == "" || p.Price == 0 {
		return errors.New("item name or price cannot be empty")
	}
	return nil
}
