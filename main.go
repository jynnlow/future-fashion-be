package main

import (
	"fmt"
	"log"
	"net/http"

	"future-fashion/handlers"
	"future-fashion/helpers"
	"future-fashion/infra"
	"future-fashion/models"

	"github.com/gorilla/mux"
)

func main() {
	// Init DB
	db, err := infra.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// Init Logger
	logger := helpers.InitLogger()

	// Init Models
	userModel := &models.UserCRUDOperationsImpl{
		DB:     db,
		Logger: logger,
	}

	credentialModel := &models.CredentialOperationsImpl{
		DB: db,
	}

	productModel := &models.ProductCRUDOperationsImpl{
		DB:     db,
		Logger: logger,
	}

	orderModel := &models.OrderCRUDOperationsImpl{
		DB:     db,
		Logger: logger,
	}

	// Init Handlers
	userHandler := &handlers.UserHandler{
		UserModel:       userModel,
		CredentialModel: credentialModel,
		Logger:          logger,
	}

	adminHandler := &handlers.AdminHandler{
		UserModel:       userModel,
		CredentialModel: credentialModel,
		Logger:          logger,
	}

	productHandler := &handlers.ProductHandler{
		ProductModel:    productModel,
		CredentialModel: credentialModel,
		Logger:          logger,
	}

	orderHandler := &handlers.OrderHandler{
		OrderModel:      orderModel,
		CredentialModel: credentialModel,
		Logger:          logger,
	}

	r := mux.NewRouter()
	//User Handlers
	r.HandleFunc("/user/signup", userHandler.SignUp).Methods("POST")
	r.HandleFunc("/user/login", userHandler.Login).Methods("POST")
	//Admin Handlers
	r.HandleFunc("/admin/create-customer", adminHandler.CreateCustomer).Methods("POST")
	r.HandleFunc("/admin/delete-customer", adminHandler.DeleteCustomer).Methods("DELETE")
	r.HandleFunc("/admin/list-customers", adminHandler.ListCustomers).Methods("GET")
	r.HandleFunc("/admin/edit-customer", adminHandler.EditCustomer).Methods("PATCH")
	//Product Handlers
	r.HandleFunc("/product/create-product", productHandler.CreateProduct).Methods("POST")
	r.HandleFunc("/product/delete-product", productHandler.DeleteProduct).Methods("DELETE")
	r.HandleFunc("/product/list-products", productHandler.ListProducts).Methods("GET")
	r.HandleFunc("/product/edit-product", productHandler.EditProduct).Methods("PATCH")

	r.HandleFunc("/product/create-order", orderHandler.CreateOrder).Methods("POST")

	fmt.Println("HTTP server running on http://127.0.0.1:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
