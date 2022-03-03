package handlers

import (
	"encoding/json"
	"fmt"
	"future-fashion/dto"
	"future-fashion/helpers"
	"future-fashion/models"
	"strconv"

	"net/http"

	"go.uber.org/zap"
)

type ProductHandlerActions interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
	ListProducts(w http.ResponseWriter, r *http.Request)
	EditProduct(w http.ResponseWriter, r *http.Request)
}

type ProductHandler struct {
	ProductModel    *models.ProductCRUDOperationsImpl
	CredentialModel *models.CredentialOperationsImpl
	Logger          *zap.SugaredLogger
}

func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	tokenKey, err := p.CredentialModel.GetTokenKey()
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

	productReq := &dto.ProductRequest{}
	err = json.NewDecoder(r.Body).Decode(productReq)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	err = productReq.Validate()
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	productModel, err := p.convertProductDTOToProductModel(productReq)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	dbProductRes, err := p.ProductModel.Insert(productModel)
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
		fmt.Sprintf("%v is inserted successfully", dbProductRes.Item),
		dbProductRes,
	)
}

func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	tokenKey, err := p.CredentialModel.GetTokenKey()
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

	deletedProduct, err := p.ProductModel.Delete(uint(uintID))
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
		fmt.Sprintf("%v is deleted successfully", deletedProduct.Item),
		deletedProduct,
	)
}

func (p *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// tokenKey, err := p.CredentialModel.GetTokenKey()
	// if err != nil {
	// 	helpers.JsonResponse(
	// 		w,
	// 		"FAIL",
	// 		err.Error(),
	// 		nil,
	// 	)
	// 	return
	// }

	// _, err = helpers.GetVerifiedToken(tokenKey, r)
	// if err != nil {
	// 	helpers.JsonResponse(
	// 		w,
	// 		"FAIL",
	// 		err.Error(),
	// 		nil,
	// 	)
	// 	return
	// }

	products, err := p.ProductModel.GetAll()
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	// err = helpers.JsonUserList(products, w)
	// if err != nil {
	// 	helpers.JsonResponse(
	// 		w,
	// 		"FAIL",
	// 		err.Error(),
	// 		nil,
	// 	)
	// 	return
	// }

	helpers.JsonResponse(
		w,
		"SUCCESS",
		"SUCCESS",
		products,
	)
}

func (p *ProductHandler) EditProduct(w http.ResponseWriter, r *http.Request) {
	tokenKey, err := p.CredentialModel.GetTokenKey()
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

	//dto user does not have all user fields - etc: chest, waist, hip
	editReq := &models.Product{}
	err = json.NewDecoder(r.Body).Decode(editReq)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	//check if customer ID is provided
	if editReq.ID == 0 {
		helpers.JsonResponse(
			w,
			"FAIL",
			"User request ID does not exist",
			nil,
		)
		return
	}

	dbProductRes, err := p.ProductModel.Update(editReq)
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
		fmt.Sprintf("%v is updated successfully", dbProductRes.Item),
		dbProductRes,
	)
}

func (p *ProductHandler) convertProductDTOToProductModel(productReq *dto.ProductRequest) (*models.Product, error) {
	picturesJsonByte, err := json.Marshal(productReq.Pictures)
	if err != nil {
		return nil, err
	}

	xsJsonByte, err := json.Marshal(productReq.XS)
	if err != nil {
		return nil, err
	}

	sJsonByte, err := json.Marshal(productReq.S)
	if err != nil {
		return nil, err
	}

	mJsonByte, err := json.Marshal(productReq.M)
	if err != nil {
		return nil, err
	}

	lJsonByte, err := json.Marshal(productReq.L)
	if err != nil {
		return nil, err
	}

	xlJsonByte, err := json.Marshal(productReq.XL)
	if err != nil {
		return nil, err
	}

	return &models.Product{
		Item:     productReq.Item,
		Price:    productReq.Price,
		Stock:    productReq.Stock,
		Pictures: string(picturesJsonByte),
		XS:       string(xsJsonByte),
		S:        string(sJsonByte),
		M:        string(mJsonByte),
		L:        string(lJsonByte),
		XL:       string(xlJsonByte),
	}, nil
}
