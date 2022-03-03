package models

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductCRUDOperation interface {
	GetByID(id uint) (*Product, error)
	GetAll() ([]*Product, error)
	Insert(*Product) (*Product, error)
	Delete(id uint) (*Product, error)
	Update(productReq *Product) (*Product, error)
}

type Product struct {
	gorm.Model
	Item     string  `json:"item"`
	Price    float32 `json:"price"`
	Stock    int     `json:"stock"`
	Pictures string  `json:"picture"`
	XS       string  `json:"xs"`
	S        string  `json:"s"`
	M        string  `json:"m"`
	L        string  `json:"l"`
	XL       string  `json:"xl"`
}

type ProductCRUDOperationsImpl struct {
	DB     *gorm.DB
	Logger *zap.SugaredLogger
}

func (p *ProductCRUDOperationsImpl) GetByID(id uint) (*Product, error) {
	product := &Product{}
	err := p.DB.First(product, id).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductCRUDOperationsImpl) GetAll() ([]*Product, error) {
	var product []*Product
	err := p.DB.Find(&product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductCRUDOperationsImpl) Insert(product *Product) (*Product, error) {
	err := p.DB.Create(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductCRUDOperationsImpl) Delete(id uint) (*Product, error) {
	foundProduct, err := p.GetByID(id)
	if err != nil {
		return nil, err
	}
	//delete only if the user exists
	//permanently deleted with Unscoped().Delete()
	err = p.DB.Unscoped().Delete(foundProduct, id).Error
	if err != nil {
		return nil, err
	}
	return foundProduct, nil
}

func (p *ProductCRUDOperationsImpl) Update(productReq *Product) (*Product, error) {
	foundProduct, err := p.GetByID(productReq.ID)
	if err != nil {
		return nil, err
	}
	if productReq.Item != "" && foundProduct.Item != productReq.Item {
		foundProduct.Item = productReq.Item
	}
	if productReq.Price != 0 && foundProduct.Price != productReq.Price {
		foundProduct.Price = productReq.Price
	}
	if productReq.Stock != 0 && foundProduct.Stock != productReq.Stock {
		foundProduct.Stock = productReq.Stock
	}
	if productReq.Pictures != "" && foundProduct.Pictures != productReq.Pictures {
		foundProduct.Pictures = productReq.Pictures
	}
	if productReq.XS != "" && foundProduct.XS != productReq.XS {
		foundProduct.XS = productReq.XS
	}
	if productReq.S != "" && foundProduct.S != productReq.S {
		foundProduct.S = productReq.S
	}
	if productReq.M != "" && foundProduct.M != productReq.M {
		foundProduct.M = productReq.M
	}
	if productReq.L != "" && foundProduct.L != productReq.L {
		foundProduct.L = productReq.L
	}
	if productReq.XL != "" && foundProduct.XL != productReq.XL {
		foundProduct.XL = productReq.XL
	}

	//update user with all field
	err = p.DB.Save(foundProduct).Error
	if err != nil {
		return nil, err
	}
	return foundProduct, nil
}
