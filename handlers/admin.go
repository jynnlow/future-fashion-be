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
	"gorm.io/gorm"
)

type AdminHandlerActions interface {
	CreateCustomer(w http.ResponseWriter, r *http.Request)
	DeleteCustomer(w http.ResponseWriter, r *http.Request)
	ListCustomers(w http.ResponseWriter, r *http.Request)
	EditCustomers(w http.ResponseWriter, r *http.Request)
	GetCustomerInfo(w http.ResponseWriter, r *http.Request)
	AdminLogin(w http.ResponseWriter, r *http.Request)
}

type AdminHandler struct {
	UserModel       *models.UserCRUDOperationsImpl
	CredentialModel *models.CredentialOperationsImpl
	Logger          *zap.SugaredLogger
}

// Login ...
func (a *AdminHandler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	loginReq := &dto.UserRequest{}
	err := json.NewDecoder(r.Body).Decode(loginReq)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		helpers.JsonResponse(
			w,
			"FAIL",
			"NOTE: Username or password cannot be empty",
			nil,
		)
		return
	}

	foundUser, err := a.UserModel.GetByUsername(loginReq.Username)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			"NOTE: User does not exist. Please create a user account.",
			nil,
		)
		return
	}

	if foundUser.Role != "admin" {
		helpers.JsonResponse(
			w,
			"FAIL",
			"You are not admin",
			nil,
		)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginReq.Password))
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			"NOTE: Incorrect Password. Please try again.",
			nil,
		)
		return
	}

	token := helpers.NewClaim(
		foundUser.ID,
		foundUser.Username,
		foundUser.Role,
	)

	tokenKey, err := a.CredentialModel.GetTokenKey()
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			"NOTE: Failed to get token key",
			nil,
		)
		return
	}

	tokenEncodedString, err := token.CreateToken(tokenKey)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			"NOTE: Failed to create token",
			nil,
		)
		return
	}

	helpers.JsonResponse(
		w,
		"SUCCESS",
		fmt.Sprintf("%v logged in successfully", foundUser.Username),
		tokenEncodedString,
	)
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

	usersResponse := &dto.ListUsersResponse{
		Users: []*dto.UserResponse{},
	}

	for _, user := range users {
		usersResponse.Users = append(usersResponse.Users, &dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Password: user.Password,
			DOB:      user.DOB,
			Role:     user.Role,
			Chest:    user.Chest,
			Waist:    user.Waist,
			Hip:      user.Hip,
		})
	}

	helpers.JsonResponse(
		w,
		"SUCCESS",
		"SUCCESS",
		usersResponse,
	)
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
	editUserReq := &dto.EditUserReq{}
	err = json.NewDecoder(r.Body).Decode(editUserReq)
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
	if editUserReq.ID == 0 {
		helpers.JsonResponse(
			w,
			"FAIL",
			"User request ID does not exist",
			nil,
		)
		return
	}

	var hashedPassword []byte
	if editUserReq.Password != "" {
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(editUserReq.Password), 8)
		if err != nil {
			helpers.JsonResponse(
				w,
				"FAIL",
				err.Error(),
				nil,
			)
			return
		}
		editUserReq.Password = string(hashedPassword)
	}

	userModel, err := a.convertEditUserDTOToUserModel(editUserReq)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	dbUserRes, err := a.UserModel.Update(userModel)
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

func (a *AdminHandler) GetCustomerInfo(w http.ResponseWriter, r *http.Request) {
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

	foundUser, err := a.UserModel.GetByID(uint(uintID))
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	userRes := &dto.UserResponse{
		ID:       foundUser.ID,
		Username: foundUser.Username,
		Password: foundUser.Password,
		DOB:      foundUser.DOB,
		Role:     foundUser.Role,
		Chest:    foundUser.Chest,
		Waist:    foundUser.Waist,
		Hip:      foundUser.Hip,
	}

	helpers.JsonResponse(
		w,
		"SUCCESS",
		fmt.Sprintf("%v is deleted successfully", userRes.Username),
		userRes,
	)
}

func (a *AdminHandler) convertEditUserDTOToUserModel(userReq *dto.EditUserReq) (*models.User, error) {
	return &models.User{
		Model: gorm.Model{
			ID: userReq.ID,
		},
		Username: userReq.Username,
		Password: userReq.Password,
		Role:     userReq.Role,
		DOB:      userReq.DOB,
		Chest:    userReq.Chest,
		Waist:    userReq.Waist,
		Hip:      userReq.Hip,
	}, nil
}
