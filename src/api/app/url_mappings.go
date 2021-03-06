package app

import (
	"github.com/nickmonks/microservices-go/src/api/controllers/polo"
	"github.com/nickmonks/microservices-go/src/api/controllers/repositories"
)

func mapUrls() {
	// we use the router from init
	router.POST("/repository", repositories.CreateRepo)

	// Route to create several repositories
	router.POST("/repositories", repositories.CreateRepos)

	// create a symbolic GET request for AWS services; to check if server is responding, we need to do "marco/polo"
	router.GET("/marco", polo.Polo)
}
