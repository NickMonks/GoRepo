package repositories

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nickmonks/microservices-go/src/api/domain/repositories"
	"github.com/nickmonks/microservices-go/src/api/services"
	"github.com/nickmonks/microservices-go/src/api/utils/errors"
)

// Handles POST requests from /repositories endpoint
func CreateRepo(c *gin.Context) {
	var request repositories.CreateRepoRequest

	// will try to bind an empty struct to check if its bindable
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := errors.NewBadRequestError("Invalid JSON body")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Remember: RepositoryService is instantiated in the init() function, created when imported
	// and is an interface that implements createRepo.
	result, err := services.RepositoryService.CreateRepo(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, result)
}
