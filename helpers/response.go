package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func JsonResponse(w http.ResponseWriter, status, message string, details interface{}) {
	var err error
	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	response := &Response{
		Status:  status,
		Message: message,
		Details: details,
	}

	jsonRes, err := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonRes)
}

func JsonUserList(objectList interface{}, w http.ResponseWriter) error {
	jsonRes, err := json.Marshal(objectList)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	_, err = w.Write(jsonRes)
	if err != nil {
		return err
	}

	return nil
}
