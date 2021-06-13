package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nickmonks/microservices-go/mvc/services"
	"github.com/nickmonks/microservices-go/mvc/utils"
)

func GetUser(resp http.ResponseWriter, req *http.Request) {
	userIdParam := req.URL.Query().Get("user")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		apiErr := &utils.ApplicationError{
			Message:    "user_id must be a number!",
			StatusCode: http.StatusBadRequest,
			Code:       "not found",
		}

		jsonValue, _ := json.Marshal(apiErr)
		resp.WriteHeader(apiErr.StatusCode)
		resp.Write(jsonValue)
		return
	}

	// call a service, which retrieves the data from the backend
	user, apiErr := services.GetUser(userId)
	if apiErr != nil {
		//TODO: Handle error
		resp.WriteHeader(apiErr.StatusCode)
		resp.Write([]byte(err.Error()))
		return
	}

	//return user to client
	jsonValue, _ := json.Marshal(user)
	resp.Write(jsonValue)

}
