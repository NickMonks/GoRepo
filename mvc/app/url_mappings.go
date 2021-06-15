package app

import (
	"github.com/nickmonks/microservices-go/mvc/controllers"
)

func mapUrls() {
	// we use the router from init
	router.GET("/users/:user", controllers.GetUser)
}
