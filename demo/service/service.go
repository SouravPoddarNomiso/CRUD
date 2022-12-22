package service

import (
	"demo/model"
	"demo/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsersInfo() []byte {
	client, ctx, cancel, err := repository.Connect()
	if err != nil {
		panic(err)
	}

	defer repository.Close(client, ctx, cancel)

	var filter, option interface{}
	filter = bson.D{
		//{"maths", bson.D{{"$gt", 70}}},
	}
	option = bson.D{{"_id", 0}}

	cursor, err := repository.Query(client, ctx, "local",
		"users", filter, option)

	if err != nil {
		panic(err)
	}

	var results []bson.D

	if err := cursor.All(ctx, &results); err != nil {

		// handle the error
		panic(err)
	}

	fmt.Println("Query Result")
	for _, doc := range results {
		fmt.Println(doc)
	}
	json, _ := json.Marshal(results)
	return json
}

func SaveUser(userDetails model.User, httpError model.ErrorResponse, response http.ResponseWriter, request *http.Request) {
	httpError.Code = http.StatusBadRequest
	if userDetails.Name == "" {
		httpError.Message = "First Name can't be empty"
		returnErrorResponse(response, request, httpError)
	} else if userDetails.Lname == "" {
		httpError.Message = "Last Name can't be empty"
		returnErrorResponse(response, request, httpError)
	} else if userDetails.Country == "" {
		httpError.Message = "Country can't be empty"
		returnErrorResponse(response, request, httpError)
	} else {
		client, ctx, cancel, err := repository.Connect()
		if err != nil {
			panic(err)
		}
		defer repository.Close(client, ctx, cancel)
		isInserted := repository.InsertOne(client, ctx, "local", "users", userDetails)
		//isInserted := repository.InsertUserInDB(userDetails)
		if isInserted {
			jsonResponse := GetUsersInfo()
			if jsonResponse == nil {
				returnErrorResponse(response, request, httpError)
			} else {
				response.Header().Set("Content-Type", "application/json")
				response.Write(jsonResponse)
			}
		} else {
			returnErrorResponse(response, request, httpError)
		}
	}
}

func DeleteUser(userId string) bool {
	client, ctx, cancel, err := repository.Connect()
	if err != nil {
		panic(err)
	}
	defer repository.Close(client, ctx, cancel)
	intUserId, convErr := strconv.Atoi(userId)
	if convErr != nil {
		return false
	}
	query := bson.D{
		{"id", bson.D{{"$eq", intUserId}}},
	}
	isDeleted := repository.DeleteOne(client, ctx, "local", "users", query)
	if isDeleted {
		jsonResponse := GetUsersInfo()
		if jsonResponse == nil {
			return false
		} else {
			return true
		}
	} else {
		return false
	}

}

func UpdateUser(userDetails model.User) (result *mongo.UpdateResult, mongoErr error) {
	client, ctx, cancel, err := repository.Connect()
	if err != nil {
		panic(err)
	}
	defer repository.Close(client, ctx, cancel)
	filter := bson.D{
		{"id", bson.D{{"$eq", userDetails.ID}}},
	}
	result, mongoErr = repository.UpdateOne(client, ctx, "local", "users", filter, userDetails)
	return
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
