package user

import (
	"aherman/src/container"
	"aherman/src/http/response"

	"github.com/gin-gonic/gin"
)

// LogoutAll removes all user tokens from the whitelist, effectively
// logging them out on all devices
func LogoutAll(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		err := app.Facades.User.RevokeTokensByUser(app.Current.User.ID)
		ok, httpResponse := app.Facades.Error.ShouldContinue(err, &response.ErrUnknown)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		c.JSON(response.Success(response.SuccessLogout, struct{}{}))
	}
}