package handlers

import (
	"encoding/json"
	"fmt"
	"future-fashion/dto"
	"future-fashion/helpers"
	"future-fashion/models"
	"net/http"
	"strconv"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderHandlerActions interface {
	CreateOrder(w http.ResponseWriter, r *http.Request)
	DeleteOrder(w http.ResponseWriter, r *http.Request)
	ListOrders(w http.ResponseWriter, r *http.Request)
	ListOrdersByUserID(w http.ResponseWriter, r *http.Request)
	EditOrderStatus(w http.ResponseWriter, r *http.Request)
}

type OrderHandler struct {
	OrderModel      *models.OrderCRUDOperationsImpl
	CredentialModel *models.CredentialOperationsImpl
	Logger          *zap.SugaredLogger
}

func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	tokenKey, err := o.CredentialModel.GetTokenKey()
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	verifiedToken, err := helpers.GetVerifiedToken(tokenKey, r)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}
	var orderReq *dto.OrderRequest
	err = json.NewDecoder(r.Body).Decode(&orderReq)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}
	orderReq.UserID = verifiedToken.Id
	orderReq.Status = "Order is comfirmed"

	orderModel, err := o.convertOrderDTOToOrderModel(orderReq)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	dbOrderRes, err := o.OrderModel.Insert(orderModel)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	helpers.JsonResponse(
		w,
		"SUCCESS",
		fmt.Sprintf("%v is inserted successfully", dbOrderRes.ID),
		convertOrderModelToCreateOrderRes(dbOrderRes),
	)
}

func (o *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	tokenKey, err := o.CredentialModel.GetTokenKey()
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	verifiedToken, err := helpers.GetVerifiedToken(tokenKey, r)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	if verifiedToken.Role != "admin" {
		helpers.JsonResponse(
			w,
			"FAIL",
			"NOTE: Only admin is allowed for this operation",
			nil,
		)
		return
	}

	//retrieve parameter from url
	param, ok := r.URL.Query()["id"]
	if !ok || len(param[0]) < 1 {
		helpers.JsonResponse(
			w,
			"FAIL",
			"Url param key not exist",
			nil,
		)
		return
	}

	// convert id to uint64 type
	uintID, err := strconv.ParseUint(param[0], 10, 64)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	deletedOrder, err := o.OrderModel.Delete(uint(uintID))
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	helpers.JsonResponse(
		w,
		"SUCCESS",
		fmt.Sprintf("%v is deleted successfully", deletedOrder),
		deletedOrder,
	)

}

func (o *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	tokenKey, err := o.CredentialModel.GetTokenKey()
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	verifiedToken, err := helpers.GetVerifiedToken(tokenKey, r)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	if verifiedToken.Role != "admin" {
		helpers.JsonResponse(
			w,
			"FAIL",
			"NOTE: Only admin is allowed for this operation",
			nil,
		)
		return
	}

	orders, err := o.OrderModel.GetAll()
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	orderResponse := &dto.ListOrdersResponse{
		Orders: []*dto.OrderResponse{},
	}
	for _, order := range orders {
		snapshotsObj, err := unmarshalSnapshots(order.Snapshots)
		if err != nil {
			helpers.JsonResponse(
				w,
				"FAIL",
				err.Error(),
				nil,
			)
			return
		}
		orderResponse.Orders = append(orderResponse.Orders, &dto.OrderResponse{
			ID:        order.ID,
			Total:     order.Total,
			Status:    order.Status,
			Snapshots: snapshotsObj,
			UserID:    order.UserID,
			CreatedAt: order.CreatedAt,
		})
	}

	helpers.JsonResponse(
		w,
		"SUCCESS",
		"SUCCESS",
		orderResponse,
	)
}

func (o *OrderHandler) ListOrdersByUserID(w http.ResponseWriter, r *http.Request) {
	tokenKey, err := o.CredentialModel.GetTokenKey()
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	verifiedToken, err := helpers.GetVerifiedToken(tokenKey, r)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	orders, err := o.OrderModel.GetByUserID(verifiedToken.Id)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	orderResponse := &dto.ListOrdersResponse{
		Orders: []*dto.OrderResponse{},
	}
	for _, order := range orders {
		snapshotsObj, err := unmarshalSnapshots(order.Snapshots)
		if err != nil {
			helpers.JsonResponse(
				w,
				"FAIL",
				err.Error(),
				nil,
			)
			return
		}
		orderResponse.Orders = append(orderResponse.Orders, &dto.OrderResponse{
			ID:        order.ID,
			Total:     order.Total,
			Status:    order.Status,
			Snapshots: snapshotsObj,
			CreatedAt: order.CreatedAt,
		})
	}

	helpers.JsonResponse(
		w,
		"SUCCESS",
		"SUCCESS",
		orderResponse,
	)
}

func (o *OrderHandler) EditOrderStatus(w http.ResponseWriter, r *http.Request) {
	tokenKey, err := o.CredentialModel.GetTokenKey()
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	verifiedToken, err := helpers.GetVerifiedToken(tokenKey, r)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	if verifiedToken.Role != "admin" {
		helpers.JsonResponse(
			w,
			"FAIL",
			"NOTE: Only admin is allowed for this operation",
			nil,
		)
		return
	}

	updateOrderReq := &dto.EditOrderRequest{}
	err = json.NewDecoder(r.Body).Decode(updateOrderReq)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	//check if order ID is provided
	if updateOrderReq.ID == 0 {
		helpers.JsonResponse(
			w,
			"FAIL",
			"Order request ID does not exist",
			nil,
		)
		return
	}

	orderModel, err := o.convertEditOrderDTOToOrderModel(updateOrderReq)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	dbOrderRes, err := o.OrderModel.Update(orderModel)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	helpers.JsonResponse(
		w,
		"SUCCESS",
		"SUCCESS",
		dbOrderRes,
	)
}

func (o *OrderHandler) convertEditOrderDTOToOrderModel(orderReq *dto.EditOrderRequest) (*models.Order, error) {
	return &models.Order{
		Model: gorm.Model{
			ID: orderReq.ID,
		},
		Status: orderReq.Status,
		UserID: orderReq.UserID,
	}, nil
}

func (o *OrderHandler) convertOrderDTOToOrderModel(orderReq *dto.OrderRequest) (*models.Order, error) {
	snapshotJsonByte, err := json.Marshal(orderReq.Snapshots)
	if err != nil {
		return nil, err
	}

	return &models.Order{
		Total:     orderReq.Total,
		Status:    orderReq.Status,
		Snapshots: string(snapshotJsonByte),
		UserID:    orderReq.UserID,
	}, nil
}

func convertOrderModelToCreateOrderRes(orderModel *models.Order) *dto.OrderResponse {
	return &dto.OrderResponse{
		Total:  orderModel.Total,
		Status: orderModel.Status,
	}
}

func unmarshalSnapshots(snapshots string) ([]*dto.CartModel, error) {
	var cartModels []*dto.CartModel
	err := json.Unmarshal([]byte(snapshots), &cartModels)
	if err != nil {
		return nil, err
	}
	return cartModels, nil
}
