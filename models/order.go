package models

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderCRUDOperation interface {
	GetByID(id uint) (*Order, error)
	GetByUserID(user_id uint) ([]*Order, error)
	GetAll() ([]*Order, error)
	Insert(*Order) (*Product, error)
	Delete(id uint) (*Product, error)
	Update(orderReq *Order) (*Order, error)
}

type Order struct {
	gorm.Model
	Total     float32 `json:"total"`
	Status    string  `json:"status"`
	Snapshots string  `json:"snapshots"`
	UserID    uint    `json:"user_id"`
}

type OrderCRUDOperationsImpl struct {
	DB     *gorm.DB
	Logger *zap.SugaredLogger
}

func (o *OrderCRUDOperationsImpl) GetByID(id uint) (*Order, error) {
	order := &Order{}
	err := o.DB.First(order, id).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *OrderCRUDOperationsImpl) GetAll() ([]*Order, error) {
	var order []*Order
	err := o.DB.Find(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *OrderCRUDOperationsImpl) GetByUserID(user_id uint) ([]*Order, error) {
	var order []*Order
	err := o.DB.Where("user_id = ?", user_id).Find(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *OrderCRUDOperationsImpl) Insert(order *Order) (*Order, error) {
	err := o.DB.Create(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *OrderCRUDOperationsImpl) Delete(id uint) (*Order, error) {
	foundOrder, err := o.GetByID(id)
	if err != nil {
		return nil, err
	}
	//delete only if the user exists
	//permanently deleted with Unscoped().Delete()
	err = o.DB.Unscoped().Delete(foundOrder, id).Error
	if err != nil {
		return nil, err
	}
	return foundOrder, nil
}

func (o *OrderCRUDOperationsImpl) Update(orderReq *Order) (*Order, error) {
	foundOrder, err := o.GetByID(orderReq.ID)
	if err != nil {
		return nil, err
	}
	// if orderReq.Total != 0 && foundOrder.Total != orderReq.Total {
	// 	foundOrder.Total = orderReq.Total
	// }
	if orderReq.Status != "" && foundOrder.Status != orderReq.Status {
		foundOrder.Status = orderReq.Status
	}
	// if orderReq.Snapshots != "" && foundOrder.Snapshots != orderReq.Snapshots {
	// 	foundOrder.Snapshots = orderReq.Snapshots
	// }

	//update order with all field except userID
	err = o.DB.Save(foundOrder).Error
	if err != nil {
		return nil, err
	}
	return foundOrder, nil
}
