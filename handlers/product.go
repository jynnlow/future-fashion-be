package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"future-fashion/dto"
	"future-fashion/helpers"
	"future-fashion/models"
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
	productsResponse := &dto.ListProductsResponse{
		Products: []*dto.ProductResponse{},
	}

	for _, product := range products {
		var picturesList []string
		err := json.Unmarshal([]byte(product.Pictures), &picturesList)
		if err != nil {
			helpers.JsonResponse(
				w,
				"FAIL",
				err.Error(),
				nil,
			)
			return
		}

		xsModel, err := unmarshalSizing(product.XS)
		if err != nil {
			helpers.JsonResponse(
				w,
				"FAIL",
				err.Error(),
				nil,
			)
			return
		}
		sModel, err := unmarshalSizing(product.S)
		if err != nil {
			helpers.JsonResponse(
				w,
				"FAIL",
				err.Error(),
				nil,
			)
			return
		}
		mModel, err := unmarshalSizing(product.M)
		if err != nil {
			helpers.JsonResponse(
				w,
				"FAIL",
				err.Error(),
				nil,
			)
			return
		}
		lModel, err := unmarshalSizing(product.L)
		if err != nil {
			helpers.JsonResponse(
				w,
				"FAIL",
				err.Error(),
				nil,
			)
			return
		}
		xlModel, err := unmarshalSizing(product.XL)
		if err != nil {
			helpers.JsonResponse(
				w,
				"FAIL",
				err.Error(),
				nil,
			)
			return
		}

		productsResponse.Products = append(productsResponse.Products, &dto.ProductResponse{
			ID:       product.ID,
			Item:     product.Item,
			Price:    product.Price,
			Stock:    product.Stock,
			Pictures: picturesList,
			XS:       xsModel,
			S:        sModel,
			M:        mModel,
			L:        lModel,
			XL:       xlModel,
		})
	}

	helpers.JsonResponse(
		w,
		"SUCCESS",
		"SUCCESS",
		productsResponse,
	)
}

func unmarshalSizing(sizing string) (*dto.Sizing, error) {
	var sizingModel *dto.Sizing
	err := json.Unmarshal([]byte(sizing), &sizingModel)
	if err != nil {
		return nil, err
	}
	return sizingModel, nil
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
	updateProductReq := &dto.UpdateProductRequest{}
	err = json.NewDecoder(r.Body).Decode(updateProductReq)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	//check if product ID is provided
	if updateProductReq.ID == 0 {
		helpers.JsonResponse(
			w,
			"FAIL",
			"Product request ID does not exist",
			nil,
		)
		return
	}

	productModel, err := p.convertUpdateProductDTOToProductModel(updateProductReq)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	dbProductRes, err := p.ProductModel.Update(productModel)
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

func (p *ProductHandler) convertUpdateProductDTOToProductModel(productReq *dto.UpdateProductRequest) (*models.Product, error) {
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
		Model: gorm.Model{
			ID: productReq.ID,
		},
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
