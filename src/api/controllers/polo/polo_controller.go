package polo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	polo = "polo"
)

// this endpoint will simply return a statusOk
func Polo(c *gin.Context) {
	c.String(http.StatusOK, polo)
}
