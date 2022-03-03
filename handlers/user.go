package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"future-fashion/dto"
	"future-fashion/helpers"
	"future-fashion/models"

	"go.uber.org/zap"
)

type UserHandlerActions interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	UserModel       *models.UserCRUDOperationsImpl
	CredentialModel *models.CredentialOperationsImpl
	Logger          *zap.SugaredLogger
}

// Sign Up ...
func (u *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	signupReq := &dto.UserRequest{}
	err := json.NewDecoder(r.Body).Decode(signupReq)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	if signupReq.Username == "" || signupReq.Password == "" || signupReq.DOB == "" {
		helpers.JsonResponse(
			w,
			"FAIL",
			"Please fill in all the required information to sign up an accouont",
			nil,
		)
		return
	}

	if signupReq.Role == "" {
		signupReq.Role = "customer"
	} else {
		signupReq.Role = "admin"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupReq.Password), 8)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			err.Error(),
			nil,
		)
		return
	}

	signupReq.Password = string(hashedPassword)

	dbUserRes, err := u.UserModel.Insert(signupReq)
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

// Login ...
func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	foundUser, err := u.UserModel.GetByUsername(loginReq.Username)
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			"NOTE: User does not exist",
			nil,
		)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginReq.Password))
	if err != nil {
		helpers.JsonResponse(
			w,
			"FAIL",
			"NOTE: Incorrect Password",
			nil,
		)
		return
	}

	token := helpers.NewClaim(
		foundUser.ID,
		foundUser.Username,
		foundUser.Role,
	)

	tokenKey, err := u.CredentialModel.GetTokenKey()
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