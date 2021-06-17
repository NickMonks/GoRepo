package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nickmonks/microservices-go/mvc/services"
	"github.com/nickmonks/microservices-go/mvc/utils"
)

// the gin controller needs to implement the interface of a function that takes a *gin.Context
func GetUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("user"), 10, 64)
	if err != nil {
		apiErr := &utils.ApplicationError{
			Message:    "user_id must be a number!",
			StatusCode: http.StatusBadRequest,
			Code:       "not found",
		}

		// this will send the response back to the client
		utils.Respond(c, apiErr.StatusCode, apiErr)
		return // we could not returning and non block the code.
	}

	// call a service, which retrieves the data from the backend
	user, apiErr := services.GetUser(userId)
	if apiErr != nil {
		utils.Respond(c, apiErr.StatusCode, apiErr)
		return
	}

	//return user to client
	// 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
	utils.Respond(c, http.StatusOK, user)

}
