package user

import (
	"aherman/src/container"
	"aherman/src/http/response"
	userModels "aherman/src/models/user"

	"github.com/gin-gonic/gin"
)

// Delete removes a user from the database
func Delete(app *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			user       *userModels.User   = app.Current.User
			userPublic *userModels.Public = &userModels.Public{}
		)

		userPublic.BindAttributes(user)

		result := app.Facades.User.DB.Delete(&userModels.User{}, user.ID)
		ok, httpRequest := app.Facades.Error.ShouldContinue(result.Error, &response.ErrUnknown)
		if !ok {
			c.JSON(response.Error(httpRequest))
			return
		}

		c.JSON(response.Success(response.SuccessDelete, userPublic))
	}
}