package handlers

import (
	"future-fashion/models"
	"net/http"

	"go.uber.org/zap"
)

type OrderHandlerActions interface {
	CreateOrder(w http.ResponseWriter, r *http.Request)
}

type OrderHandler struct {
	OrderModel      *models.OrderCRUDOperationsImpl
	CredentialModel *models.CredentialOperationsImpl
	Logger          *zap.SugaredLogger
}

func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {

}
