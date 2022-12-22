package controller

import (
	"demo/model"
	"demo/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func Home(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "../views/index.html")
}

func GetUsers(response http.ResponseWriter, request *http.Request) {
	var httpError = model.ErrorResponse{
		Code: http.StatusInternalServerError, Message: "Error getting user details.",
	}
	jsonResponse := service.GetUsersInfo()

	if jsonResponse == nil {
		returnErrorResponse(response, request, httpError)
	} else {
		response.Header().Set("Content-Type", "application/json")
		response.Write(jsonResponse)
	}
}

func InsertUser(response http.ResponseWriter, request *http.Request) {
	var httpError = model.ErrorResponse{
		Code: http.StatusInternalServerError, Message: "Error while inserting user details.",
	}
	var userDetails model.User
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&userDetails)
	defer request.Body.Close()
	if err != nil {
		returnErrorResponse(response, request, httpError)
	} else {
		service.SaveUser(userDetails, httpError, response, request)
	}
}

func DeleteUser(response http.ResponseWriter, request *http.Request) {
	var httpError = model.ErrorResponse{
		Code: http.StatusInternalServerError, Message: "Error deleting user.",
	}
	userID := mux.Vars(request)["id"]
	if userID == "" {
		httpError.Message = "User id can't be empty"
		returnErrorResponse(response, request, httpError)
	} else {
		isdeleted := service.DeleteUser(userID)
		if isdeleted {
			GetUsers(response, request)
		} else {
			returnErrorResponse(response, request, httpError)
		}
	}
}

func UpdateUser(response http.ResponseWriter, request *http.Request) {
	var httpError = model.ErrorResponse{
		Code: http.StatusInternalServerError, Message: "Error in updating user.",
	}
	var userDetails model.User
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&userDetails)
	defer request.Body.Close()
	if err != nil {
		returnErrorResponse(response, request, httpError)
	} else {
		httpError.Code = http.StatusBadRequest
		if userDetails.Name == "" {
			httpError.Message = "First Name can't be empty"
			returnErrorResponse(response, request, httpError)
		} else if userDetails.ID == 0 {
			httpError.Message = "user Id can't be empty"
			returnErrorResponse(response, request, httpError)
		} else if userDetails.Lname == "" {
			httpError.Message = "Last Name can't be empty"
			returnErrorResponse(response, request, httpError)
		} else if userDetails.Country == "" {
			httpError.Message = "Country can't be empty"
			returnErrorResponse(response, request, httpError)
		} else {
			_, mongoErr := service.UpdateUser(userDetails)
			if mongoErr != nil {
				returnErrorResponse(response, request, httpError)
			} else {
				GetUsers(response, request)
			}
		}
	}
}

func returnErrorResponse(response http.ResponseWriter, request *http.Request, errorMesage model.ErrorResponse) {
	httpResponse := &model.ErrorResponse{Code: errorMesage.Code, Message: errorMesage.Message}
	jsonResponse, err := json.Marshal(httpResponse)
	if err != nil {
		panic(err)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(errorMesage.Code)
	response.Write(jsonResponse)
}
