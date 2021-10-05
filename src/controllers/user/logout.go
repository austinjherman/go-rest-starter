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

		err := app.User.RevokeSession(app.Current.User.ID, app.Current.Token.SessionID)
		if err != nil {
			c.Error(err)
			res := response.ErrLogout
			c.JSON(response.Error(res))
			return
		}

		c.JSON(response.Success(response.SuccessLogout, struct{}{}))

	}
}