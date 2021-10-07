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
		err := c.ShouldBindBodyWith(incoming, binding.JSON)
		ok, httpResponse := app.Facades.Error.ShouldContinue(err, &response.ErrValidation)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		if incoming.GrantType == enums.OAuthGrantTypePassword {
			PasswordGrant(app)(c)
			return
		}

		res := response.ErrNotImplemented
		c.JSON(res.Status, res)
	}
}