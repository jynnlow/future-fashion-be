package models

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Total     float32 `json:"total"`
	Status    string  `json:"status"`
	Snapshots string  `json:"snapshots"`
	UserID    uint
}

type OrderCRUDOperationsImpl struct {
	DB     *gorm.DB
	Logger *zap.SugaredLogger
}
