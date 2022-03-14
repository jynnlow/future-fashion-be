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
	"golang.org/x/crypto/bcrypt"
)

type AdminHandlerActions interface {
	CreateCustomer(w http.ResponseWriter, r *http.Request)
	DeleteCustomer(w http.ResponseWriter, r *http.Request)
	ListCustomers(w http.ResponseWriter, r *http.Request)
	EditCustomers(w http.ResponseWriter, r *http.Request)
}

type AdminHandler struct {
	UserModel       *models.UserCRUDOperationsImpl
	CredentialModel *models.CredentialOperationsImpl
	Logger          *zap.SugaredLogger
}

func (a *AdminHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	tokenKey, err := a.CredentialModel.GetTokenKey()
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

	newCustomer := &dto.UserRequest{}
	err = json.NewDecoder(r.Body).Decode(newCustomer)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	if newCustomer.Username == "" || newCustomer.Password == "" || newCustomer.DOB == "" {
		helpers.JsonResponse(
			w,
			"FAIL",
			"Please fill in all the required information to sign up an accouont",
			nil,
		)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newCustomer.Password), 8)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	newCustomer.Password = string(hashedPassword)
	newCustomer.Role = "customer"

	dbUserRes, err := a.UserModel.Insert(newCustomer)
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
		fmt.Sprintf("%v is inserted successfully", dbUserRes.Username),
		dbUserRes,
	)
}

func (a *AdminHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	tokenKey, err := a.CredentialModel.GetTokenKey()
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

	deletedUser, err := a.UserModel.Delete(uint(uintID))
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
		fmt.Sprintf("%v is deleted successfully", deletedUser.Username),
		deletedUser,
	)
}

func (a *AdminHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	tokenKey, err := a.CredentialModel.GetTokenKey()
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

	users, err := a.UserModel.GetAll()
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	err = helpers.JsonUserList(users, w)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}
}

func (a *AdminHandler) EditCustomer(w http.ResponseWriter, r *http.Request) {
	tokenKey, err := a.CredentialModel.GetTokenKey()
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
	editReq := &models.User{}
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

	// check if customer ID is provided
	if editReq.ID == 0 {
		helpers.JsonResponse(
			w,
			"FAIL",
			"User request ID does not exist",
			nil,
		)
		return
	}

	var hashedPassword []byte
	if editReq.Password != "" {
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(editReq.Password), 8)
		if err != nil {
			helpers.JsonResponse(
				w,
				"FAIL",
				err.Error(),
				nil,
			)
			return
		}
		editReq.Password = string(hashedPassword)
	}

	dbUserRes, err := a.UserModel.Update(editReq)
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
		fmt.Sprintf("%v is updated successfully", dbUserRes.Username),
		dbUserRes,
	)
}
