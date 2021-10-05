package oauth

import (
	"aherman/src/http/response"
	"aherman/src/container"

	"github.com/gin-gonic/gin"
)

// AuthorizationCode todo
func AuthorizationCode(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(response.Error(response.ErrNotImplemented))
		return
	}
}