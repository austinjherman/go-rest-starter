package user

import (
	"aherman/src/container"
	"aherman/src/http/response"

	"github.com/gin-gonic/gin"
)

// Logout removes the user's access token and its associated
// refresh token from the token whitelist.
func Logout(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		err := app.Facades.User.RevokeTokensBySession(app.Current.User.ID, app.Current.Token.SessionID)
		ok, httpResponse := app.Facades.Error.ShouldContinue(err, &response.ErrUnknown)
		if !ok {
			c.JSON(response.Error(httpResponse))
			return
		}

		c.JSON(response.Success(response.SuccessLogout, struct{}{}))
	}
}