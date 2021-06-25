package app

import (
	"github.com/gin-gonic/gin"
	"github.com/nickmonks/microservices-go/src/api/log"
)

var (
	router *gin.Engine
)

// IMPORTANT: init function gets call when it is imported into another package.
// because app is being imported from main, it will have tge router initializations
func init() {
	// Defaults wraps the gin.New() and allows recovery, which recovers from panic
	router = gin.Default()
}

func StartApp() {
	log.Log.Info("About to map the urls")
	mapUrls()

	// use router.Run
	if err := router.Run(":8081"); err != nil {
		panic(err)
	}
}
