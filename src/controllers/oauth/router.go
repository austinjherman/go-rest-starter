package oauth

import (
	"aherman/src/container"
	"aherman/src/enums"
	"aherman/src/http/response"
	oauthModels "aherman/src/models/oauth"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Router todo
func Router(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		incoming := &oauthModels.RouterRequest{}

		// bind input. if we can't, there was a validation error.
		if err := c.ShouldBindBodyWith(incoming, binding.JSON); err != nil {
			c.Error(err)
			res := response.ErrValidation
			res.Description = err.Error()
			c.AbortWithStatusJSON(res.Status, res)
			return
		}

		if incoming.GrantType == enums.OAuthGrantTypePassword {
			PasswordGrant(app)(c)
			return
		}

		res := response.ErrNotImplemented
		c.AbortWithStatusJSON(res.Status, res)
	}
}